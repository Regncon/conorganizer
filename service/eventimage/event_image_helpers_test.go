package eventimage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetEventImageUrl_WhenImageExists_ReturnsVersionedEventImageURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event image file exists on disk.",
		When:  "When the event image URL is built.",
		Then:  "Then the URL includes the image modification time as a cache-busting version.",
	})

	// Given
	eventID := "test-event"
	kind := "card"
	imageDir := t.TempDir()
	filename := eventID + "_" + kind + ".webp"
	imagePath := filepath.Join(imageDir, filename)
	modTime := time.Unix(1_700_000_000, 0)
	expectedURL := "/event-images/test-event_card.webp?v=1700000000000000000"

	if err := os.WriteFile(imagePath, []byte("image"), 0o600); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}
	if err := os.Chtimes(imagePath, modTime, modTime); err != nil {
		t.Fatalf("failed to set test image modtime: %v", err)
	}

	// When
	actualURL := GetEventImageUrl(eventID, kind, &imageDir)

	// Then
	if actualURL != expectedURL {
		t.Fatalf("image URL mismatch\nexpected: %q\nactual:   %q", expectedURL, actualURL)
	}
}

func TestGetEventImageUrl_WhenImageIsMissing_ReturnsPlaceholderURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given no event image file exists on disk.",
		When:  "When the event image URL is built.",
		Then:  "Then the placeholder URL is returned without a cache-busting version.",
	})

	// Given
	eventID := "test-event"
	kind := "banner"
	imageDir := t.TempDir()
	expectedURL := "/static/placeholder_banner.svg"

	// When
	actualURL := GetEventImageUrl(eventID, kind, &imageDir)

	// Then
	if actualURL != expectedURL {
		t.Fatalf("image URL mismatch\nexpected: %q\nactual:   %q", expectedURL, actualURL)
	}
}
