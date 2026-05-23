package templtest

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HasSelector(doc *goquery.Document, selector string) bool {
	return doc.Find(selector).Length() > 0
}

func CollectTexts(doc *goquery.Document, selector string) []string {
	matches := doc.Find(selector)
	texts := make([]string, 0, matches.Length())

	matches.Each(func(_ int, selection *goquery.Selection) {
		texts = append(texts, normalizeText(selection.Text()))
	})

	return texts
}

func normalizeText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
