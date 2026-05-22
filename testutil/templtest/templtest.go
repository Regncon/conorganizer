package templtest

import (
	"bytes"
	"context"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"golang.org/x/net/html"
)

func Render(t *testing.T, component templ.Component) *html.Node {
	t.Helper()

	var buf bytes.Buffer
	if err := component.Render(context.Background(), &buf); err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := html.Parse(&buf)
	if err != nil {
		t.Fatalf("failed to parse rendered component HTML: %v", err)
	}

	return doc
}

func HasSelector(doc *html.Node, selector string) bool {
	return findFirst(doc, selector) != nil
}

func CollectTexts(doc *html.Node, selector string) []string {
	matches := collect(doc, selector)
	texts := make([]string, 0, len(matches))
	for _, match := range matches {
		texts = append(texts, normalizeText(textContent(match)))
	}

	return texts
}

func CollectUniqueHrefs(doc *html.Node) []string {
	links := collect(doc, "a")
	seen := make(map[string]struct{}, len(links))
	hrefs := make([]string, 0, len(links))

	for _, link := range links {
		href, ok := attr(link, "href")
		if !ok {
			continue
		}
		if _, exists := seen[href]; exists {
			continue
		}
		seen[href] = struct{}{}
		hrefs = append(hrefs, href)
	}

	return hrefs
}

func AssertSameHrefs(t *testing.T, expected []string, actual []string) {
	t.Helper()

	expectedCopy := slices.Clone(expected)
	actualCopy := slices.Clone(actual)
	slices.Sort(expectedCopy)
	slices.Sort(actualCopy)

	if !slices.Equal(expectedCopy, actualCopy) {
		t.Fatalf("hrefs mismatch\nexpected: %v\nactual:   %v", expectedCopy, actualCopy)
	}
}

func collect(node *html.Node, selector string) []*html.Node {
	var matches []*html.Node
	walk(node, func(candidate *html.Node) {
		if matchesSelector(candidate, selector) {
			matches = append(matches, candidate)
		}
	})

	return matches
}

func findFirst(node *html.Node, selector string) *html.Node {
	var match *html.Node
	walk(node, func(candidate *html.Node) {
		if match != nil {
			return
		}
		if matchesSelector(candidate, selector) {
			match = candidate
		}
	})

	return match
}

func walk(node *html.Node, visit func(*html.Node)) {
	if node == nil {
		return
	}

	visit(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walk(child, visit)
	}
}

func matchesSelector(node *html.Node, selector string) bool {
	if node.Type != html.ElementNode {
		return false
	}

	if strings.HasPrefix(selector, ".") {
		return hasClass(node, strings.TrimPrefix(selector, "."))
	}

	return node.Data == selector
}

func hasClass(node *html.Node, className string) bool {
	value, ok := attr(node, "class")
	if !ok {
		return false
	}

	for _, class := range strings.Fields(value) {
		if class == className {
			return true
		}
	}

	return false
}

func attr(node *html.Node, key string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func textContent(node *html.Node) string {
	var builder strings.Builder
	writeText(&builder, node)
	return builder.String()
}

func writeText(writer io.StringWriter, node *html.Node) {
	if node == nil {
		return
	}

	if node.Type == html.TextNode {
		_, _ = writer.WriteString(node.Data)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		writeText(writer, child)
	}
}

func normalizeText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
