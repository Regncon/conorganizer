package components

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventCard_RendersTheFullOrdinaryCardInAPuljeContext(t *testing.T) {
	event := models.EventCardModel{
		Id:                "ordinary-event",
		Title:             "Ordinary Event",
		Intro:             "An ordinary event introduction",
		System:            "Ordinary System",
		HostName:          "Ordinary Host",
		EventType:         models.EventTypeRoleplay,
		AgeGroup:          models.AgeGroupAdultsOnly,
		Runtime:           models.RunTimeLongRunning,
		BeginnerFriendly:  true,
		CanBeRunInEnglish: true,
	}

	doc := templtest.Render(t, EventCard(event, nil, "LordagMorgen", "2026-10-10", nil))
	card := doc.Find(".event-card-container")

	if got := card.Length(); got != 1 {
		t.Fatalf("ordinary card count = %d, want 1", got)
	}
	if got := card.AttrOr("href", ""); got != "/event/ordinary-event?date=2026-10-10&pulje=LordagMorgen" {
		t.Fatalf("ordinary card href = %q", got)
	}
	if got := doc.Find(".raffle-event-card").Length(); got != 0 {
		t.Fatalf("ordinary card rendered %d raffle card classes", got)
	}
	if got := doc.Find(".event-card-subtitle").Length(); got != 1 {
		t.Fatalf("ordinary card subtitle count = %d, want 1", got)
	}
	if got := doc.Find(".event-card-main-body .event-card-description").Length(); got != 1 {
		t.Fatalf("ordinary card body description count = %d, want 1", got)
	}
	if got := doc.Find(".event-card-footer-description").Length(); got != 0 {
		t.Fatalf("ordinary card footer description count = %d, want 0", got)
	}
	if got := doc.Find(".event-card-footer-gamemaster").Length(); got != 1 {
		t.Fatalf("ordinary card game-master count = %d, want 1", got)
	}
	if got := doc.Find(".event-card-tagicon-container").Length(); got != 5 {
		t.Fatalf("ordinary card tag icon count = %d, want 5", got)
	}
}

func TestProgramEventCard_RendersTheReducedProgramPresentation(t *testing.T) {
	event := models.EventCardModel{
		Id:                "program-event",
		Title:             "Program Event",
		Intro:             "A program event introduction",
		System:            "Hidden System",
		HostName:          "Hidden Host",
		EventType:         models.EventTypeBoardGame,
		AgeGroup:          models.AgeGroupChildFriendly,
		Runtime:           models.RunTimeShortRunning,
		BeginnerFriendly:  true,
		CanBeRunInEnglish: true,
	}

	doc := templtest.Render(t, ProgramEventCard(event, nil, "LordagKveld", "2026-10-10", nil))
	card := doc.Find(".program-event-card")

	if got := card.Length(); got != 1 {
		t.Fatalf("program card count = %d, want 1", got)
	}
	if got := card.AttrOr("href", ""); got != "/event/program-event?date=2026-10-10&pulje=LordagKveld" {
		t.Fatalf("program card href = %q", got)
	}
	if got := card.Find(".event-card-subtitle").Length(); got != 0 {
		t.Fatalf("program card subtitle count = %d, want 0", got)
	}
	if got := card.Find(".event-card-footer-gamemaster").Length(); got != 0 {
		t.Fatalf("program card game-master count = %d, want 0", got)
	}
	if got := card.Find(".event-card-main-body .event-card-description").Length(); got != 0 {
		t.Fatalf("program card body description count = %d, want 0", got)
	}
	if got := card.Find(".event-card-footer-description").Length(); got != 1 {
		t.Fatalf("program card footer description count = %d, want 1", got)
	}
	if got := card.Find(".event-card-tagicon-container").Length(); got != 5 {
		t.Fatalf("program card tag icon count = %d, want 5", got)
	}
}

func TestProgramEventCard_ShowsEachTimeOnItsOwnLine(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement med to tidspunkter samme dag.",
		When:  "Når programkortet vises.",
		Then:  "Så står hvert tidspunkt på egen linje.",
	})

	// Given
	expectedTimes := []string{"10:00–12:00", "19:00–21:00"}
	event := models.EventCardModel{Id: "program-event", Title: "Program Event"}

	// When
	doc := templtest.Render(t, ProgramEventCard(event, nil, "LordagMorgen", "2026-10-03", expectedTimes))

	// Then
	got := templtest.CollectTexts(doc, ".program-event-card-times li")
	if !slices.Equal(got, expectedTimes) {
		t.Fatalf("program card times = %q, want %q", got, expectedTimes)
	}
}

func TestProgramEventCard_ShowsNoTimeWhenNoneIsWritten(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement uten tidspunkt.",
		When:  "Når programkortet vises.",
		Then:  "Så vises ingen tid på kortet.",
	})

	// Given
	expectedTimeLists := 0
	event := models.EventCardModel{Id: "program-event", Title: "Program Event"}

	// When
	doc := templtest.Render(t, ProgramEventCard(event, nil, "LordagMorgen", "2026-10-03", nil))

	// Then
	if got := doc.Find(".program-event-card-times").Length(); got != expectedTimeLists {
		t.Fatalf("program card time lists = %d, want %d", got, expectedTimeLists)
	}
}
