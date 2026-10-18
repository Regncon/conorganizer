package templtest

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var signalsModifierBeforeKey = regexp.MustCompile(`^data-signals__[a-z.]+:`)

func SignalsAttrsWithModifierBeforeKey(doc *goquery.Document) []string {
	var names []string
	doc.Find("*").Each(func(_ int, el *goquery.Selection) {
		for _, attr := range el.Nodes[0].Attr {
			if signalsModifierBeforeKey.MatchString(attr.Key) {
				names = append(names, attr.Key)
			}
		}
	})
	return names
}

func SignalsObject(t testing.TB, el *goquery.Selection, attr string) map[string]any {
	t.Helper()

	raw, ok := el.Attr(attr)
	if !ok {
		t.Fatalf("expected element to have a %s attribute", attr)
	}
	var signals map[string]any
	if err := json.Unmarshal([]byte(raw), &signals); err != nil {
		t.Fatalf("%s is not a JSON object: %q: %v", attr, raw, err)
	}
	return signals
}
