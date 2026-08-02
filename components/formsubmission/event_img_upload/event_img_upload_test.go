package eventimgupload

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestBroadcastEventImageUpdate_BroadcastsEventBucket(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a cropped event image update and a live broadcaster.",
		When:  "When the image update is broadcast.",
		Then:  "Then the events live bucket is notified.",
	})

	// Given
	expectedBuckets := []live.Bucket{live.BucketEvents}
	broadcaster := &recordingEventImageBroadcaster{}

	// When
	err := broadcastEventImageUpdate(context.Background(), broadcaster)

	// Then
	if err != nil {
		t.Fatalf("expected broadcast to succeed: %v", err)
	}
	assertBroadcastBuckets(t, expectedBuckets, broadcaster.buckets)
}

func TestBroadcastEventImageUpdate_ReturnsBroadcastError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a live broadcaster that cannot notify clients.",
		When:  "When the image update is broadcast.",
		Then:  "Then the broadcast error is returned.",
	})

	// Given
	expectedErr := errors.New("broadcast failed")
	broadcaster := &recordingEventImageBroadcaster{err: expectedErr}

	// When
	err := broadcastEventImageUpdate(context.Background(), broadcaster)

	// Then
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

type recordingEventImageBroadcaster struct {
	buckets []live.Bucket
	err     error
}

func (b *recordingEventImageBroadcaster) Broadcast(_ context.Context, buckets ...live.Bucket) error {
	b.buckets = slices.Clone(buckets)
	return b.err
}

func assertBroadcastBuckets(t *testing.T, expected, actual []live.Bucket) {
	t.Helper()

	if !slices.Equal(actual, expected) {
		t.Fatalf("expected broadcast buckets %v, got %v", expected, actual)
	}
}
