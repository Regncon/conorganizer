package eventimgupload

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetUploadedSourceImageURL_ReturnsStampedSourceURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a saved source image with a known modification time.",
		When:  "When the uploaded source image URL is requested.",
		Then:  "Then the returned source URL includes a short timestamp suffix.",
	})

	// Given
	modTime := time.UnixMilli(1_779_123_456_789)
	expectedStamp := strconv.FormatInt(modTime.UnixMilli(), 36)
	expectedURL := "/event-images/event123_source_" + expectedStamp + ".jpg"
	eventImageDir := t.TempDir()
	writeSourceImageFixtureAt(t, eventImageDir, "event123_source.jpg", modTime)

	// When
	actualURL := getUploadedSourceImageURL("event123", &eventImageDir)

	// Then
	if actualURL != expectedURL {
		t.Fatalf("expected source image URL %q, got %q", expectedURL, actualURL)
	}
}

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

func writeSourceImageFixtureAt(t *testing.T, eventImageDir, filename string, modTime time.Time) {
	t.Helper()

	imagePath := filepath.Join(eventImageDir, filename)
	if err := os.WriteFile(imagePath, []byte("source image"), 0o644); err != nil {
		t.Fatalf("failed to create source image fixture %q: %v", imagePath, err)
	}
	if err := os.Chtimes(imagePath, modTime, modTime); err != nil {
		t.Fatalf("failed to set source image fixture modtime for %q: %v", imagePath, err)
	}
}
