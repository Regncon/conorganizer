package eventservice

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

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
