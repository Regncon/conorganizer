package formsubmission

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

// The placement endpoints answer an under-18-in-an-18+-game attempt with a
// warning instead of a seat, so the page has to carry the dialog that turns that
// answer into a deliberate confirmation.
func TestApprovalAgeWarningDialog_RendersConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt aldersgrense-dialogen for godkjenningssida.",
		When:  "Når han blir rendra.",
		Then:  "Så finst dialogen som stadfestar plassering av nokon under 18 i eit 18+-spel.",
	})

	doc := templtest.Render(t, ApprovalAgeWarningDialog())
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}

	if got := doc.Find("#approval-age-dialog").Length(); got != 1 {
		t.Fatalf("expected exactly one age confirmation dialog, got %d", got)
	}
	if !strings.Contains(html, "assignmentAgeConfirmed") {
		t.Errorf("the confirm button should re-post with assignmentAgeConfirmed set")
	}
	if !strings.Contains(html, "$ageWarningUrl") {
		t.Errorf("the confirm button should re-post to the endpoint the warning came from")
	}
	if !strings.Contains(html, "Plasser likevel") {
		t.Errorf("expected a Norwegian confirm label, got: %s", html)
	}

	// The dialog lives outside the live region, so it must own its signal
	// defaults instead of relying on the streamed content to declare them.
	signals := doc.Find("#approval-age-dialog").AttrOr("data-signals", "")
	for _, name := range []string{"ageWarningText", "ageWarningUrl", "ageWarningMethod", "assignmentAgeConfirmed"} {
		if !strings.Contains(signals, name) {
			t.Errorf("dialog data-signals should initialise %s; got %q", name, signals)
		}
	}
}

// datastar 1.0.x compiles every expression with a plain Function(), so `await`
// is a SyntaxError and the handler silently never runs. The confirmation flag is
// page-global and sent with every request, so it must be reset on the statement
// right after the action call — datastar serialises the payload synchronously —
// never later in a promise callback, which would pre-confirm any placement
// fired while the re-post is in flight. The dialog must also clear its warning
// on close so an identical warning can re-open it.
func TestApprovalAgeWarningDialog_ConfirmResetsFlagSynchronouslyAndClearsOnClose(t *testing.T) {
	doc := templtest.Render(t, ApprovalAgeWarningDialog())
	dialog := doc.Find("#approval-age-dialog")

	confirm := dialog.Find("button").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(s.AttrOr("data-on:click", ""), "assignmentAgeConfirmed = true")
	})
	if confirm.Length() != 1 {
		t.Fatalf("expected exactly one confirm button setting assignmentAgeConfirmed, got %d", confirm.Length())
	}
	action := confirm.AttrOr("data-on:click", "")
	if strings.Contains(action, "await") {
		t.Errorf("confirm action must not use await (datastar compiles with Function, not AsyncFunction); got %q", action)
	}
	if strings.Contains(action, ".finally(") || strings.Contains(action, ".then(") {
		t.Errorf("confirm action must reset the flag synchronously, not in a promise callback; got %q", action)
	}
	post := strings.Index(action, "@post(")
	reset := strings.Index(action, "$assignmentAgeConfirmed = false")
	if post < 0 || reset < 0 || reset < post {
		t.Errorf("confirm action should reset assignmentAgeConfirmed right after dispatching the request; got %q", action)
	}
	if idReset := strings.Index(action, "$assignmentBillettholderId = 0"); idReset < 0 || idReset < post {
		t.Errorf("confirm action should reset the picker id to 0 after dispatching, like the picker buttons do; got %q", action)
	}

	onClose := dialog.AttrOr("data-on:close", "")
	if !strings.Contains(onClose, "$ageWarningText = ''") {
		t.Errorf("age dialog should clear $ageWarningText on close; got %q", onClose)
	}
}
