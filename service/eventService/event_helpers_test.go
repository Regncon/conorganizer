package eventservice

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetEventById_LoadsOpenRegistration(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event configured for open registration.",
		When:  "When the event is loaded by ID.",
		Then:  "Then its open-registration configuration is returned.",
	})

	// Given
	expectedOpenRegistration := true
	db := testutil.CreateTestDB(t, "get_open_registration_event")
	testutil.MustExec(t, db, `
		INSERT INTO events(
			id, title, intro, description, host_name, email, phone_number,
			max_players, is_open_registration
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, "open-event", "Blood in the Clock Tower", "Intro", "Description", "Host", "host@example.com", "12345678", 8, expectedOpenRegistration)

	// When
	event, err := GetEventById("open-event", db)

	// Then
	if err != nil {
		t.Fatalf("expected event lookup to succeed: %v", err)
	}
	if event == nil {
		t.Fatal("expected event to be returned")
	}
	if event.IsOpenRegistration != expectedOpenRegistration {
		t.Fatalf("open-registration mismatch\nexpected: %t\nactual:   %t", expectedOpenRegistration, event.IsOpenRegistration)
	}
}

func TestGetEventById_OpenRegistrationDefaultsToFalse(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event without explicit open-registration configuration.",
		When:  "When the event is loaded by ID.",
		Then:  "Then open registration is disabled by default.",
	})

	// Given
	expectedOpenRegistration := false
	db := testutil.CreateTestDB(t, "get_regular_event")
	testutil.MustExec(t, db, `
		INSERT INTO events(
			id, title, intro, description, host_name, email, phone_number, max_players
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, "regular-event", "Call of Cthulhu", "Intro", "Description", "Host", "host@example.com", "12345678", 5)

	// When
	event, err := GetEventById("regular-event", db)

	// Then
	if err != nil {
		t.Fatalf("expected event lookup to succeed: %v", err)
	}
	if event == nil {
		t.Fatal("expected event to be returned")
	}
	if event.IsOpenRegistration != expectedOpenRegistration {
		t.Fatalf("open-registration mismatch\nexpected: %t\nactual:   %t", expectedOpenRegistration, event.IsOpenRegistration)
	}
}

func TestSanitizeMdToHTML_RendersSafeMarkdownAndRawHTML(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given markdown with safe formatting, links, and raw HTML.",
		When:  "When the markdown is rendered and sanitized.",
		Then:  "Then safe markdown and safe raw HTML remain in the output.",
	})

	// Given
	expectedFragments := []string{
		`<h1 id="event-title">Event Title</h1>`,
		`<strong>bold</strong>`,
		`<a href="https://example.com"`,
		`>event page</a>`,
		`<div>Safe raw HTML</div>`,
		`<img src="https://example.com/banner.png" alt="Banner">`,
	}
	md := []byte(`# Event Title

This is **bold** text with an [event page](https://example.com).

<div>Safe raw HTML</div>

<img src="https://example.com/banner.png" alt="Banner">
`)

	// When
	actual := string(SanitizeMdToHTML(md))

	// Then
	assertStringContainsAll(t, actual, expectedFragments)
}

func TestSanitizeMdToHTML_RemovesExecutableContent(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given markdown with scripts, unsafe links, and unsafe raw HTML attributes.",
		When:  "When the markdown is rendered and sanitized.",
		Then:  "Then executable content is removed while safe content remains.",
	})

	// Given
	expectedSafeFragments := []string{
		`Intro text`,
		`<a href="https://example.com"`,
		`>safe link</a>`,
		`bad link`,
		`<img src="https://example.com/banner.png" alt="Banner">`,
	}
	forbiddenFragments := []string{
		`<script`,
		`</script`,
		`javascript:`,
		`onerror`,
		`stealCookies`,
	}
	md := []byte(`Intro text

<script>stealCookies()</script>

[safe link](https://example.com) [bad link](javascript:stealCookies)

<img src="https://example.com/banner.png" alt="Banner" onerror="stealCookies()">
`)

	// When
	actual := string(SanitizeMdToHTML(md))

	// Then
	assertStringContainsAll(t, actual, expectedSafeFragments)
	assertStringExcludesAll(t, actual, forbiddenFragments)
}

func assertStringContainsAll(t testing.TB, actual string, expectedFragments []string) {
	t.Helper()

	for _, expected := range expectedFragments {
		if !strings.Contains(actual, expected) {
			t.Fatalf("expected sanitized HTML to contain %q\nactual: %s", expected, actual)
		}
	}
}

func assertStringExcludesAll(t testing.TB, actual string, forbiddenFragments []string) {
	t.Helper()

	lowerActual := strings.ToLower(actual)
	for _, forbidden := range forbiddenFragments {
		if strings.Contains(lowerActual, strings.ToLower(forbidden)) {
			t.Fatalf("expected sanitized HTML not to contain %q\nactual: %s", forbidden, actual)
		}
	}
}
