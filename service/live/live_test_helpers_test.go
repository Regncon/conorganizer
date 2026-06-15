package live

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()

	store := sessions.NewCookieStore([]byte("live-test-session-secret"))
	store.MaxAge(int((24 * time.Hour) / time.Second))

	buckets := make(map[Bucket]keyValue)
	for _, bucket := range allBuckets {
		buckets[bucket] = newFakeKeyValue(bucket, DefaultTTL)
	}

	return &Manager{
		store:   store,
		buckets: buckets,
		ttl:     DefaultTTL,
		now:     time.Now,
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func mustFakeKeyValue(t *testing.T, manager *Manager, bucket Bucket) *fakeKeyValue {
	t.Helper()

	kv, err := manager.keyValue(bucket)
	if err != nil {
		t.Fatalf("get key value bucket %s: %v", bucket, err)
	}
	fakeKV, ok := kv.(*fakeKeyValue)
	if !ok {
		t.Fatalf("expected fake key value bucket %s, got %T", bucket, kv)
	}
	return fakeKV
}

func assertResponseHasCookie(t *testing.T, recorder *httptest.ResponseRecorder, name string) {
	t.Helper()

	_ = responseCookie(t, recorder, name)
}

func responseCookie(t *testing.T, recorder *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()

	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("expected response cookie %q", name)
	return nil
}

func assertLiveKeyExists(t *testing.T, manager *Manager, bucket Bucket, key string) {
	t.Helper()

	kv := mustFakeKeyValue(t, manager, bucket)
	entry, err := kv.Get(context.Background(), key)
	if err != nil {
		t.Fatalf("expected live key %s in bucket %s: %v", key, bucket, err)
	}
	assertTimestampValue(t, entry.Value())
}

func assertTimestampValue(t *testing.T, value []byte) {
	t.Helper()

	if _, err := time.Parse(time.RFC3339Nano, string(value)); err != nil {
		t.Fatalf("expected RFC3339Nano timestamp value, got %q: %v", string(value), err)
	}
}

type fakeKeyValue struct {
	mu             sync.Mutex
	bucket         Bucket
	ttl            time.Duration
	values         map[string][]byte
	watchers       map[string][]*fakeWatcher
	revision       uint64
	putErr         error
	watcherStopErr error
	onWatch        func()
}

func newFakeKeyValue(bucket Bucket, ttl time.Duration) *fakeKeyValue {
	return &fakeKeyValue{
		bucket:   bucket,
		ttl:      ttl,
		values:   make(map[string][]byte),
		watchers: make(map[string][]*fakeWatcher),
	}
}

func (kv *fakeKeyValue) Get(_ context.Context, key string) (jetstream.KeyValueEntry, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	value, ok := kv.values[key]
	if !ok {
		return nil, jetstream.ErrKeyNotFound
	}
	return &fakeEntry{
		bucket:   string(kv.bucket),
		key:      key,
		value:    cloneBytes(value),
		revision: kv.revision,
		created:  time.Now(),
	}, nil
}

func (kv *fakeKeyValue) Put(_ context.Context, key string, value []byte) (uint64, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if kv.putErr != nil {
		return 0, kv.putErr
	}

	kv.revision++
	stored := cloneBytes(value)
	kv.values[key] = stored
	entry := &fakeEntry{
		bucket:   string(kv.bucket),
		key:      key,
		value:    cloneBytes(stored),
		revision: kv.revision,
		created:  time.Now(),
	}
	for _, watcher := range kv.watchers[key] {
		watcher.send(entry)
	}
	return kv.revision, nil
}

func (kv *fakeKeyValue) Purge(_ context.Context, key string, _ ...jetstream.KVDeleteOpt) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.values, key)
	return nil
}

func (kv *fakeKeyValue) Keys(_ context.Context, _ ...jetstream.WatchOpt) ([]string, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if len(kv.values) == 0 {
		return nil, jetstream.ErrNoKeysFound
	}
	keys := make([]string, 0, len(kv.values))
	for key := range kv.values {
		keys = append(keys, key)
	}
	return keys, nil
}

func (kv *fakeKeyValue) Watch(_ context.Context, key string, _ ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	kv.mu.Lock()

	watcher := &fakeWatcher{
		updates: make(chan jetstream.KeyValueEntry, 16),
		stopErr: kv.watcherStopErr,
	}
	kv.watchers[key] = append(kv.watchers[key], watcher)
	onWatch := kv.onWatch
	kv.mu.Unlock()

	if onWatch != nil {
		onWatch()
	}
	return watcher, nil
}

type fakeWatcher struct {
	updates chan jetstream.KeyValueEntry
	stopErr error
}

func (w *fakeWatcher) Updates() <-chan jetstream.KeyValueEntry {
	return w.updates
}

func (w *fakeWatcher) Stop() error {
	if w.stopErr != nil {
		return w.stopErr
	}
	close(w.updates)
	return nil
}

func (w *fakeWatcher) send(entry jetstream.KeyValueEntry) {
	select {
	case w.updates <- entry:
	default:
	}
}

type fakeEntry struct {
	bucket   string
	key      string
	value    []byte
	revision uint64
	created  time.Time
}

func (e *fakeEntry) Bucket() string                  { return e.bucket }
func (e *fakeEntry) Key() string                     { return e.key }
func (e *fakeEntry) Value() []byte                   { return cloneBytes(e.value) }
func (e *fakeEntry) Revision() uint64                { return e.revision }
func (e *fakeEntry) Created() time.Time              { return e.created }
func (e *fakeEntry) Delta() uint64                   { return 0 }
func (e *fakeEntry) Operation() jetstream.KeyValueOp { return jetstream.KeyValuePut }

func cloneBytes(value []byte) []byte {
	cloned := make([]byte, len(value))
	copy(cloned, value)
	return cloned
}

func waitForWatcherUpdate(t *testing.T, watcher jetstream.KeyWatcher) jetstream.KeyValueEntry {
	t.Helper()

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	for {
		select {
		case entry := <-watcher.Updates():
			if entry != nil {
				return entry
			}
		case <-timer.C:
			t.Fatalf("timed out waiting for watcher update")
		}
	}
}
