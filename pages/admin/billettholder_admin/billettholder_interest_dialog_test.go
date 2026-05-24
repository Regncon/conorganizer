package billettholderadmin

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
)

func TestMockBillettholderInterestSections(t *testing.T) {
	// Given the expected pulje layout and interest level counts,
	// when the mock interest sections are built,
	// then every pulje has assigned rows and sorted interest groups.

	// Given
	expectedPuljes := []struct {
		id   models.Pulje
		name string
	}{
		{id: models.PuljeFredagKveld, name: "Fredag kveld"},
		{id: models.PuljeLordagMorgen, name: "Lørdag morgen"},
		{id: models.PuljeLordagKveld, name: "Lørdag kveld"},
		{id: models.PuljeSondagMorgen, name: "Søndag morgen"},
	}
	expectedAssignedRowsPerPulje := 2
	expectedInterestRowsPerLevel := 10
	expectedInterestGroups := []struct {
		name  string
		level models.InterestLevel
	}{
		{name: "high", level: models.InterestLevelHigh},
		{name: "medium", level: models.InterestLevelMedium},
		{name: "low", level: models.InterestLevelLow},
	}
	billettholder := models.Billettholder{ID: 6}

	// When
	sections := mockBillettholderInterestSections(billettholder)

	// Then
	if len(sections) != len(expectedPuljes) {
		t.Fatalf("expected %d pulje sections, got %d", len(expectedPuljes), len(sections))
	}

	for i, expectedPulje := range expectedPuljes {
		section := sections[i]
		if section.PuljeID != expectedPulje.id {
			t.Fatalf("section %d expected pulje ID %q, got %q", i, expectedPulje.id, section.PuljeID)
		}
		if section.Name != expectedPulje.name {
			t.Fatalf("section %d expected pulje name %q, got %q", i, expectedPulje.name, section.Name)
		}
		if len(section.Assigned) != expectedAssignedRowsPerPulje {
			t.Fatalf("section %q expected %d assigned rows, got %d", section.Name, expectedAssignedRowsPerPulje, len(section.Assigned))
		}

		interestGroups := []struct {
			name string
			rows []billettholderInterestMockRow
		}{
			{name: "high", rows: section.High},
			{name: "medium", rows: section.Medium},
			{name: "low", rows: section.Low},
		}

		for j, expectedGroup := range expectedInterestGroups {
			group := interestGroups[j]
			if group.name != expectedGroup.name {
				t.Fatalf("section %q expected interest group %q at position %d, got %q", section.Name, expectedGroup.name, j, group.name)
			}
			if len(group.rows) != expectedInterestRowsPerLevel {
				t.Fatalf("section %q group %q expected %d rows, got %d", section.Name, group.name, expectedInterestRowsPerLevel, len(group.rows))
			}
			for rowIndex, row := range group.rows {
				if row.InterestLevel != expectedGroup.level {
					t.Fatalf("section %q group %q row %d expected interest level %q, got %q", section.Name, group.name, rowIndex, expectedGroup.level, row.InterestLevel)
				}
			}
		}
	}
}
