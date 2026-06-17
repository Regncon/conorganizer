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
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5/middleware"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestManager_EnsureConnection_WhenCookieMissing_CreatesSessionCookieAndLiveKey(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a request without the live session cookie.",
		When:  "When the manager ensures the connection.",
		Then:  "Then it creates a connection id, saves the cookie, and stores a live key.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing live session without a matching KV key.",
		When:  "When the manager ensures the connection.",
		Then:  "Then the same connection id is reused and the KV key is recreated.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the live manager bucket configuration.",
		When:  "When each bucket config is built.",
		Then:  "Then every live bucket uses the configured 26 hour TTL.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a live bucket with existing connection keys.",
		When:  "When the manager broadcasts the bucket.",
		Then:  "Then every connection key receives a fresh timestamp value.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an empty live bucket.",
		When:  "When the manager broadcasts the bucket.",
		Then:  "Then no error is returned.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a watcher for a connection key.",
		When:  "When the manager broadcasts the bucket.",
		Then:  "Then the watcher receives an update for that key.",
	})

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
		`"nats_touch_duration_ms":`,
	} {
		if !strings.Contains(logOutput, expectedPart) {
			t.Fatalf("expected log output to contain %q, got %q", expectedPart, logOutput)
		}
	}
}

func TestManager_Stream_WhenWatcherAlreadyStopped_LogsAtInfo(t *testing.T) {
	// Given a watcher that has already stopped by the time stream cleanup runs,
	// when the live stream exits,
	// then the expected cleanup failure is logged at info level, not warn level.

	// Given
	manager := newTestManager(t)
	var logs bytes.Buffer
	manager.logger = slog.New(slog.NewJSONHandler(&logs, nil)).With("component", "live")
	kv := mustFakeKeyValue(t, manager, BucketEvents)
	kv.watcherStopErr = nats.ErrBadSubscription

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/live", nil)
	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()
	kv.onWatch = cancel
	request = request.WithContext(context.WithValue(ctx, middleware.RequestIDKey, "request-456"))

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
	logOutput := logs.String()
	for _, expectedPart := range []string{
		`"level":"INFO"`,
		`"msg":"live watcher already stopped: nats: invalid subscription"`,
		`"method":"GET"`,
		`"path":"/live"`,
		`"request_id":"request-456"`,
		`"buckets":["events"]`,
	} {
		if !strings.Contains(logOutput, expectedPart) {
			t.Fatalf("expected log output to contain %q, got %q", expectedPart, logOutput)
		}
	}
	if strings.Contains(logOutput, `"level":"WARN"`) {
		t.Fatalf("expected watcher cleanup log not to be warning, got %q", logOutput)
	}
}

func TestManager_Stream_WhenWatcherCloses_ExitsStreamForReconnect(t *testing.T) {
	// Given a live stream with a watcher that closes unexpectedly,
	// when the watcher closes while the request is still active,
	// then the stream exits so Datastar can reconnect.

	// Given
	manager := newTestManager(t)
	kv := mustFakeKeyValue(t, manager, BucketEvents)
	watcherClosed := make(chan struct{})
	kv.onWatchWatcher = func(watcher *fakeWatcher) {
		watcher.close()
		close(watcherClosed)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/live", nil)
	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()
	request = request.WithContext(ctx)

	done := make(chan struct{})

	// When
	go func() {
		manager.Stream(recorder, request, Page{
			Buckets: []Bucket{BucketEvents},
			Render: func(ctx context.Context, r *http.Request) templ.Component {
				return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					_, err := io.WriteString(w, `<div id="live-content">Ready</div>`)
					return err
				})
			},
		})
		close(done)
	}()

	<-watcherClosed

	// Then
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		cancel()
		<-done
		t.Fatalf("expected stream to exit when watcher closes")
	}
}

func TestDatastarInit_ReturnsRestartResilientGetExpression(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a live endpoint path.",
		When:  "When the Datastar init expression is generated.",
		Then:  "Then it includes retry settings that survive server restarts.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a live endpoint URL expression.",
		When:  "When the Datastar init expression is generated.",
		Then:  "Then the expression is passed through and retry settings are included.",
	})

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
