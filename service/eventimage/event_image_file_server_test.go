package eventimage

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestFileServer_WhenEventImageExists_SetsETagHeader(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event image file exists on disk.",
		When:  "When the image is requested through the event image file server.",
		Then:  "Then the response includes an ETag based on the file metadata.",
	})

	// Given
	imageDir := t.TempDir()
	imageName := "test-event_card.webp"
	imageBody := []byte("image")
	imagePath := filepath.Join(imageDir, imageName)
	modTime := time.Unix(1_700_000_000, 0)

	if err := os.WriteFile(imagePath, imageBody, 0o600); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}
	if err := os.Chtimes(imagePath, modTime, modTime); err != nil {
		t.Fatalf("failed to set test image modtime: %v", err)
	}

	handler := FileServer(imageDir)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/"+imageName, nil)

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("status mismatch\nexpected: %d\nactual:   %d", http.StatusOK, recorder.Code)
	}

	expectedETag := eventImageETag(modTime, int64(len(imageBody)))
	actualETag := recorder.Header().Get("ETag")
	if actualETag != expectedETag {
		t.Fatalf("ETag mismatch\nexpected: %q\nactual:   %q", expectedETag, actualETag)
	}
}

func TestFileServer_WhenIfNoneMatchMatchesETag_ReturnsNotModified(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the browser already has the current event image ETag.",
		When:  "When the image is requested with If-None-Match.",
		Then:  "Then the server returns 304 Not Modified without resending the image body.",
	})

	// Given
	imageDir := t.TempDir()
	imageName := "test-event_banner.webp"
	imageBody := []byte("banner image")
	imagePath := filepath.Join(imageDir, imageName)
	modTime := time.Unix(1_700_000_000, 0)

	if err := os.WriteFile(imagePath, imageBody, 0o600); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}
	if err := os.Chtimes(imagePath, modTime, modTime); err != nil {
		t.Fatalf("failed to set test image modtime: %v", err)
	}

	handler := FileServer(imageDir)
	request := httptest.NewRequest(http.MethodGet, "/"+imageName, nil)
	request.Header.Set("If-None-Match", eventImageETag(modTime, int64(len(imageBody))))
	recorder := httptest.NewRecorder()

	// When
	handler.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusNotModified {
		t.Fatalf("status mismatch\nexpected: %d\nactual:   %d", http.StatusNotModified, recorder.Code)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("expected empty response body for 304, got %d bytes", recorder.Body.Len())
	}
}
