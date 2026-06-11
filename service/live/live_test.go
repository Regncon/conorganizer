package live

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
)

func TestManager_EnsureConnection_WhenCookieMissing_CreatesSessionCookieAndLiveKey(t *testing.T) {
	// Given a request without the live session cookie,
	// when the manager ensures the connection,
	// then it creates a connection id, saves the cookie, and stores a live key.

	// Given
	expectedBucket := BucketEvents
	manager := newTestManager(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/live", nil)

	// When
	connectionID, err := manager.EnsureConnection(recorder, request, expectedBucket)

	// Then
	if err != nil {
		t.Fatalf("expected connection ensure to succeed: %v", err)
	}
	if connectionID == "" {
		t.Fatalf("expected connection id to be set")
	}
	assertResponseHasCookie(t, recorder, "connections")
	assertLiveKeyExists(t, manager, expectedBucket, connectionID)
}

func TestManager_EnsureConnection_WhenCookieExistsAndKeyMissing_RecreatesLiveKey(t *testing.T) {
	// Given an existing live session without a matching KV key,
	// when the manager ensures the connection,
	// then the same connection id is reused and the KV key is recreated.

	// Given
	expectedBucket := BucketEvents
	manager := newTestManager(t)
	firstRecorder := httptest.NewRecorder()
	firstRequest := httptest.NewRequest(http.MethodGet, "/live", nil)
	expectedConnectionID, err := manager.EnsureConnection(firstRecorder, firstRequest, expectedBucket)
	if err != nil {
		t.Fatalf("expected initial connection ensure to succeed: %v", err)
	}
	sessionCookie := responseCookie(t, firstRecorder, "connections")
	kv := mustFakeKeyValue(t, manager, expectedBucket)
	if err := kv.Purge(context.Background(), expectedConnectionID); err != nil {
		t.Fatalf("purge live key: %v", err)
	}

	secondRecorder := httptest.NewRecorder()
	secondRequest := httptest.NewRequest(http.MethodGet, "/live", nil)
	secondRequest.AddCookie(sessionCookie)

	// When
	actualConnectionID, err := manager.EnsureConnection(secondRecorder, secondRequest, expectedBucket)

	// Then
	if err != nil {
		t.Fatalf("expected connection ensure to succeed: %v", err)
	}
	if actualConnectionID != expectedConnectionID {
		t.Fatalf("connection id mismatch\nexpected: %s\nactual:   %s", expectedConnectionID, actualConnectionID)
	}
	assertLiveKeyExists(t, manager, expectedBucket, expectedConnectionID)
}

func TestManager_BucketConfig_UsesTwentySixHourTTLForEveryLiveBucket(t *testing.T) {
	// Given the live manager bucket configuration,
	// when each bucket config is built,
	// then every live bucket uses the configured 26 hour TTL.

	// Given
	expectedTTL := DefaultTTL
	manager := &Manager{ttl: expectedTTL}

	// When / Then
	for _, bucket := range allBuckets {
		config := manager.bucketConfig(bucket)
		if config.TTL != expectedTTL {
			t.Fatalf("TTL mismatch for bucket %s\nexpected: %s\nactual:   %s", bucket, expectedTTL, config.TTL)
		}
	}
}

func TestManager_Broadcast_WhenBucketHasKeys_WritesTimestampToEveryLiveKey(t *testing.T) {
	// Given a live bucket with existing connection keys,
	// when the manager broadcasts the bucket,
	// then every connection key receives a fresh timestamp value.

	// Given
	expectedKeys := []string{"connection-a", "connection-b"}
	manager := newTestManager(t)
	kv := mustFakeKeyValue(t, manager, BucketEvents)
	for _, key := range expectedKeys {
		if _, err := kv.Put(context.Background(), key, []byte("old-value")); err != nil {
			t.Fatalf("put fixture key %s: %v", key, err)
		}
	}

	// When
	err := manager.Broadcast(context.Background(), BucketEvents)

	// Then
	if err != nil {
		t.Fatalf("expected broadcast to succeed: %v", err)
	}
	for _, key := range expectedKeys {
		entry, err := kv.Get(context.Background(), key)
		if err != nil {
			t.Fatalf("get broadcast key %s: %v", key, err)
		}
		if string(entry.Value()) == "old-value" {
			t.Fatalf("expected key %s to receive new value", key)
		}
		assertTimestampValue(t, entry.Value())
	}
}

func TestManager_Broadcast_WhenBucketHasNoKeys_Succeeds(t *testing.T) {
	// Given an empty live bucket,
	// when the manager broadcasts the bucket,
	// then no error is returned.

	// Given
	manager := newTestManager(t)

	// When
	err := manager.Broadcast(context.Background(), BucketEvents)

	// Then
	if err != nil {
		t.Fatalf("expected empty bucket broadcast to succeed: %v", err)
	}
}

func TestManager_Broadcast_WhenWatcherIsOpen_SendsUpdateToWatcher(t *testing.T) {
	// Given a watcher for a connection key,
	// when the manager broadcasts the bucket,
	// then the watcher receives an update for that key.

	// Given
	expectedKey := "connection-a"
	manager := newTestManager(t)
	kv := mustFakeKeyValue(t, manager, BucketEvents)
	if _, err := kv.Put(context.Background(), expectedKey, []byte("old-value")); err != nil {
		t.Fatalf("put fixture key: %v", err)
	}
	watcher, err := kv.Watch(context.Background(), expectedKey, jetstream.UpdatesOnly())
	if err != nil {
		t.Fatalf("watch fixture key: %v", err)
	}
	defer func() {
		if err := watcher.Stop(); err != nil {
			t.Fatalf("stop watcher: %v", err)
		}
	}()

	// When
	if err := manager.Broadcast(context.Background(), BucketEvents); err != nil {
		t.Fatalf("expected broadcast to succeed: %v", err)
	}

	// Then
	entry := waitForWatcherUpdate(t, watcher)
	if entry.Key() != expectedKey {
		t.Fatalf("watcher key mismatch\nexpected: %s\nactual:   %s", expectedKey, entry.Key())
	}
	assertTimestampValue(t, entry.Value())
}

func TestManager_Stream_WhenTouchConnectionFails_SendsInitialPatch(t *testing.T) {
	// Given live key storage is temporarily unavailable,
	// when a live stream starts,
	// then the initial page patch is still sent instead of failing the HTTP request.

	// Given
	manager := newTestManager(t)
	var logs bytes.Buffer
	manager.logger = slog.New(slog.NewJSONHandler(&logs, nil)).With("component", "live")
	kv := mustFakeKeyValue(t, manager, BucketEvents)
	kv.putErr = errors.New("live key store unavailable")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/live", nil)
	request = request.WithContext(context.WithValue(request.Context(), middleware.RequestIDKey, "request-123"))

	// When
	manager.Stream(recorder, request, Page{
		Buckets: []Bucket{BucketEvents},
		Render: func(ctx context.Context, r *http.Request) templ.Component {
			return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<div id="live-content">Ready</div>`)
				return err
			})
		},
	})

	// Then
	result := recorder.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected live stream to start successfully, got status %d", result.StatusCode)
	}
	body := recorder.Body.String()
	for _, expectedPart := range []string{"datastar-patch-elements", "live-content", "Ready"} {
		if !strings.Contains(body, expectedPart) {
			t.Fatalf("expected stream body to contain %q, got %q", expectedPart, body)
		}
	}
	logOutput := logs.String()
	for _, expectedPart := range []string{
		`"component":"live"`,
		`"msg":"failed to touch live key before watching: touch live key in bucket events: live key store unavailable"`,
		`"method":"GET"`,
		`"path":"/live"`,
		`"request_id":"request-123"`,
		`"buckets":["events"]`,
	} {
		if !strings.Contains(logOutput, expectedPart) {
			t.Fatalf("expected log output to contain %q, got %q", expectedPart, logOutput)
		}
	}
}

func TestDatastarInit_ReturnsRestartResilientGetExpression(t *testing.T) {
	// Given a live endpoint path,
	// when the Datastar init expression is generated,
	// then it includes retry settings that survive server restarts.

	// Given
	expectedParts := []string{
		"@get('/root/api'",
		"requestCancellation: 'disabled'",
		"retryMaxCount: Infinity",
		"retryInterval: 1000",
		"retryMaxWaitMs: 30000",
	}

	// When
	actual := DatastarInit("/root/api")

	// Then
	for _, expectedPart := range expectedParts {
		if !strings.Contains(actual, expectedPart) {
			t.Fatalf("Datastar init expression missing %q in %q", expectedPart, actual)
		}
	}
}

func TestDatastarInitExpression_ReturnsRestartResilientGetExpressionWithDynamicURL(t *testing.T) {
	// Given a live endpoint URL expression,
	// when the Datastar init expression is generated,
	// then the expression is passed through and retry settings are included.

	// Given
	expectedParts := []string{
		"@get('/profile/api' + window.location.search",
		"requestCancellation: 'disabled'",
		"retryMaxCount: Infinity",
		"retryInterval: 1000",
		"retryMaxWaitMs: 30000",
	}

	// When
	actual := DatastarInitExpression("'/profile/api' + window.location.search")

	// Then
	for _, expectedPart := range expectedParts {
		if !strings.Contains(actual, expectedPart) {
			t.Fatalf("Datastar init expression missing %q in %q", expectedPart, actual)
		}
	}
}

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
	mu       sync.Mutex
	bucket   Bucket
	ttl      time.Duration
	values   map[string][]byte
	watchers map[string][]*fakeWatcher
	revision uint64
	putErr   error
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
	defer kv.mu.Unlock()

	watcher := &fakeWatcher{updates: make(chan jetstream.KeyValueEntry, 16)}
	kv.watchers[key] = append(kv.watchers[key], watcher)
	return watcher, nil
}

type fakeWatcher struct {
	updates chan jetstream.KeyValueEntry
}

func (w *fakeWatcher) Updates() <-chan jetstream.KeyValueEntry {
	return w.updates
}

func (w *fakeWatcher) Stop() error {
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
