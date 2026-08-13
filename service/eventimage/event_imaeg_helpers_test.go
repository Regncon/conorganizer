package eventimage

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetEventImageUrlReturnsStampedPublicImageURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a stable public event image file with a known modification time.",
		When:  "When the public image URL is requested.",
		Then:  "Then the returned URL includes a short timestamp suffix derived from the file.",
	})

	// Given
	modTime := time.UnixMilli(1_779_123_456_789)
	expectedStamp := strconv.FormatInt(modTime.UnixMilli(), 36)
	expectedURL := "/event-images/event123_card_" + expectedStamp + ".webp"
	eventImageDir := t.TempDir()
	eventID := "event123"
	kind := "card"
	writeImageFixtureAt(t, eventImageDir, eventID+"_"+kind+".webp", "card image", modTime)

	// When
	actualURL := GetEventImageUrl(eventID, kind, &eventImageDir)

	// Then
	assertString(t, expectedURL, actualURL)
}

func TestGetEventImageUrlReturnsPlaceholderWhenImageIsMissing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given no public event image file for an event.",
		When:  "When the public image URL is requested.",
		Then:  "Then the kind-specific placeholder URL is returned.",
	})

	// Given
	expectedURL := "/static/placeholder_banner.svg"
	eventImageDir := t.TempDir()

	// When
	actualURL := GetEventImageUrl("event123", "banner", &eventImageDir)

	// Then
	assertString(t, expectedURL, actualURL)
}

func TestFileServerServesStampedPublicImageFromStableFile(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a stable public event image file and a stamped public URL for that image.",
		When:  "When the stamped image URL is requested.",
		Then:  "Then the stable image file is served.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedBody := "card image"
	eventImageDir := t.TempDir()
	writeImageFixture(t, eventImageDir, "event123_card.webp", expectedBody)

	// When
	recorder := performImageRequest(eventImageDir, "/event123_card_cachebust.webp")

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatus, expectedBody)
}

func TestFileServerStillServesDirectImageFiles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a direct source image file that is not part of the public stamped image scheme.",
		When:  "When the direct source image URL is requested.",
		Then:  "Then the source image file is served without rewriting the URL.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedBody := "source image"
	eventImageDir := t.TempDir()
	writeImageFixture(t, eventImageDir, "event123_source.jpg", expectedBody)

	// When
	recorder := performImageRequest(eventImageDir, "/event123_source.jpg")

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatus, expectedBody)
}

func TestFileServerServesStampedSourceImageFromStableFile(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a stable source image file and a stamped source URL for that image.",
		When:  "When the stamped source image URL is requested.",
		Then:  "Then the stable source image file is served.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedBody := "source image"
	eventImageDir := t.TempDir()
	writeImageFixture(t, eventImageDir, "event123_source.jpg", expectedBody)

	// When
	recorder := performImageRequest(eventImageDir, "/event123_source_cachebust.jpg")

	// Then
	assertHTTPStatusAndBody(t, recorder, expectedStatus, expectedBody)
}

func writeImageFixture(t *testing.T, eventImageDir, filename, body string) string {
	t.Helper()

	imagePath := filepath.Join(eventImageDir, filename)
	if err := os.WriteFile(imagePath, []byte(body), 0o644); err != nil {
		t.Fatalf("failed to create image fixture %q: %v", imagePath, err)
	}
	return imagePath
}

func writeImageFixtureAt(t *testing.T, eventImageDir, filename, body string, modTime time.Time) {
	t.Helper()

	imagePath := writeImageFixture(t, eventImageDir, filename, body)
	if err := os.Chtimes(imagePath, modTime, modTime); err != nil {
		t.Fatalf("failed to set image fixture modtime for %q: %v", imagePath, err)
	}
}

func performImageRequest(eventImageDir, path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	FileServer(eventImageDir).ServeHTTP(recorder, request)
	return recorder
}

func assertString(t *testing.T, expected, actual string) {
	t.Helper()

	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func assertHTTPStatusAndBody(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	t.Helper()

	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	if recorder.Body.String() != expectedBody {
		t.Fatalf("expected body %q, got %q", expectedBody, recorder.Body.String())
	}
}
