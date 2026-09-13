package eventimgupload

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestEventImageFormSubmission_ProgressUploadReturnsPersistedSourceURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing event and an authenticated multipart source image upload.",
		When:  "When the upload action requests progress responses and Datastar handling.",
		Then:  "Then the image is saved and its cache-busted URL is returned as JSON without redirecting.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedImage := []byte("source image content")
	expectedPath := "/event-images/event-123_source.png"
	db, router, imageDir := newEventImageUploadFixture(t)
	request := newImageUploadRequest(t, "/profile/api/new/event-123/upload", "Poster.PNG", expectedImage, "")
	request.Header.Set("Upload-Progress-Request", "true")
	request.Header.Set("Datastar-Request", "true")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	assertUploadStatus(t, recorder, expectedStatus)
	assertNoUploadRedirect(t, recorder)
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected JSON response, got content type %q", got)
	}
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode source upload response: %v", err)
	}
	assertSourceImageURL(t, response["sourceImageUrl"], expectedPath)
	assertUploadedImage(t, imageDir, "event-123_source.png", expectedImage)
	assertEventImageAudit(t, db)
}

func TestEventImageFormSubmission_DatastarUploadPatchesSourceURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing event and an authenticated multipart source image upload.",
		When:  "When only the Datastar response header is requested.",
		Then:  "Then the saved image URL is patched through an event stream without redirecting.",
	})

	// Given
	expectedStatus := http.StatusOK
	expectedImage := []byte("source image content")
	expectedPath := "/event-images/event-123_source.png"
	_, router, imageDir := newEventImageUploadFixture(t)
	request := newImageUploadRequest(t, "/profile/api/new/event-123/upload", "source.png", expectedImage, "")
	request.Header.Set("Datastar-Request", "true")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	assertUploadStatus(t, recorder, expectedStatus)
	assertNoUploadRedirect(t, recorder)
	if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/event-stream") {
		t.Fatalf("expected event stream, got content type %q", got)
	}
	if !strings.Contains(recorder.Body.String(), "event: datastar-patch-signals\n") {
		t.Fatalf("expected a signal patch event, got %q", recorder.Body.String())
	}
	var signals map[string]string
	for line := range strings.SplitSeq(recorder.Body.String(), "\n") {
		if encoded, ok := strings.CutPrefix(line, "data: signals "); ok {
			if err := json.Unmarshal([]byte(encoded), &signals); err != nil {
				t.Fatalf("decode source image signal patch: %v", err)
			}
		}
	}
	assertSourceImageURL(t, signals["sourceImageUrl"], expectedPath)
	assertUploadedImage(t, imageDir, "event-123_source.png", expectedImage)
}

func TestEventImageCroppedSubmission_DatastarSavesCropWithoutRedirect(t *testing.T) {
	for _, kind := range []string{"card", "banner"} {
		t.Run(kind, func(t *testing.T) {
			bdd.Behavior(t, bdd.BDD{
				Given: "Given an existing event and an authenticated cropped " + kind + " image.",
				When:  "When Datastar submits the crop as multipart form data.",
				Then:  "Then the crop and audit update are saved with an empty 204 response and no redirect.",
			})

			// Given
			expectedStatus := http.StatusNoContent
			expectedImage := []byte("cropped " + kind + " image content")
			db, router, imageDir := newEventImageUploadFixture(t)
			request := newImageUploadRequest(t, "/profile/api/new/event-123/upload-cropped", "crop.webp", expectedImage, kind)
			request.Header.Set("Datastar-Request", "true")
			recorder := httptest.NewRecorder()

			// When
			router.ServeHTTP(recorder, request)

			// Then
			assertUploadStatus(t, recorder, expectedStatus)
			assertNoUploadRedirect(t, recorder)
			if recorder.Body.Len() != 0 {
				t.Fatalf("expected empty crop success response, got %q", recorder.Body.String())
			}
			assertUploadedImage(t, imageDir, "event-123_"+kind+".webp", expectedImage)
			assertEventImageAudit(t, db)
		})
	}
}

func TestEventImageCroppedSubmission_OrdinaryFormRedirectsAfterSaving(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing event and an authenticated cropped image.",
		When:  "When an ordinary multipart form submits the crop.",
		Then:  "Then the image is saved and the browser is redirected to the image page.",
	})

	// Given
	expectedStatus := http.StatusSeeOther
	expectedLocation := "/profile/new/event-123/image/"
	expectedImage := []byte("cropped card image content")
	_, router, imageDir := newEventImageUploadFixture(t)
	request := newImageUploadRequest(t, "/profile/api/new/event-123/upload-cropped", "crop.webp", expectedImage, "card")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	assertUploadStatus(t, recorder, expectedStatus)
	if got := recorder.Header().Get("Location"); got != expectedLocation {
		t.Fatalf("expected redirect to %q, got %q", expectedLocation, got)
	}
	assertUploadedImage(t, imageDir, "event-123_card.webp", expectedImage)
}

func TestEventImageCroppedSubmission_FailedDatastarUploadDoesNotReportSuccess(t *testing.T) {
	for _, scenario := range []struct {
		name           string
		kind           string
		filename       string
		expectedStatus int
		prepare        func(*testing.T, *sql.DB, string)
	}{
		{name: "invalid image kind", kind: "source", filename: "crop.webp", expectedStatus: http.StatusBadRequest},
		{name: "missing image", kind: "card", expectedStatus: http.StatusBadRequest},
		{
			name: "image cannot be saved", kind: "card", filename: "crop.webp", expectedStatus: http.StatusInternalServerError,
			prepare: func(t *testing.T, _ *sql.DB, imageDir string) {
				if err := os.Mkdir(filepath.Join(imageDir, "event-123_card.webp"), 0700); err != nil {
					t.Fatalf("block crop destination with a directory: %v", err)
				}
			},
		},
		{
			name: "audit cannot be saved", kind: "card", filename: "crop.webp", expectedStatus: http.StatusInternalServerError,
			prepare: func(t *testing.T, db *sql.DB, _ string) {
				testutil.MustExec(t, db, `CREATE TRIGGER reject_image_audit BEFORE UPDATE ON events BEGIN SELECT RAISE(FAIL, 'audit unavailable'); END`)
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			bdd.Behavior(t, bdd.BDD{
				Given: "Given a cropped image request where " + scenario.name + ".",
				When:  "When Datastar submits the image.",
				Then:  "Then the handler returns an error response without redirecting or updating the audit record.",
			})

			// Given
			expectedStatus := scenario.expectedStatus
			db, router, imageDir := newEventImageUploadFixture(t)
			if scenario.prepare != nil {
				scenario.prepare(t, db, imageDir)
			}
			request := newImageUploadRequest(t, "/profile/api/new/event-123/upload-cropped", scenario.filename, []byte("cropped image"), scenario.kind)
			request.Header.Set("Datastar-Request", "true")
			recorder := httptest.NewRecorder()

			// When
			router.ServeHTTP(recorder, request)

			// Then
			assertUploadStatus(t, recorder, expectedStatus)
			assertNoUploadRedirect(t, recorder)
			if recorder.Body.Len() == 0 {
				t.Fatal("expected an error message in the failed upload response")
			}
			if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM events WHERE id = 'event-123' AND updated_by_id IS NOT NULL`); got != 0 {
				t.Fatal("expected failed upload to leave the event audit unchanged")
			}
		})
	}
}

func newEventImageUploadFixture(t *testing.T) (*sql.DB, http.Handler, string) {
	t.Helper()
	db, logger := testutil.CreateTestDBAndLogger(t, "event_image_upload")
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (42, 'image-uploader', 'uploader@example.com')`)
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players, updated_at)
		VALUES ('event-123', 'Image upload event', '', '', 'Host', 'uploader@example.com', '', 4, '2000-01-01T00:00:00.000Z')`)
	imageDir := t.TempDir()
	router := chi.NewRouter()
	router.Route("/profile/api/new/{id}/upload", func(uploadRouter chi.Router) {
		EventImageFormSubmission(uploadRouter, db, &imageDir, logger)
	})
	router.Route("/profile/api/new/{id}/upload-cropped", func(uploadRouter chi.Router) {
		EventImageCroppedSubmission(uploadRouter, db, &imageDir, logger)
	})
	return db, router, imageDir
}

func newImageUploadRequest(t *testing.T, path, filename string, content []byte, kind string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if filename != "" {
		part, err := writer.CreateFormFile("image", filename)
		if err != nil {
			t.Fatalf("create multipart image field: %v", err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatalf("write multipart image: %v", err)
		}
	}
	if kind != "" {
		if err := writer.WriteField("kind", kind); err != nil {
			t.Fatalf("write image kind: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart form: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, path, &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request.WithContext(authctx.WithUserToken(request.Context(), "image-uploader", "uploader@example.com"))
}

func assertUploadStatus(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if recorder.Code != expected {
		t.Fatalf("expected upload status %d, got %d: %s", expected, recorder.Code, recorder.Body.String())
	}
}

func assertNoUploadRedirect(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()
	if location := recorder.Header().Get("Location"); location != "" {
		t.Fatalf("expected no redirect, got Location %q", location)
	}
}

func assertSourceImageURL(t *testing.T, rawURL, expectedPath string) {
	t.Helper()
	imageURL, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse source image URL: %v", err)
	}
	if imageURL.Path != expectedPath {
		t.Fatalf("expected source image path %q, got %q", expectedPath, imageURL.Path)
	}
	if version, err := strconv.ParseInt(imageURL.Query().Get("v"), 10, 64); err != nil || version <= 0 {
		t.Fatalf("expected a positive cache version in source image URL, got %q", rawURL)
	}
}

func assertUploadedImage(t *testing.T, imageDir, filename string, expected []byte) {
	t.Helper()
	actual, err := os.ReadFile(filepath.Join(imageDir, filename))
	if err != nil {
		t.Fatalf("read persisted image: %v", err)
	}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("expected persisted image %q, got %q", expected, actual)
	}
}

func assertEventImageAudit(t *testing.T, db *sql.DB) {
	t.Helper()
	var updatedBy int
	var updatedAt string
	if err := db.QueryRow(`SELECT updated_by_id, updated_at FROM events WHERE id = 'event-123'`).Scan(&updatedBy, &updatedAt); err != nil {
		t.Fatalf("read image upload audit: %v", err)
	}
	if updatedBy != 42 || updatedAt == "" || updatedAt == "2000-01-01T00:00:00.000Z" {
		t.Fatalf("expected current uploader and refreshed timestamp, got updated_by_id=%d, updated_at=%q", updatedBy, updatedAt)
	}
}
