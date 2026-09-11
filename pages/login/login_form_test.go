package login

import (
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/Regncon/conorganizer/service/authctx"
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
		"project-id": authctx.DescopeProjectID,
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
	if strings.Contains(actualInlineScript, "console.log('Email:'") || strings.Contains(actualInlineScript, "console.log('User:'") {
		t.Fatal("login script must not write user identity details to the browser console")
	}
}

func TestLoginForm_PreservesSafeProfileReturnURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at innlogging ble startet fra en profil-lenke.",
		When:  "Når innloggingskomponenten rendres.",
		Then:  "Så skal den trygge profilsiden følge med til post-login.",
	})

	// Given
	expectedReturnURL := "/profile?b_id=123&pulje=FredagKveld#mitt-program"

	// When
	doc := templtest.Render(t, loginFormForReturnURL(expectedReturnURL))
	wrapper := doc.Find("[data-login-return-url]")
	actualReturnURL, exists := wrapper.Attr("data-login-return-url")
	actualScript := doc.Find("script:not([src])").Text()

	// Then
	if !exists || actualReturnURL != expectedReturnURL {
		t.Fatalf("login return URL mismatch\nexpected: %q\nactual:   %q", expectedReturnURL, actualReturnURL)
	}
	if !strings.Contains(actualScript, "encodeURIComponent(returnURL)") {
		t.Fatalf("expected post-login script to preserve the return URL, got %q", actualScript)
	}
}

func TestAlreadyLoggedIn_PreservesProfileReturnURLInFallbackLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker følger en lenke til festivalprogrammet.",
		When:  "Når siden for allerede innlogget bruker rendres.",
		Then:  "Så skal både automatisk videresending og fallback-lenken gå til festivalprogrammet.",
	})

	// Given
	expectedReturnURL := "/profile?pulje=FredagKveld#mitt-program"

	// When
	doc := templtest.Render(t, alreadyLogedInForReturnURL(expectedReturnURL))
	link := doc.Find("[data-login-return-link]")
	actualHref, hrefExists := link.Attr("href")

	// Then
	if !hrefExists || actualHref != expectedReturnURL {
		t.Fatalf("already logged-in fallback URL mismatch\nexpected: %q\nactual:   %q", expectedReturnURL, actualHref)
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
