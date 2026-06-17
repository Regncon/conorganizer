package login

import (
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestLoginForm_RendersDescopeWidgetAndPostLoginRedirect(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en bruker åpner innloggingssiden.",
		When:  "Når innloggingskomponenten rendres.",
		Then:  "Så skal Descope-widgeten være konfigurert og sende vellykket innlogging til post-login.",
	})

	// Given
	expectedWidgetAttributes := map[string]string{
		"project-id": "P2ufzqahlYUHDIprVXtkuCx8MH5C",
		"flow-id":    "sign-up-or-in-passwords-social",
		"theme":      "dark",
	}
	expectedScriptSources := []string{
		"https://descopecdn.com/npm/@descope/web-component@3.21.0/dist/index.js",
		"https://descopecdn.com/npm/@descope/web-js-sdk@1.16.0/dist/index.umd.js",
		"https://static.descope.com/npm/@descope/user-management-widget@0.4.116/dist/index.js",
	}
	expectedInlineScriptParts := []string{
		"descopeSdk.getSessionToken",
		"descopeSdk.getRefreshToken",
		"fetch('/auth/session'",
		"sessionJwt",
		"refreshJwt",
		"window.location.href = '/auth/post-login';",
	}

	// When
	doc := templtest.Render(t, loginForm())
	widget := doc.Find("descope-wc")
	actualScriptSources := collectScriptSources(doc)
	actualInlineScript := doc.Find("script:not([src])").Text()

	// Then
	if widget.Length() != 1 {
		t.Fatalf("expected one Descope widget, got %d", widget.Length())
	}
	for attribute, expectedValue := range expectedWidgetAttributes {
		actualValue, exists := widget.Attr(attribute)
		if !exists || actualValue != expectedValue {
			t.Fatalf("Descope widget attribute %q mismatch\nexpected: %q\nactual:   %q", attribute, expectedValue, actualValue)
		}
	}
	if !slices.Equal(expectedScriptSources, actualScriptSources) {
		t.Fatalf("Descope script sources mismatch\nexpected: %v\nactual:   %v", expectedScriptSources, actualScriptSources)
	}
	for _, expectedInlineScriptPart := range expectedInlineScriptParts {
		if !strings.Contains(actualInlineScript, expectedInlineScriptPart) {
			t.Fatalf("inline script mismatch\nexpected script to contain: %q\nactual script:              %q", expectedInlineScriptPart, actualInlineScript)
		}
	}
}

func collectScriptSources(doc *goquery.Document) []string {
	sources := make([]string, 0)
	doc.Find("script[src]").Each(func(_ int, selection *goquery.Selection) {
		source, exists := selection.Attr("src")
		if exists {
			sources = append(sources, source)
		}
	})
	slices.Sort(sources)
	return sources
}
