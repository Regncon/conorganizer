package billettholderadmin

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestGetBillettholderInterestSectionsByBillettholderID_ReturnsAssignedAndInterestedEventsGroupedByPulje(t *testing.T) {
	// Given expected pulje sections with assigned rows separated from interest rows and pulje publication state,
	// when the real billettholder interest loader reads the database,
	// then each requested billettholder gets assigned rows first and interest rows grouped by level with the correct pulje-specific publication status.

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
					IsPublished:   true,
					InterestLevel: models.InterestLevelNone,
					AssignedRole:  models.EventPlayerRoleGM,
				},
				{
					EventID:       "assigned-player-first-choice",
					EventTitle:    "Assigned Player First Choice",
					EventStatus:   models.EventStatusAnnounced,
					IsPublished:   true,
					InterestLevel: models.InterestLevelHigh,
					AssignedRole:  models.EventPlayerRolePlayer,
				},
			},
			High: []expectedBillettholderInterestRow{
				{
					EventID:       "high-interest",
					EventTitle:    "High Interest",
					EventStatus:   models.EventStatusAnnounced,
					IsPublished:   false,
					InterestLevel: models.InterestLevelHigh,
				},
			},
			Medium: []expectedBillettholderInterestRow{
				{
					EventID:       "medium-interest",
					EventTitle:    "Medium Interest",
					EventStatus:   models.EventStatusApproved,
					IsPublished:   true,
					InterestLevel: models.InterestLevelMedium,
				},
			},
			Low: []expectedBillettholderInterestRow{
				{
					EventID:       "low-interest",
					EventTitle:    "Low Interest",
					EventStatus:   models.EventStatusDraft,
					IsPublished:   false,
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
					IsPublished:   true,
					InterestLevel: models.InterestLevelNone,
					AssignedRole:  models.EventPlayerRolePlayer,
				},
			},
			High: []expectedBillettholderInterestRow{
				{
					EventID:       "saturday-high-interest",
					EventTitle:    "Saturday High Interest",
					EventStatus:   models.EventStatusApproved,
					IsPublished:   true,
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
