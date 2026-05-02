package templtest

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
)

// Render renders a templ component and parses the output as a goquery document.
func Render(t testing.TB, component templ.Component) *goquery.Document {
	t.Helper()

	var html bytes.Buffer
	if err := component.Render(context.Background(), &html); err != nil {
		t.Fatalf("render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&html)
	if err != nil {
		t.Fatalf("parse component html: %v", err)
	}

	return doc
}

// CollectUniqueInternalHrefs returns unique app-internal hrefs from all links in the document.
func CollectUniqueInternalHrefs(doc *goquery.Document) []string {
	uniqueHrefs := make(map[string]struct{})

	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		href = strings.TrimSpace(href)
		if href == "" || !strings.HasPrefix(href, "/") || strings.HasPrefix(href, "//") {
			return
		}

		uniqueHrefs[href] = struct{}{}
	})

	hrefs := make([]string, 0, len(uniqueHrefs))
	for href := range uniqueHrefs {
		hrefs = append(hrefs, href)
	}

	slices.Sort(hrefs)
	return hrefs
}

// AssertSameHrefs compares href sets without depending on their order.
func AssertSameHrefs(t testing.TB, expectedHrefs []string, actualHrefs []string) {
	t.Helper()

	expectedHrefs = sortedCopy(expectedHrefs)
	actualHrefs = sortedCopy(actualHrefs)

	if !slices.Equal(expectedHrefs, actualHrefs) {
		t.Fatalf("expected internal hrefs %v, got %v", expectedHrefs, actualHrefs)
	}
}

func sortedCopy(values []string) []string {
	copiedValues := append([]string(nil), values...)
	slices.Sort(copiedValues)
	return copiedValues
}
