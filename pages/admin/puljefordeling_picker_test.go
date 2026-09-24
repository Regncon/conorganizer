package admin

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeEventBox_MissingGMShortcutRequiresEditablePulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Et arrangement mangler spilleder eller har allerede en spilleder.", When: "Kortet vises i en redigerbar eller publisert pulje.", Then: "Bare et arrangement uten spilleder i en redigerbar pulje tilbyr en knapp for å legge til spilleder."})
	for _, tc := range []struct {
		name      string
		gm        string
		published bool
		expected  int
	}{
		{name: "missing GM", expected: 1},
		{name: "published", published: true},
		{name: "already assigned", gm: "Kari Nordmann"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			expectedButtons := tc.expected
			event := puljefordeling.EmulatedEvent{EventID: "evA", Title: "Drager", GMName: tc.gm}

			// When
			doc := templtest.Render(t, puljeEventBox(models.PuljeFredagKveld, event, tc.published, nil))

			// Then
			button := doc.Find("button:contains('Mangler spilleder · Legg til spilleder')")
			if button.Length() != expectedButtons {
				t.Fatalf("missing-GM buttons = %d, want %d", button.Length(), expectedButtons)
			}
			if expectedButtons > 0 && (button.AttrOr("type", "") != "button" || button.AttrOr("aria-controls", "") != "puljefordeling-assign-dialog") {
				t.Error("missing-GM shortcut must be a native button identifying its dialog")
			}
		})
	}
}

func TestPuljeAssignmentPicker_MissingGMShortcutSetsContextAndGenericAddRestoresRoles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "En administrator vil legge til spilleder på ett arrangement, deretter deltaker på et annet.", When: "De renderte knappene åpner og bruker den delte velgeren.", Then: "Spillederflyten sender riktig arrangement, pulje og rolle, fokuserer søket, og vanlig Legg til gjenoppretter begge rollevalg."})
	// Given
	expectedGMTitle := "Legg til spilleder på Drager's \"tårn\""
	expectedGenericTitle := "Legg til deltaker på Arrangement B"
	db, _ := tildelingsFixture(t)
	missingGM := templtest.Render(t, puljeEventBox(models.PuljeFredagKveld, puljefordeling.EmulatedEvent{EventID: "evA", Title: "Drager's \"tårn\""}, false, nil))
	generic := templtest.Render(t, puljeEventBox(models.PuljeLordagKveld, puljefordeling.EmulatedEvent{EventID: "evB", Title: "Arrangement B"}, false, nil))
	page := templtest.Render(t, puljefordelingIndex(db, testutil.NewTestLogger(), models.PuljeFredagKveld, nil))
	picker := page.Find("#puljefordeling-assign-dialog")
	if picker.Find("admin-billettholder-search").AttrOr("data-attr:data-clear-input", "") != "$clearInput" {
		t.Fatal("search reset must bind to the component observed data-clear-input attribute")
	}
	payload, err := json.Marshal(map[string]string{
		"gmOpen":               missingGM.Find("button:contains('Mangler spilleder · Legg til spilleder')").AttrOr("data-on:click", ""),
		"genericOpen":          generic.Find(".pulje-add").AttrOr("data-on:click", ""),
		"dialogEffect":         picker.AttrOr("data-effect", ""),
		"interestsShown":       picker.Find(".pulje-assignment-interests").AttrOr("data-show", "true"),
		"interestActionsShown": picker.Find(".pulje-assignment-interesse-actions").AttrOr("data-show", "true"),
		"title":                picker.Find("h3").AttrOr("data-text", "''"),
		"playerShown":          picker.Find("button:contains('Tildel som spiller')").AttrOr("data-show", "true"),
		"gmPrimary":            picker.Find("button:contains('Tildel som spilleder')").AttrOr("data-class:btn--primary", "false"),
		"gmSubmission":         picker.Find("button:contains('Tildel som spilleder')").AttrOr("data-on:click", ""),
	})
	if err != nil {
		t.Fatal(err)
	}

	// When
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node is required to execute the rendered picker actions")
	}
	command := exec.Command(node, "-e", pickerActionScript)
	command.Stdin = bytes.NewReader(payload)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("execute picker actions: %v\n%s", err, output)
	}
	var actual struct {
		GMTitle                     string
		GenericTitle                string
		GMPlayerShown               bool
		GMInterestsShown            bool
		GMInterestActionsShown      bool
		GenericInterestsShown       bool
		GenericInterestActionsShown bool
		GMPrimary                   bool
		GenericPrimary              bool
		PlayerShown                 bool
		Focused                     bool
		Cleared                     bool
		EventID                     string
		PuljeID                     string
		Role                        string
		HolderID                    int
		GenericEvent                string
		GenericPulje                string
	}
	if err := json.Unmarshal(output, &actual); err != nil {
		t.Fatalf("decode picker action result: %v\n%s", err, output)
	}

	// Then
	if actual.GMTitle != expectedGMTitle || actual.GenericTitle != expectedGenericTitle {
		t.Errorf("picker titles = %q / %q, want %q / %q", actual.GMTitle, actual.GenericTitle, expectedGMTitle, expectedGenericTitle)
	}
	if actual.GMInterestsShown || actual.GMInterestActionsShown || !actual.GenericInterestsShown || !actual.GenericInterestActionsShown || actual.GMPlayerShown || !actual.PlayerShown || !actual.GMPrimary || actual.GenericPrimary || !actual.Focused || !actual.Cleared {
		t.Errorf("picker mode, focus or stale selection incorrect: %+v", actual)
	}
	if actual.EventID != "evA" || actual.PuljeID != "FredagKveld" || actual.Role != "GM" || actual.HolderID != 42 {
		t.Errorf("GM submission lost target or role: %+v", actual)
	}
	if actual.GenericEvent != "evB" || actual.GenericPulje != "LordagKveld" {
		t.Errorf("generic add retained the previous event or pulje: %+v", actual)
	}
}

const pickerActionScript = `
const vm = require('node:vm');
const actions = JSON.parse(require('node:fs').readFileSync(0, 'utf8'));
let focused = false;
const input = { value: 'Previously selected person', focus() { focused = true; } };
const search = { clearSearch() { input.value = ''; }, shadowRoot: { querySelector() { return input; } } };
const dialog = { open: false, showModal() { this.open = true; }, close() { this.open = false; }, querySelector() { return search; } };
const context = vm.createContext({
  el: dialog,
  document: { getElementById() { return dialog; } },
  evt: { currentTarget: { closest() { return { querySelector() { return search; } }; } } },
  $assignmentBillettholderId: 99, $clearInput: 0,
});
const runAction = (action) => vm.runInContext('(function() {' + action + '})()', context);
runAction(actions.gmOpen);
if (context.$clearInput > 0) search.clearSearch();
runAction(actions.dialogEffect);
const result = {
  GMInterestsShown: vm.runInContext(actions.interestsShown, context),
  GMInterestActionsShown: vm.runInContext(actions.interestActionsShown, context),
  GMTitle: vm.runInContext(actions.title, context),
  GMPlayerShown: vm.runInContext(actions.playerShown, context),
  GMPrimary: vm.runInContext(actions.gmPrimary, context),
  Focused: dialog.open && focused,
  Cleared: input.value === '' && context.$assignmentBillettholderId === 0,
};
context.$assignmentBillettholderId = 42;
context.post = () => Object.assign(result, {
  EventID: context.$assignmentEventId, PuljeID: context.$assignmentPuljeId,
  Role: context.$assignmentRole, HolderID: context.$assignmentBillettholderId,
});
runAction(actions.gmSubmission.replaceAll('@post(', 'post('));
runAction(actions.genericOpen);
Object.assign(result, {
  GenericInterestsShown: vm.runInContext(actions.interestsShown, context),
  GenericInterestActionsShown: vm.runInContext(actions.interestActionsShown, context),
  GenericTitle: vm.runInContext(actions.title, context),
  PlayerShown: vm.runInContext(actions.playerShown, context),
  GenericPrimary: vm.runInContext(actions.gmPrimary, context),
  GenericEvent: context.$assignmentEventId, GenericPulje: context.$assignmentPuljeId,
});
process.stdout.write(JSON.stringify(result));
`
