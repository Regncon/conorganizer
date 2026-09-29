package admin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	billettholderadmin "github.com/Regncon/conorganizer/pages/admin/billettholder_admin"
	printfriendly "github.com/Regncon/conorganizer/pages/print-friendly"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/program"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

// This file tests the admin landing page against the shared design contract
// (breadcrumbs unchanged; three sections in a fixed order; a fixed, ordered
// set of tool cards per section; the Publiser program card restyled but with
// its publish toggle behaviour preserved). It intentionally does not assume
// anything about the current implementation beyond that contract.

const insertPublishingStateSQL = `
	INSERT INTO program_publishing_state(id, is_published)
	VALUES(1, ?)
	ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
`

// expectedTool describes one admin-tool-card. A tool with an empty href is a
// non-link card (only "publiser" today) and is asserted separately, since its
// shape (wide, contains the print link and the publish toggle) differs from
// an ordinary link card.
type expectedTool struct {
	slug        string
	title       string
	href        string
	graphic     string
	description string
}

type expectedSection struct {
	slug    string
	heading string
	tools   []expectedTool
}

var adminContractSections = []expectedSection{
	{
		slug:    "puljefordeling",
		heading: "Puljefordeling",
		tools: []expectedTool{
			{
				slug:        "fordeling",
				title:       "Fordeling av deltakere",
				href:        "/admin/puljefordeling/",
				graphic:     "/static/admin/algorithm.svg",
				description: "Fordel festivaldeltakerne til arrangementene som skjer i løpet av festivalen.",
			},
			{
				slug:        "puljeoppsett",
				title:       "Sett arrangementer i puljer",
				href:        "/admin/puljeoppsett/",
				graphic:     "/static/admin/blueprint.svg",
				description: "Organiser spill inn i puljene som finner sted under festivalen.",
			},
		},
	},
	{
		slug:    "programoppsett",
		heading: "Programoppsett",
		tools: []expectedTool{
			{slug: "publiser"}, // non-link wide card, asserted separately below
			{
				slug:        "rom",
				title:       "Administrer rom",
				href:        "/admin/rooms/",
				graphic:     "/static/admin/blueprint.svg",
				description: "Legg til, rediger og få en oversikt over alle tilgjengelige rom under festivalen.",
			},
			{
				slug:        "godkjenn",
				title:       "Godkjenn arrangementer",
				href:        "/admin/approval/",
				graphic:     "/static/admin/checklist.svg",
				description: "Se alle arrangementer som er sendt inn av deltakerne og som venter på godkjenning.",
			},
		},
	},
	{
		slug:    "ovrige",
		heading: "Øvrige verktøy",
		tools: []expectedTool{
			{
				slug:        "billettholdere",
				title:       "Administrer billettholdere",
				href:        "/admin/billettholder/",
				graphic:     "/static/admin/access-badge.svg",
				description: "Se en oversikt over alle billettholdere som kommer til festivalen.",
			},
			{
				slug:        "tilbakemeldinger",
				title:       "Tilbakemeldinger",
				href:        "/admin/tilbakemeldinger/",
				graphic:     "/static/admin/inbox.svg",
				description: "Les tilbakemeldinger fra deltakerne om nettsiden og festivalen.",
			},
		},
	},
}

const publiserGraphic = "/static/admin/launch.svg"
const publiserFrameGraphic = "/static/admin/frame.svg"

func collectCardSlugs(t *testing.T, section *goquery.Selection) []string {
	t.Helper()
	var slugs []string
	section.Find("[data-admin-tool]").Each(func(_ int, card *goquery.Selection) {
		slug, _ := card.Attr("data-admin-tool")
		slugs = append(slugs, slug)
	})
	return slugs
}

func TestAdminPage_RendersSectionsAndCardsInContractOrder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt adminforsidens kontrakt for overskrift, seksjoner og verktøykort.",
		When:  "Når adminforsiden rendres.",
		Then:  "Så vises overskriften, de tre seksjonene i riktig rekkefølge, og hvert korts lenke, tittel, beskrivelse og grafikk stemmer.",
	})

	// Given
	expectedHeading := "Adminverktøy"
	db := testutil.CreateTestDB(t, "admin_landing_structure")
	testutil.MustExec(t, db, insertPublishingStateSQL, 0)

	// When
	doc := templtest.Render(t, adminPage(db))

	// Then
	actualHeading := strings.TrimSpace(doc.Find("h1.page-heading").Text())
	if actualHeading != expectedHeading {
		t.Fatalf("h1 text = %q, want %q", actualHeading, expectedHeading)
	}

	sections := doc.Find("section.admin-section")
	if sections.Length() != len(adminContractSections) {
		t.Fatalf("expected %d admin sections, got %d", len(adminContractSections), sections.Length())
	}

	sections.Each(func(i int, section *goquery.Selection) {
		expected := adminContractSections[i]

		actualSlug, _ := section.Attr("data-admin-section")
		if actualSlug != expected.slug {
			t.Fatalf("section %d data-admin-section = %q, want %q", i, actualSlug, expected.slug)
		}

		heading := section.Find("h2")
		if heading.Length() != 1 {
			t.Fatalf("section %q: expected exactly one h2, got %d", expected.slug, heading.Length())
		}
		actualHeadingText := strings.TrimSpace(heading.Text())
		if actualHeadingText != expected.heading {
			t.Fatalf("section %q h2 text = %q, want %q", expected.slug, actualHeadingText, expected.heading)
		}
		headingID, _ := heading.Attr("id")
		labelledBy, _ := section.Attr("aria-labelledby")
		if headingID == "" || labelledBy != headingID {
			t.Fatalf("section %q: aria-labelledby %q does not point at the h2's id %q", expected.slug, labelledBy, headingID)
		}

		expectedSlugs := make([]string, len(expected.tools))
		for j, tool := range expected.tools {
			expectedSlugs[j] = tool.slug
		}
		actualSlugs := collectCardSlugs(t, section)
		if !slices.Equal(actualSlugs, expectedSlugs) {
			t.Fatalf("section %q tool order = %v, want %v", expected.slug, actualSlugs, expectedSlugs)
		}

		for _, tool := range expected.tools {
			if tool.href == "" {
				// Non-link cards (the "publiser" wide card) are asserted separately.
				continue
			}
			card := section.Find(fmt.Sprintf(`[data-admin-tool="%s"]`, tool.slug))
			if card.Length() != 1 {
				t.Fatalf("expected exactly one card for tool %q, got %d", tool.slug, card.Length())
			}
			if !card.Is("a") {
				t.Fatalf("expected tool %q to be a single <a class=admin-tool-card> link", tool.slug)
			}
			href, exists := card.Attr("href")
			if !exists || href != tool.href {
				t.Fatalf("tool %q href = %q, want %q", tool.slug, href, tool.href)
			}
			title := strings.TrimSpace(card.Find(".admin-tool-title").Text())
			if title != tool.title {
				t.Fatalf("tool %q title = %q, want %q", tool.slug, title, tool.title)
			}
			description := strings.TrimSpace(card.Find(".admin-tool-description").Text())
			if description != tool.description {
				t.Fatalf("tool %q description = %q, want %q", tool.slug, description, tool.description)
			}
			graphicSrc, exists := card.Find(".admin-tool-graphic img").Attr("src")
			if !exists || graphicSrc != tool.graphic {
				t.Fatalf("tool %q graphic src = %q, want %q", tool.slug, graphicSrc, tool.graphic)
			}
		}
	})
}

func TestAdminPage_PublisertProgramCardIsWideNonLinkWithPrintLinkAndToggle(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt kontrakten for det brede Publiser program-kortet.",
		When:  "Når adminforsiden rendres.",
		Then:  "Så er kortet ikke en lenke, spenner over hele raden, og inneholder utskriftslenken med ikon og publiseringsbryteren.",
	})

	// Given
	expectedTitle := "Publiser program"
	expectedDescription := "Publiser programmet for festivalen, og åpne for ønsker og påmeldinger."
	expectedPrintLinkText := "Hent utskriftsvennlig versjon av programmet"
	db := testutil.CreateTestDB(t, "admin_landing_publiser_card")
	testutil.MustExec(t, db, insertPublishingStateSQL, 0)

	// When
	doc := templtest.Render(t, adminPage(db))
	card := doc.Find(`[data-admin-tool="publiser"]`)

	// Then
	if card.Length() != 1 {
		t.Fatalf("expected exactly one publiser card, got %d", card.Length())
	}
	if card.Is("a") {
		t.Fatal("expected the publiser card to not itself be a link")
	}
	if !card.HasClass("admin-tool-card") || !card.HasClass("admin-tool-card--wide") {
		classAttr, _ := card.Attr("class")
		t.Fatalf("expected publiser card classes admin-tool-card admin-tool-card--wide, got %q", classAttr)
	}
	title := strings.TrimSpace(card.Find(".admin-tool-title").Text())
	if title != expectedTitle {
		t.Fatalf("publiser card title = %q, want %q", title, expectedTitle)
	}
	description := strings.TrimSpace(card.Find(".admin-tool-description").Text())
	if description != expectedDescription {
		t.Fatalf("publiser card description = %q, want %q", description, expectedDescription)
	}
	graphicSrc, exists := card.Find(".admin-tool-graphic img").Attr("src")
	if !exists || graphicSrc != publiserGraphic {
		t.Fatalf("publiser card graphic src = %q, want %q", graphicSrc, publiserGraphic)
	}

	printLink := card.Find("a.admin-tool-print-link")
	if printLink.Length() != 1 {
		t.Fatalf("expected exactly one print link, got %d", printLink.Length())
	}
	if href, _ := printLink.Attr("href"); href != "/print" {
		t.Fatalf("print link href = %q, want /print", href)
	}
	if printLinkText := strings.TrimSpace(printLink.Text()); !strings.Contains(printLinkText, expectedPrintLinkText) {
		t.Fatalf("print link text = %q, want it to contain %q", printLinkText, expectedPrintLinkText)
	}
	if printLink.Find("svg").Length() == 0 {
		t.Fatal("expected the external-link icon (an inline svg) inside the print link")
	}
	if printLink.Find(`[aria-hidden="true"]`).Length() == 0 {
		t.Fatal("expected the print link's icon to be wrapped in an aria-hidden element")
	}

	toggle := card.Find(`input[type="checkbox"][data-bind="programPublished"]`)
	if toggle.Length() != 1 {
		t.Fatalf("expected exactly one publish toggle, got %d", toggle.Length())
	}
}

// TestAdminPage_CardHrefsResolveToRegisteredRoutes is the "no dead links"
// check. It mounts the real admin routes (SetupAdminRoute), plus the two
// other real top-level route setups that some admin cards point at
// (billettholder admin and /print — see router.go), and walks the resulting
// router to confirm every href rendered on the page resolves to a real,
// registered GET route.
func TestAdminPage_CardHrefsResolveToRegisteredRoutes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt adminforsidens kort og de reelle, registrerte rutene.",
		When:  "Når man går gjennom alle lenker på siden.",
		Then:  "Så finnes det en registrert GET-rute for hver lenke - ingen døde lenker.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "admin_landing_routes")
	testutil.MustExec(t, db, insertPublishingStateSQL, 0)

	router := chi.NewRouter()
	if err := SetupAdminRoute(router, logger, &live.Manager{}, db, nil); err != nil {
		t.Fatalf("failed to set up admin routes: %v", err)
	}
	if err := billettholderadmin.SetupBillettholderAdminRoute(router, &live.Manager{}, logger, db); err != nil {
		t.Fatalf("failed to set up billettholder admin routes: %v", err)
	}
	if err := printfriendly.PrintFriendlyRoute(router, db, nil, logger); err != nil {
		t.Fatalf("failed to set up print-friendly route: %v", err)
	}

	registeredGETRoutes := map[string]bool{}
	if err := chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if method == http.MethodGet {
			registeredGETRoutes[route] = true
		}
		return nil
	}); err != nil {
		t.Fatalf("failed to walk router: %v", err)
	}

	// When
	// Scoped to the admin-tool-card and print links the CONTRACT describes,
	// not to every <a> on the page — the breadcrumb's "Hjem" link to "/" is
	// unrelated to this page's card contract and isn't set up by any of the
	// route functions mounted above.
	doc := templtest.Render(t, adminPage(db))
	linkedHrefs := map[string]bool{}
	doc.Find(".admin-tool-card[href], .admin-tool-print-link[href]").Each(func(_ int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if exists && strings.TrimSpace(href) != "" {
			linkedHrefs[href] = true
		}
	})

	// Then
	if len(linkedHrefs) == 0 {
		t.Fatal("expected the admin page to render at least one card or print link")
	}
	for href := range linkedHrefs {
		if !registeredGETRoutes[href] {
			t.Fatalf("admin page links to %q, but no GET route is registered for it (registered: %v)", href, registeredGETRoutes)
		}
	}
}

// TestAdminPage_GraphicAndFrameAssetsExistAsValidSVG asserts every graphic
// referenced by a card, and the round frame background image, exist on disk
// under static/ and look like real SVG files. Paths are resolved relative to
// this test file so the check works regardless of the working directory the
// test binary runs from.
func TestAdminPage_GraphicAndFrameAssetsExistAsValidSVG(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt filstiene til grafikkene og rammebildet som adminforsiden bruker.",
		When:  "Når filene leses fra disk under static/.",
		Then:  "Så finnes hver fil, og den er en gyldig SVG.",
	})

	// Given
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve this test file's own path")
	}
	staticDir := filepath.Join(filepath.Dir(testFile), "..", "..", "static")

	expectedAssets := []string{
		"/static/admin/algorithm.svg",
		"/static/admin/blueprint.svg",
		"/static/admin/checklist.svg",
		"/static/admin/access-badge.svg",
		"/static/admin/inbox.svg",
		publiserGraphic,
		publiserFrameGraphic,
	}

	// When / Then
	for _, asset := range expectedAssets {
		relPath := strings.TrimPrefix(asset, "/static/")
		fullPath := filepath.Join(staticDir, relPath)

		contents, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("expected asset %q to exist at %s: %v", asset, fullPath, err)
		}

		trimmed := strings.TrimSpace(string(contents))
		if !strings.HasPrefix(trimmed, "<svg") {
			preview := trimmed
			if len(preview) > 40 {
				preview = preview[:40]
			}
			t.Fatalf("expected %q to be a valid svg (start with <svg), got %q...", asset, preview)
		}
	}
}

func newLiveManagerForTest(t *testing.T, cookieName string) *live.Manager {
	t.Helper()

	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{
		Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(),
	}))
	if err != nil {
		t.Fatalf("failed to start embedded nats: %v", err)
	}
	t.Cleanup(func() { ns.Close() })
	ns.WaitForServer()

	manager, err := live.NewManager(context.Background(), ns, sessions.NewCookieStore([]byte(cookieName)))
	if err != nil {
		t.Fatalf("failed to create live manager: %v", err)
	}
	return manager
}

func assertPublishSwitchContract(t *testing.T, sw *goquery.Selection) {
	t.Helper()

	if role, _ := sw.Attr("role"); role != "switch" {
		t.Fatalf("publish switch role = %q, want %q", role, "switch")
	}
	if bind, _ := sw.Attr("data-bind"); bind != "programPublished" {
		t.Fatalf("publish switch data-bind = %q, want %q", bind, "programPublished")
	}
	handler, exists := sw.Attr("data-on:click")
	if !exists || !strings.Contains(handler, "@put('/admin/api/program-publishing')") {
		t.Fatalf("publish switch data-on:click = %q, want it to call @put('/admin/api/program-publishing')", handler)
	}
}

// TestAdminPage_PublishToggleRoundTripUpdatesStateAndMarkup covers the full
// round trip: the rendered card reflects the unpublished state, a real PUT
// through the real programPublishingRoute flips it, and a re-render reflects
// the new, published state - all while the toggle keeps its documented
// Datastar contract.
func TestAdminPage_PublishToggleRoundTripUpdatesStateAndMarkup(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet er upublisert.",
		When:  "Når admin sender PUT til /admin/api/program-publishing for å publisere det.",
		Then:  "Så oppdateres lagret tilstand, og adminforsiden viser publisert status og en avkrysset bryter ved neste rendring.",
	})

	// Given
	expectedInitialStateText := "Ikke publisert"
	expectedAfterStateText := "Publisert"
	db, logger := testutil.CreateTestDBAndLogger(t, "admin_landing_publish_round_trip")
	testutil.MustExec(t, db, insertPublishingStateSQL, 0)

	manager := newLiveManagerForTest(t, "admin-landing-publish-round-trip")
	router := chi.NewRouter()
	if err := SetupAdminRoute(router, logger, manager, db, nil); err != nil {
		t.Fatalf("failed to set up admin routes: %v", err)
	}

	// When (initial render, before publishing)
	initialDoc := templtest.Render(t, adminPage(db))
	initialSwitch := initialDoc.Find(`input[type="checkbox"][data-bind="programPublished"]`)

	// Then (initial render, before publishing)
	if initialSwitch.Length() != 1 {
		t.Fatalf("expected exactly one publish switch before publishing, got %d", initialSwitch.Length())
	}
	if _, checked := initialSwitch.Attr("checked"); checked {
		t.Fatal("expected the publish switch to be unchecked before publishing")
	}
	initialStateText := strings.TrimSpace(initialDoc.Find(".publiser-program-state-value").Text())
	if initialStateText != expectedInitialStateText {
		t.Fatalf("state text before publishing = %q, want %q", initialStateText, expectedInitialStateText)
	}
	assertPublishSwitchContract(t, initialSwitch)

	// When (the real PUT request)
	request := httptest.NewRequest(http.MethodPut, "/admin/api/program-publishing", strings.NewReader(`{"programPublished":true}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Then (the real PUT request)
	if recorder.Code < 200 || recorder.Code >= 300 {
		t.Fatalf("expected a 2xx from the publish PUT, got %d: %s", recorder.Code, recorder.Body.String())
	}
	actualPublished, err := program.IsPublished(db)
	if err != nil {
		t.Fatalf("failed to read program publishing state: %v", err)
	}
	if !actualPublished {
		t.Fatal("expected the program to be published after the PUT")
	}

	// When (re-render, after publishing)
	afterDoc := templtest.Render(t, adminPage(db))
	afterSwitch := afterDoc.Find(`input[type="checkbox"][data-bind="programPublished"]`)

	// Then (re-render, after publishing)
	if afterSwitch.Length() != 1 {
		t.Fatalf("expected exactly one publish switch after publishing, got %d", afterSwitch.Length())
	}
	if _, checked := afterSwitch.Attr("checked"); !checked {
		t.Fatal("expected the publish switch to be checked after publishing")
	}
	afterStateText := strings.TrimSpace(afterDoc.Find(".publiser-program-state-value").Text())
	if afterStateText != expectedAfterStateText {
		t.Fatalf("state text after publishing = %q, want %q", afterStateText, expectedAfterStateText)
	}
	assertPublishSwitchContract(t, afterSwitch)
}

// TestAdminPage_NonAdminCannotTogglePublishState mirrors
// TestFeedbackAdminRoute_NonAdminIsForbiddenAndSeesNoFeedback: a logged-in
// non-admin user must be forbidden from flipping the publish switch, and the
// stored state must not change.
func TestAdminPage_NonAdminCannotTogglePublishState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som ikke er admin, og et upublisert program.",
		When:  "Når brukeren sender PUT til /admin/api/program-publishing gjennom admin-kravet.",
		Then:  "Så skal tilgang nektes, og publiseringsstatusen skal fortsatt være upublisert.",
	})

	// Given
	expectedStatus := http.StatusForbidden
	db, logger := testutil.CreateTestDBAndLogger(t, "admin_landing_non_admin_publish")
	testutil.MustExec(t, db, insertPublishingStateSQL, 0)

	router := chi.NewRouter()
	adminRouter := router.With(
		userctx.UserMiddleware(logger, db),
		authctx.RequireAdmin(logger, authctx.WithForbiddenHandler(userctx.AdminForbiddenHandler(db, logger))),
	)
	if err := SetupAdminRoute(adminRouter, logger, &live.Manager{}, db, nil); err != nil {
		t.Fatalf("failed to set up admin routes: %v", err)
	}

	request := httptest.NewRequest(http.MethodPut, "/admin/api/program-publishing", strings.NewReader(`{"programPublished":true}`))
	request.Header.Set("Content-Type", "application/json")
	request = request.WithContext(authctx.WithUserToken(request.Context(), "admin-landing-non-admin", "non-admin-landing@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	actualPublished, err := program.IsPublished(db)
	if err != nil {
		t.Fatalf("failed to read program publishing state: %v", err)
	}
	if actualPublished {
		t.Fatal("expected the program to remain unpublished after a forbidden request")
	}
}
