package checkIn

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

const successfulCheckinResponse = `{"data":{"eventTickets":[{"id":1,"order_id":2,"category":"Festivalpass","category_id":3,"crm":{"first_name":"Kari","last_name":"Nordmann","email":"kari@example.com","born":"1990-01-01"}}]}}`

func TestIsOver18_WhenBirthdayIsOnConventionEnd_ReturnsTrue(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Given a person who turns eighteen on the last day of Regncon.", When: "When their age is checked.", Then: "Then they count as over eighteen for the convention."})
	if !isOver18("2008-10-04") {
		t.Fatal("expected birthday on convention end to count as over eighteen")
	}
}

func TestIsOver18_WhenBirthdayIsAfterConventionEnd_ReturnsFalse(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Given a person who turns eighteen after the last day of Regncon.", When: "When their age is checked.", Then: "Then they do not count as over eighteen for the convention."})
	if isOver18("2008-10-05") {
		t.Fatal("expected birthday after convention end not to count as over eighteen")
	}
}

func TestFetchTicketsFromCheckIn_TimeoutIsBoundedAndSanitized(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-release }))
	defer func() { close(release); server.Close() }()

	_, err := fetchTicketsFromCheckIn(context.Background(), &http.Client{Timeout: 20 * time.Millisecond}, server.URL+"?client_secret=must-not-leak")
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if outcome, _ := fetchFailureDetails(err); outcome != "timeout" {
		t.Fatalf("outcome mismatch\nexpected: timeout\nactual: %s", outcome)
	}
	if strings.Contains(err.Error(), "client_secret") || strings.Contains(err.Error(), server.URL) {
		t.Fatalf("network error leaked endpoint credentials: %q", err)
	}
}

func TestFetchTicketsFromCheckIn_ContextCancellationIsPropagated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := fetchTicketsFromCheckIn(ctx, &http.Client{}, server.URL)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context\nactual: %v", err)
	}
	if outcome, _ := fetchFailureDetails(err); outcome != "canceled" {
		t.Fatalf("outcome mismatch\nexpected: canceled\nactual: %s", outcome)
	}
}

func TestFetchTicketsFromCheckIn_RejectsInvalidResponses(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		body, outcome string
	}{
		{name: "http 500", status: http.StatusInternalServerError, outcome: "http_error"},
		{name: "http 401", status: http.StatusUnauthorized, outcome: "http_error"},
		{name: "graphql errors", status: http.StatusOK, body: `{"errors":[{"message":"not authorized"}]}`, outcome: "graphql_error"},
		{name: "malformed JSON", status: http.StatusOK, body: `{`, outcome: "decode_error"},
		{name: "missing ticket data", status: http.StatusOK, body: `{"data":{}}`, outcome: "decode_error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()
			_, err := fetchTicketsFromCheckIn(context.Background(), server.Client(), server.URL)
			if err == nil {
				t.Fatal("expected fetch error")
			}
			if outcome, status := fetchFailureDetails(err); outcome != test.outcome || (test.status >= 300 && status != test.status) {
				t.Fatalf("failure details mismatch\nexpected outcome/status: %s/%d\nactual outcome/status: %s/%d", test.outcome, test.status, outcome, status)
			}
		})
	}
}

func TestCacheGet_SuccessfulRefreshPopulatesCache(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		_, _ = w.Write([]byte(successfulCheckinResponse))
	}))
	defer server.Close()
	cache := newTestCache(server.URL, server.Client())

	result, err := cache.Get(context.Background(), testutil.NewTestLogger(), "")
	if err != nil || result.UsedStaleCache || len(result.Tickets) != 1 {
		t.Fatalf("unexpected refresh result: %+v, %v", result, err)
	}
	if cache.lastFetch.IsZero() || calls.Load() != 1 {
		t.Fatalf("cache was not populated\nlast fetch: %v\ncalls: %d", cache.lastFetch, calls.Load())
	}
}

func TestCacheGet_FailedRefreshUsesPreviousCacheWithoutReplacingIt(t *testing.T) {
	status := atomic.Int32{}
	status.Store(http.StatusOK)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(int(status.Load()))
		if status.Load() == http.StatusOK {
			_, _ = w.Write([]byte(successfulCheckinResponse))
		}
	}))
	defer server.Close()
	cache := newTestCache(server.URL, server.Client())
	logger := testutil.NewTestLogger()
	if _, err := cache.Get(context.Background(), logger, ""); err != nil {
		t.Fatalf("initial refresh failed: %v", err)
	}
	initialFetch, initialTicket := cache.lastFetch, cache.data[0]
	status.Store(http.StatusInternalServerError)
	staleFetch := cache.lastFetch.Add(-cache.ttl)
	cache.lastFetch = staleFetch

	result, err := cache.Get(context.Background(), logger, "")
	if err == nil || !result.UsedStaleCache || len(result.Tickets) != 1 || result.Tickets[0] != initialTicket {
		t.Fatalf("expected stale cached ticket\nresult: %+v\nerror: %v", result, err)
	}
	if cache.lastFetch != staleFetch || !cache.lastFetch.Before(initialFetch) || cache.data[0] != initialTicket {
		t.Fatal("failed refresh replaced successful cache state")
	}
}

func TestCacheGet_FailedRefreshWithoutCacheReturnsNoTickets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadGateway) }))
	defer server.Close()
	result, err := newTestCache(server.URL, server.Client()).Get(context.Background(), testutil.NewTestLogger(), "")
	if err == nil || result.UsedStaleCache || len(result.Tickets) != 0 {
		t.Fatalf("expected no-cache failure\nresult: %+v\nerror: %v", result, err)
	}
}

func TestCacheGet_FailureCooldownMakesOneUpstreamRequest(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	cache := newTestCache(server.URL, server.Client())
	cache.failureCooldown = time.Minute
	for range 3 {
		if _, err := cache.Get(context.Background(), testutil.NewTestLogger(), ""); err == nil {
			t.Fatal("expected unavailable Checkin")
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("upstream calls mismatch\nexpected: 1\nactual: %d", calls.Load())
	}
}

func newTestCache(endpoint string, client *http.Client) *Cache {
	cache := newTicketCache()
	cache.endpoint, cache.client, cache.ttl = endpoint, client, time.Millisecond
	return cache
}
