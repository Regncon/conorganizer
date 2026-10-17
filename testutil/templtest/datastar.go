package templtest

import (
	"regexp"

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
