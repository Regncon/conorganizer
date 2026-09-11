package requestctx

import "testing"

func TestSafeProfileReturnURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "profile query and fragment", raw: "/profile?b_id=123&pulje=FredagKveld#mitt-program", want: "/profile?b_id=123&pulje=FredagKveld#mitt-program"},
		{name: "profile subpath", raw: "/profile/tickets?tab=one", want: ""},
		{name: "dot traversal", raw: "/profile/../../auth/logout", want: ""},
		{name: "external host", raw: "//evil.example/profile", want: ""},
		{name: "absolute URL", raw: "https://evil.example/profile", want: ""},
		{name: "encoded slash", raw: "/profile%2f%2fevil.example", want: ""},
		{name: "backslash", raw: "/profile\\\\evil", want: ""},
		{name: "other path", raw: "/admin", want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := SafeProfileReturnURL(test.raw); got != test.want {
				t.Fatalf("SafeProfileReturnURL(%q) = %q, want %q", test.raw, got, test.want)
			}
		})
	}
}
