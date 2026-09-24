package billettholderadmin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestGetBillettholderInterestSectionsByBillettholderID_ReturnsAssignedAndInterestedEventsGroupedByPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given expected pulje sections with assigned rows separated from interest rows and pulje publication state.",
		When:  "When the real billettholder interest loader reads the database.",
		Then:  "Then each requested billettholder gets assigned rows first and interest rows grouped by level with the correct pulje-specific publication status.",
	})

	// Given
	expectedBillettholderID := 42
	expectedSections := []expectedBillettholderInterestSection{
		{
			PuljeID: models.PuljeFredagKveld,
			Name:    "Fredag kveld",
			Assigned: []expectedBillettholderInterestRow{
				{
					EventID:       "assigned-gm-without-interest",
					EventTitle:    "Assigned GM Without Interest",
					EventStatus:   models.EventStatusAnnounced,
					InterestLevel: models.InterestLevelNone,
					AssignedRole:  models.EventPlayerRoleGM,
				},
				{
					EventID:       "assigned-player-first-choice",
					EventTitle:    "Assigned Player First Choice",
					EventStatus:   models.EventStatusAnnounced,
					InterestLevel: models.InterestLevelHigh,
					AssignedRole:  models.EventPlayerRolePlayer,
				},
			},
			High: []expectedBillettholderInterestRow{
				{
					EventID:       "high-interest",
					EventTitle:    "High Interest",
					EventStatus:   models.EventStatusAnnounced,
					InterestLevel: models.InterestLevelHigh,
				},
			},
			Medium: []expectedBillettholderInterestRow{
				{
					EventID:       "medium-interest",
					EventTitle:    "Medium Interest",
					EventStatus:   models.EventStatusApproved,
					InterestLevel: models.InterestLevelMedium,
				},
			},
			Low: []expectedBillettholderInterestRow{
				{
					EventID:       "low-interest",
					EventTitle:    "Low Interest",
					EventStatus:   models.EventStatusDraft,
					InterestLevel: models.InterestLevelLow,
				},
			},
		},
		{
			PuljeID: models.PuljeLordagMorgen,
			Name:    "Lørdag morgen",
			Assigned: []expectedBillettholderInterestRow{
				{
					EventID:       "saturday-assigned-player",
					EventTitle:    "Saturday Assigned Player",
					EventStatus:   models.EventStatusAnnounced,
					InterestLevel: models.InterestLevelNone,
					AssignedRole:  models.EventPlayerRolePlayer,
				},
			},
			High: []expectedBillettholderInterestRow{
				{
					EventID:       "saturday-high-interest",
					EventTitle:    "Saturday High Interest",
					EventStatus:   models.EventStatusApproved,
					InterestLevel: models.InterestLevelHigh,
				},
			},
		},
	}

	db := testutil.CreateTestDB(t, "billettholder_interests")
	seedBillettholderInterestLookups(t, db)
	seedBillettholderInterestBillettholdere(t, db, expectedBillettholderID)
	seedBillettholderInterestPuljer(t, db)
	seedBillettholderInterestEvents(t, db)
	seedBillettholderInterestEventPuljer(t, db)
	seedBillettholderInterestRows(t, db, expectedBillettholderID)
	seedBillettholderInterestAssignments(t, db, expectedBillettholderID)

	// When
	actualSectionsByBillettholderID, err := getBillettholderInterestSectionsByBillettholderID(db, []int{expectedBillettholderID})

	// Then
	if err != nil {
		t.Fatalf("expected billettholder interest sections to load: %v", err)
	}

	actualSections, ok := actualSectionsByBillettholderID[expectedBillettholderID]
	if !ok {
		t.Fatalf("expected sections for billettholder ID %d", expectedBillettholderID)
	}

	assertBillettholderInterestSections(t, expectedSections, actualSections)
}

func TestGetBillettholderInterestSectionsByBillettholderID_WhenAssignedAsPlayerAndGM_ReturnsOneArrangementWithBothRoles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given one expected assigned arrangement carrying both the Player and GM roles.",
		When:  "When the billettholder interest loader and tabs render the dual assignment.",
		Then:  "Then the arrangement is counted once and both role labels are shown.",
	})

	// Given
	expectedBillettholderID := 42
	expectedAssignedTotal := 1
	expectedRoleLabels := []string{"Spiller", "Spilleder (GM/DM)"}

	db := testutil.CreateTestDB(t, "billettholder_interests_dual_role")
	seedBillettholderInterestLookups(t, db)
	seedBillettholderInterestBillettholdere(t, db, expectedBillettholderID)
	seedBillettholderInterestPuljer(t, db)
	seedBillettholderInterestEvents(t, db)
	seedBillettholderInterestEventPuljer(t, db)
	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO relation_events_players (
			event_id, pulje_id, billettholder_id, role
		) VALUES
			('assigned-player-first-choice', ?, ?, ?),
			('assigned-player-first-choice', ?, ?, ?)
	`,
		models.PuljeFredagKveld, expectedBillettholderID, models.EventPlayerRolePlayer,
		models.PuljeFredagKveld, expectedBillettholderID, models.EventPlayerRoleGM,
	)

	// When
	sectionsByBillettholderID, err := getBillettholderInterestSectionsByBillettholderID(db, []int{expectedBillettholderID})

	// Then
	if err != nil {
		t.Fatalf("expected billettholder interest sections to load: %v", err)
	}
	sections := sectionsByBillettholderID[expectedBillettholderID]
	if actualAssignedTotal := billettholderAssignedTotal(sections); actualAssignedTotal != expectedAssignedTotal {
		t.Fatalf("assigned arrangement count mismatch: expected %d, got %d", expectedAssignedTotal, actualAssignedTotal)
	}
	if len(sections) != 1 || len(sections[0].Assigned) != 1 {
		t.Fatalf("expected one assigned arrangement in one pulje, got %#v", sections)
	}
	assigned := sections[0].Assigned[0]
	if !assigned.IsPlayer || !assigned.IsGM {
		t.Fatalf("expected assigned arrangement to carry Player and GM roles, got %#v", assigned)
	}

	doc := templtest.Render(t, billettholderInterestTabPanels(sections))
	assignedText := strings.Join(strings.Fields(doc.Find(".billettholder-interest-group a").Text()), " ")
	for _, expectedRoleLabel := range expectedRoleLabels {
		if !strings.Contains(assignedText, expectedRoleLabel) {
			t.Fatalf("expected assigned arrangement to show %q, got %q", expectedRoleLabel, assignedText)
		}
	}
}
