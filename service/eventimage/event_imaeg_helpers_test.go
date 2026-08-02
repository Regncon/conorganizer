package eventimage

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestGetEventImageUrlReturnsStampedPublicImageURL(t *testing.T) {
	eventImageDir := t.TempDir()
	eventID := "event123"
	kind := "card"
	imagePath := filepath.Join(eventImageDir, eventID+"_"+kind+".webp")
	if err := os.WriteFile(imagePath, []byte("card image"), 0o644); err != nil {
		t.Fatalf("failed to create image fixture: %v", err)
	}

	modTime := time.UnixMilli(1_779_123_456_789)
	if err := os.Chtimes(imagePath, modTime, modTime); err != nil {
		t.Fatalf("failed to set image fixture modtime: %v", err)
	}

	expectedStamp := strconv.FormatInt(modTime.UnixMilli(), 36)
	expectedURL := "/event-images/event123_card_" + expectedStamp + ".webp"
	actualURL := GetEventImageUrl(eventID, kind, &eventImageDir)

	if actualURL != expectedURL {
		t.Fatalf("expected %q, got %q", expectedURL, actualURL)
	}
}

func TestGetEventImageUrlReturnsPlaceholderWhenImageIsMissing(t *testing.T) {
	eventImageDir := t.TempDir()
	expectedURL := "/static/placeholder_banner.svg"
	actualURL := GetEventImageUrl("event123", "banner", &eventImageDir)

	if actualURL != expectedURL {
		t.Fatalf("expected %q, got %q", expectedURL, actualURL)
	}
}

func TestFileServerServesStampedPublicImageFromStableFile(t *testing.T) {
	eventImageDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(eventImageDir, "event123_card.webp"), []byte("card image"), 0o644); err != nil {
		t.Fatalf("failed to create image fixture: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/event123_card_cachebust.webp", nil)
	FileServer(eventImageDir).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "card image" {
		t.Fatalf("expected stable file contents, got %q", recorder.Body.String())
	}
}

func TestFileServerStillServesDirectImageFiles(t *testing.T) {
	eventImageDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(eventImageDir, "event123_source.jpg"), []byte("source image"), 0o644); err != nil {
		t.Fatalf("failed to create source image fixture: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/event123_source.jpg", nil)
	FileServer(eventImageDir).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != "source image" {
		t.Fatalf("expected source file contents, got %q", recorder.Body.String())
	}
}
