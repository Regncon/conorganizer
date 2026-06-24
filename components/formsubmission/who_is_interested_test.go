package formsubmission

import (
	"testing"

	puljefordelingservice "github.com/Regncon/conorganizer/service/puljefordeling"
)

const (
	eventE1 = "E1"
	eventE2 = "E2"
	eventE3 = "E3"
	eventE4 = "E4"

	puljeP1 = "P1"
	puljeP2 = "P2"
	puljeP3 = "P3"
	puljeP4 = "P4"

	idPlayerAssigned             = 1
	idGMAssigned                 = 2
	idNotVeryInterested          = 3
	idUnassigned                 = 4
	idSameEventAssignee          = 5
	idGMPlayer                   = 6
	idGMAndPlayerDifferentEvents = 7
	idGMOnlyVeryInterestedOther  = 8
)

func TestFirstChoiceStatusText_PrefersOtherFirstChoiceOverCurrentPulje(t *testing.T) {
	row := InterestWithHolder{
		HasCurrentPuljeFirstChoice: true,
		HasOtherPuljeFirstChoice:   true,
	}

	got := firstChoiceStatusText(row)

	if got != "Fått førsteval frå før" {
		t.Fatalf("status text mismatch: %q", got)
	}
}

func TestFirstChoiceButtonAction_ReturnsExpectedState(t *testing.T) {
	setAction := "set"
	removeAction := "remove"

	tests := []struct {
		name string
		row  InterestWithHolder
		want ButtonInfo
	}{
		{
			name: "default can set first-choice",
			row:  InterestWithHolder{},
			want: ButtonInfo{Label: "Set førsteval", Action: setAction},
		},
		{
			name: "GM cannot set first-choice",
			row:  InterestWithHolder{IsGamemaster: true},
			want: ButtonInfo{
				Label:    "Set førsteval",
				Disabled: true,
				Title:    "GM kan ikkje ha førsteval for dette arrangementet",
			},
		},
		{
			name: "other pulje first-choice cannot set another",
			row:  InterestWithHolder{HasOtherPuljeFirstChoice: true},
			want: ButtonInfo{
				Label:    "Set førsteval",
				Disabled: true,
				Title:    "Har allereie fått førsteval",
			},
		},
		{
			name: "current first-choice can be removed",
			row:  InterestWithHolder{HasCurrentPuljeFirstChoice: true},
			want: ButtonInfo{Label: "Fjern førsteval", Action: removeAction},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := firstChoiceButtonAction(test.row, setAction, removeAction)
			if got != test.want {
				t.Fatalf("button state mismatch: want %+v, got %+v", test.want, got)
			}
		})
	}
}

func TestApplyFirstChoiceStatusesToRows_AddsServiceStatusToRows(t *testing.T) {
	rows := []InterestWithHolder{{
		BillettholderId: 42,
		EventId:         "event-1",
		PuljeId:         puljeP1,
	}}
	statuses := map[puljefordelingservice.FirstChoiceKey]puljefordelingservice.FirstChoiceStatus{
		{BillettholderID: 42, EventID: "event-1", PuljeID: puljeP1}: {
			HasCurrentPuljeFirstChoice: true,
			HasOtherPuljeFirstChoice:   false,
		},
	}

	applyFirstChoiceStatusesToRows(rows, statuses)

	if !rows[0].HasCurrentPuljeFirstChoice {
		t.Fatalf("expected current first-choice status to be applied")
	}
}
