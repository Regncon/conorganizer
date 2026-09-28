package userctx

import (
	"net/http"
	"net/url"
)

// loginHrefWithNeste builds the "/auth" link shown on the unauthenticated
// page, carrying the current request's path and query as the "neste"
// (Norwegian for "next") return target. Login uses this to send the user
// back to the page they tried to reach once they are signed in.
func loginHrefWithNeste(r *http.Request) string {
	neste := r.URL.Path
	if r.URL.RawQuery != "" {
		neste += "?" + r.URL.RawQuery
	}
	if neste == "" || neste == "/" {
		return "/auth"
	}
	return "/auth?neste=" + url.QueryEscape(neste)
}
