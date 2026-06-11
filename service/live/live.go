package live

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

const (
	sessionName  = "connections"
	sessionIDKey = "id"

	DefaultTTL = 26 * time.Hour
	maxBytes   = 16 * 1024 * 1024
)

type Bucket string

const (
	BucketEvents         Bucket = "events"
	BucketInterests      Bucket = "interests"
	BucketBillettholders Bucket = "billettholders"
	BucketRooms          Bucket = "rooms"
)

var allBuckets = []Bucket{BucketEvents, BucketInterests, BucketBillettholders, BucketRooms}

type Option func(*Manager)

func WithTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.ttl = ttl
	}
}

type Manager struct {
	store   sessions.Store
	buckets map[Bucket]keyValue
	ttl     time.Duration
	now     func() time.Time
}

type keyValue interface {
	Put(ctx context.Context, key string, value []byte) (uint64, error)
	Keys(ctx context.Context, opts ...jetstream.WatchOpt) ([]string, error)
	Watch(ctx context.Context, keys string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error)
}

type Page struct {
	Buckets []Bucket
	Render  func(ctx context.Context, r *http.Request) templ.Component
}

func NewManager(ctx context.Context, ns *embeddednats.Server, store sessions.Store, opts ...Option) (*Manager, error) {
	nc, err := ns.Client()
	if err != nil {
		return nil, fmt.Errorf("create nats client: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("create jetstream client: %w", err)
	}

	manager := &Manager{
		store:   store,
		buckets: make(map[Bucket]keyValue),
		ttl:     DefaultTTL,
		now:     time.Now,
	}
	for _, opt := range opts {
		opt(manager)
	}

	for _, bucket := range allBuckets {
		kv, err := js.CreateOrUpdateKeyValue(ctx, manager.bucketConfig(bucket))
		if err != nil {
			return nil, fmt.Errorf("create live bucket %s: %w", bucket, err)
		}
		manager.buckets[bucket] = kv
	}

	return manager, nil
}

func (m *Manager) EnsureConnection(w http.ResponseWriter, r *http.Request, buckets ...Bucket) (string, error) {
	connectionID, err := m.ensureSession(w, r)
	if err != nil {
		return "", err
	}

	if err := m.touchConnection(r.Context(), connectionID, buckets...); err != nil {
		return "", err
	}

	return connectionID, nil
}

func (m *Manager) ensureSession(w http.ResponseWriter, r *http.Request) (string, error) {
	sess, err := m.store.Get(r, sessionName)
	if err != nil {
		return "", fmt.Errorf("get live session: %w", err)
	}

	connectionID, ok := sess.Values[sessionIDKey].(string)
	if !ok || connectionID == "" {
		connectionID = uuid.NewString()
		sess.Values[sessionIDKey] = connectionID
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("save live session: %w", err)
		}
	}

	return connectionID, nil
}

func (m *Manager) touchConnection(ctx context.Context, connectionID string, buckets ...Bucket) error {
	value := m.liveValue()
	for _, bucket := range buckets {
		kv, err := m.keyValue(bucket)
		if err != nil {
			return err
		}
		if _, err := kv.Put(ctx, connectionID, value); err != nil {
			return fmt.Errorf("touch live key %s in bucket %s: %w", connectionID, bucket, err)
		}
	}

	return nil
}

func (m *Manager) Broadcast(ctx context.Context, buckets ...Bucket) error {
	for _, bucket := range buckets {
		kv, err := m.keyValue(bucket)
		if err != nil {
			return err
		}

		keys, err := kv.Keys(ctx)
		if err != nil {
			if errors.Is(err, jetstream.ErrNoKeysFound) {
				continue
			}
			return fmt.Errorf("list live keys in bucket %s: %w", bucket, err)
		}

		for _, key := range keys {
			if _, err := kv.Put(ctx, key, m.liveValue()); err != nil {
				return fmt.Errorf("broadcast live key %s in bucket %s: %w", key, bucket, err)
			}
		}
	}

	return nil
}

func (m *Manager) Stream(w http.ResponseWriter, r *http.Request, page Page) {
	connectionID, err := m.ensureSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	ctx := r.Context()

	patch := func() error {
		if page.Render == nil {
			return fmt.Errorf("live page renderer is nil")
		}
		return sse.PatchElementTempl(page.Render(ctx, r))
	}

	if err := patch(); err != nil {
		_ = sse.ConsoleError(err)
		return
	}

	if err := m.touchConnection(ctx, connectionID, page.Buckets...); err != nil {
		_ = sse.ConsoleError(err)
		return
	}

	watchers := make([]jetstream.KeyWatcher, 0, len(page.Buckets))
	for _, bucket := range page.Buckets {
		kv, err := m.keyValue(bucket)
		if err != nil {
			_ = sse.ConsoleError(err)
			return
		}
		watcher, err := kv.Watch(ctx, connectionID, jetstream.UpdatesOnly())
		if err != nil {
			_ = sse.ConsoleError(err)
			return
		}
		watchers = append(watchers, watcher)
	}
	defer func() {
		for _, watcher := range watchers {
			_ = watcher.Stop()
		}
	}()

	updates := make(chan struct{}, 1)
	for _, watcher := range watchers {
		go forwardWatcherUpdates(ctx, watcher, updates)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-updates:
			if err := patch(); err != nil {
				_ = sse.ConsoleError(err)
				return
			}
		}
	}
}

func DatastarInit(url string) string {
	escapedURL := strings.ReplaceAll(url, "'", "\\'")
	return DatastarInitExpression(fmt.Sprintf("'%s'", escapedURL))
}

func DatastarInitExpression(urlExpression string) string {
	return fmt.Sprintf(
		"@get(%s, {requestCancellation: 'disabled', retryMaxCount: Infinity, retryInterval: 1000, retryMaxWaitMs: 30000})",
		urlExpression,
	)
}

func (m *Manager) keyValue(bucket Bucket) (keyValue, error) {
	kv, ok := m.buckets[bucket]
	if !ok {
		return nil, fmt.Errorf("unknown live bucket: %s", bucket)
	}
	return kv, nil
}

func (m *Manager) liveValue() []byte {
	now := m.now
	if now == nil {
		now = time.Now
	}
	return []byte(now().UTC().Format(time.RFC3339Nano))
}

func (m *Manager) bucketConfig(bucket Bucket) jetstream.KeyValueConfig {
	ttl := m.ttl
	if ttl == 0 {
		ttl = DefaultTTL
	}
	return jetstream.KeyValueConfig{
		Bucket:      string(bucket),
		Description: fmt.Sprintf("Conorganizer live update bucket: %s", bucket),
		Compression: true,
		TTL:         ttl,
		MaxBytes:    maxBytes,
	}
}

func forwardWatcherUpdates(ctx context.Context, watcher jetstream.KeyWatcher, updates chan<- struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-watcher.Updates():
			if !ok {
				return
			}
			if entry == nil {
				continue
			}
			select {
			case updates <- struct{}{}:
			default:
			}
		}
	}
}
