package live

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	nats "github.com/nats-io/nats.go"
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

func WithLogger(logger *slog.Logger) Option {
	return func(m *Manager) {
		if logger != nil {
			m.logger = logger.With("component", "live")
		}
	}
}

type Manager struct {
	store    sessions.Store
	buckets  map[Bucket]keyValue
	ttl      time.Duration
	now      func() time.Time
	logger   *slog.Logger
	natsConn *nats.Conn
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
		store:    store,
		buckets:  make(map[Bucket]keyValue),
		ttl:      DefaultTTL,
		now:      time.Now,
		logger:   slog.Default().With("component", "live"),
		natsConn: nc,
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
		m.logRequestError(r, buckets, "failed to ensure live session", err)
		return "", err
	}

	touchStart := time.Now()
	if err := m.touchConnection(r.Context(), connectionID, buckets...); err != nil {
		m.logLiveKeyTouchError(r, buckets, "failed to touch live key", err, time.Since(touchStart))
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
			return fmt.Errorf("touch live key in bucket %s: %w", bucket, err)
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
				return fmt.Errorf("broadcast live key in bucket %s: %w", bucket, err)
			}
		}
	}

	return nil
}

func (m *Manager) Stream(w http.ResponseWriter, r *http.Request, page Page) {
	connectionID, err := m.ensureSession(w, r)
	if err != nil {
		m.logRequestError(r, page.Buckets, "failed to ensure live session", err)
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
		m.logRequestError(r, page.Buckets, "failed to send initial live patch", err)
		_ = sse.ConsoleError(err)
		return
	}

	touchStart := time.Now()
	if err := m.touchConnection(ctx, connectionID, page.Buckets...); err != nil {
		m.logLiveKeyTouchError(r, page.Buckets, "failed to touch live key before watching", err, time.Since(touchStart))
		_ = sse.ConsoleError(err)
		return
	}

	watchers := make([]bucketWatcher, 0, len(page.Buckets))
	for _, bucket := range page.Buckets {
		kv, err := m.keyValue(bucket)
		if err != nil {
			m.logRequestError(r, []Bucket{bucket}, "failed to prepare live watcher", err)
			_ = sse.ConsoleError(err)
			return
		}
		watcher, err := kv.Watch(ctx, connectionID, jetstream.UpdatesOnly())
		if err != nil {
			m.logRequestError(r, []Bucket{bucket}, "failed to start live watcher", err)
			_ = sse.ConsoleError(err)
			return
		}
		watchers = append(watchers, bucketWatcher{bucket: bucket, watcher: watcher})
	}
	defer func() {
		for _, watcher := range watchers {
			if err := watcher.watcher.Stop(); err != nil {
				if errors.Is(err, nats.ErrBadSubscription) {
					m.logRequestInfo(r, []Bucket{watcher.bucket}, "live watcher already stopped", err)
					continue
				}
				m.logRequestWarn(r, []Bucket{watcher.bucket}, "failed to stop live watcher", err)
			}
		}
	}()

	updates := make(chan struct{}, 1)
	watcherClosed := make(chan Bucket, len(watchers))
	for _, watcher := range watchers {
		go forwardWatcherUpdates(ctx, watcher, updates, watcherClosed)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case bucket := <-watcherClosed:
			if ctx.Err() != nil {
				return
			}
			m.logRequestWarnMessage(r, []Bucket{bucket}, "live watcher closed; closing stream for reconnect")
			return
		case <-updates:
			if err := patch(); err != nil {
				m.logRequestError(r, page.Buckets, "failed to send live patch after update", err)
				_ = sse.ConsoleError(err)
				return
			}
		}
	}
}

type bucketWatcher struct {
	bucket  Bucket
	watcher jetstream.KeyWatcher
}

func (m *Manager) logRequestError(r *http.Request, buckets []Bucket, message string, err error) {
	m.log().Error(fmt.Errorf("%s: %w", message, err).Error(), liveRequestLogArgs(r, buckets)...)
}

func (m *Manager) logRequestWarn(r *http.Request, buckets []Bucket, message string, err error) {
	m.log().Warn(fmt.Errorf("%s: %w", message, err).Error(), liveRequestLogArgs(r, buckets)...)
}

func (m *Manager) logRequestWarnMessage(r *http.Request, buckets []Bucket, message string) {
	args := liveRequestLogArgs(r, buckets)
	args = append(args, m.natsLogArgs()...)
	m.log().Warn(message, args...)
}

func (m *Manager) logRequestInfo(r *http.Request, buckets []Bucket, message string, err error) {
	m.log().Info(fmt.Errorf("%s: %w", message, err).Error(), liveRequestLogArgs(r, buckets)...)
}

func (m *Manager) logLiveKeyTouchError(r *http.Request, buckets []Bucket, message string, err error, duration time.Duration) {
	args := liveRequestLogArgs(r, buckets)
	args = append(args, "nats_touch_duration_ms", duration.Milliseconds())
	args = append(args, m.natsLogArgs()...)
	m.log().Error(fmt.Errorf("%s: %w", message, err).Error(), args...)
}

func (m *Manager) log() *slog.Logger {
	if m.logger != nil {
		return m.logger
	}
	return slog.Default().With("component", "live")
}

func (m *Manager) natsLogArgs() []any {
	if m.natsConn == nil {
		return nil
	}

	args := []any{"nats_status", m.natsConn.Status().String()}
	if lastErr := m.natsConn.LastError(); lastErr != nil {
		args = append(args, "nats_last_error", lastErr.Error())
	}
	return args
}

func liveRequestLogArgs(r *http.Request, buckets []Bucket) []any {
	args := make([]any, 0, 8)
	if r != nil {
		args = append(args, "method", r.Method)
		if r.URL != nil {
			args = append(args, "path", r.URL.Path)
		}
		if requestID := middleware.GetReqID(r.Context()); requestID != "" {
			args = append(args, "request_id", requestID)
		}
	}
	if len(buckets) > 0 {
		args = append(args, "buckets", bucketNames(buckets))
	}
	return args
}

func bucketNames(buckets []Bucket) []string {
	names := make([]string, 0, len(buckets))
	for _, bucket := range buckets {
		names = append(names, string(bucket))
	}
	return names
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

func forwardWatcherUpdates(ctx context.Context, watcher bucketWatcher, updates chan<- struct{}, watcherClosed chan<- Bucket) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-watcher.watcher.Updates():
			if !ok {
				select {
				case <-ctx.Done():
				case watcherClosed <- watcher.bucket:
				}
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
