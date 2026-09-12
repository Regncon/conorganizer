package requestctx

import (
	"net/url"
	"strings"
)

// SafeProfileReturnURL accepts profile destinations that may be used after login.
// Fragments are intentionally preserved when they are supplied by a browser.
func SafeProfileReturnURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.ContainsAny(raw, "\\\r\n\x00") {
		return ""
	}
	if strings.Contains(strings.ToLower(raw), "%2f") || strings.Contains(strings.ToLower(raw), "%5c") {
		return ""
	}

	destination, err := url.Parse(raw)
	if err != nil || destination.IsAbs() || destination.Host != "" || destination.User != nil {
		return ""
	}
	if destination.Path != "/profile" {
		return ""
	}
	if strings.HasPrefix(destination.Path, "//") || strings.Contains(destination.Path, "\\") {
		return ""
	}
	result := destination.RequestURI()
	if destination.Fragment != "" {
		result += "#" + destination.EscapedFragment()
	}
	return result
}
