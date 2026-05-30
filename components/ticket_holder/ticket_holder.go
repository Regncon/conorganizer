package ticketholder

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/Regncon/conorganizer/models"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
	"github.com/Regncon/conorganizer/service/requestctx"
)

type BillettHolder struct {
	Email string
	Name  string
	Id    int
	Color string
}

type PuljeInterestAvailability string

const (
	PuljeInterestOpen          PuljeInterestAvailability = "open"
	PuljeInterestWarning       PuljeInterestAvailability = "warning"
	PuljeInterestUrgentWarning PuljeInterestAvailability = "urgent-warning"
	PuljeInterestLocked        PuljeInterestAvailability = "locked"
	PuljeInterestCompleted     PuljeInterestAvailability = "completed"
)

type PuljeInterestState struct {
	PuljeID         models.Pulje
	PuljeName       string
	Availability    PuljeInterestAvailability
	Message         string
	CanEdit         bool
	ShowProfileLink bool
	Priority        int
}

func (state PuljeInterestState) ClassName() string {
	return fmt.Sprintf("pulje-interest-state--%s", state.Availability)
}

func (state PuljeInterestState) HasMessage() bool {
	return state.Message != ""
}

func (state PuljeInterestState) IsWarning() bool {
	return state.Availability == PuljeInterestWarning || state.Availability == PuljeInterestUrgentWarning
}

func (state PuljeInterestState) IsLocked() bool {
	return state.Availability == PuljeInterestLocked
}

func (state PuljeInterestState) IsCompleted() bool {
	return state.Availability == PuljeInterestCompleted
}

func (state PuljeInterestState) SignalPatch() string {
	return fmt.Sprintf(
		"$puljeAvailability = %q; $puljeCanEdit = %t; $puljeStatusMessage = %q; $puljeShowProfileLink = %t;",
		state.Availability,
		state.CanEdit,
		state.Message,
		state.ShowProfileLink,
	)
}

func BuildPuljeInterestState(pulje models.PuljeRow, now time.Time) PuljeInterestState {
	state := PuljeInterestState{
		PuljeID:      pulje.ID,
		PuljeName:    pulje.Name,
		Availability: PuljeInterestOpen,
		CanEdit:      true,
		Priority:     0,
	}

	switch pulje.Status {
	case models.PuljeStatusLocked:
		state.Availability = PuljeInterestLocked
		state.Message = "Puljen er låst. Du kan ikke melde eller endre interesse lenger. Vi jobber med å fordele spillere."
		state.CanEdit = false
		state.Priority = 1
		return state
	case models.PuljeStatusCompleted:
		state.Availability = PuljeInterestCompleted
		state.Message = "Puljefordelingen er klar. Se hva du fikk på profilen din."
		state.CanEdit = false
		state.ShowProfileLink = true
		state.Priority = 0
		return state
	}

	lockAt, hasLockAt := puljeLockAt(pulje)
	if !hasLockAt {
		return state
	}

	urgentStartsAt := lockAt.Add(-30 * time.Minute)
	warningStartsAt := lockAt.Add(-2 * time.Hour)
	lockTimeLabel := lockAt.Format("15:04")

	if !now.Before(urgentStartsAt) {
		state.Availability = PuljeInterestUrgentWarning
		state.Message = fmt.Sprintf("Puljen låses straks, kl %s. Gjør endringer nå hvis du vil endre interessen din.", lockTimeLabel)
		state.Priority = 3
		return state
	}
	if !now.Before(warningStartsAt) {
		state.Availability = PuljeInterestWarning
		state.Message = fmt.Sprintf("Puljen låses snart, kl %s.", lockTimeLabel)
		state.Priority = 2
		return state
	}

	return state
}

func puljeLockAt(pulje models.PuljeRow) (time.Time, bool) {
	if pulje.StartAt.IsZero() {
		return time.Time{}, false
	}
	return pulje.StartAt.TimeOrZero().Add(-30 * time.Minute), true
}

func BuildSelectedPuljeInterestState(puljer []models.PuljeRow, puljeID string, now time.Time) PuljeInterestState {
	for _, pulje := range puljer {
		if string(pulje.ID) == puljeID {
			return BuildPuljeInterestState(pulje, now)
		}
	}
	if len(puljer) > 0 {
		return BuildPuljeInterestState(puljer[0], now)
	}
	return PuljeInterestState{Availability: PuljeInterestOpen, CanEdit: true}
}

func BuildMostUrgentPuljeInterestState(puljer []models.PuljeRow, now time.Time) (PuljeInterestState, bool) {
	var selected PuljeInterestState
	hasSelected := false

	for _, pulje := range puljer {
		state := BuildPuljeInterestState(pulje, now)
		if !state.HasMessage() {
			continue
		}
		if !hasSelected || state.Priority > selected.Priority {
			selected = state
			hasSelected = true
			continue
		}
		if state.Priority == selected.Priority && pulje.StartAt.TimeOrZero().Before(selectedStartTime(puljer, selected.PuljeID)) {
			selected = state
		}
	}

	return selected, hasSelected
}

func selectedStartTime(puljer []models.PuljeRow, puljeID models.Pulje) time.Time {
	for _, pulje := range puljer {
		if pulje.ID == puljeID {
			return pulje.StartAt.TimeOrZero()
		}
	}
	return time.Time{}
}

func GetTicketHolders(userInfo requestctx.UserRequestInfo, db *sql.DB) ([]BillettHolder, error) {
	// todo: use the correct way to get billettholders (billettholderservice.GetBillettholdere has a fallback to get all billettholders)
	query := `
    SELECT
        [be].email,
        [be].billettholder_id,
        [bh].first_name,
        [bh].last_name
    FROM
        relation_billettholder_emails [be]
        LEFT JOIN billettholdere [bh] ON [be].billettholder_id = [bh].id
	    WHERE
	        [be].kind = ?
	        AND [be].billettholder_id IN (
            SELECT
                billettholder_id
            FROM
                relation_billettholder_emails
            WHERE
                email = ?
        )
`
	rows, ticketHolderQueryErr := db.Query(query, models.BillettholderEmailKindTicket, userInfo.Email)
	if ticketHolderQueryErr != nil {
		return nil, fmt.Errorf("failed to query ticket holders for email %q: %w", userInfo.Email, ticketHolderQueryErr)
	}
	defer rows.Close()

	var email, firstName, lastName string
	var associatedTicketholders []BillettHolder

	for rows.Next() {
		var billettHolderId int

		if ticketHolderScanErr := rows.Scan(&email, &billettHolderId, &firstName, &lastName); ticketHolderScanErr != nil {
			return nil, fmt.Errorf("failed to scan ticket holder row: %w", ticketHolderScanErr)
		}
		associatedTicketholders = append(associatedTicketholders, BillettHolder{
			Email: email,
			Name:  fmt.Sprintf("%s %s", firstName, lastName),
			Id:    billettHolderId,
			Color: ColorForName(fmt.Sprintf("%s %s", firstName, lastName)),
		})

	}
	if ticketHolderRowsErr := rows.Err(); ticketHolderRowsErr != nil {
		return nil, fmt.Errorf("error iterating over ticket holder rows: %w", ticketHolderRowsErr)
	}

	return associatedTicketholders, nil

}

func GetPuljerFromEventId(eventId string, db *sql.DB) ([]models.PuljeRow, error) {
	puljer, err := puljerService.GetActivePuljeForEvent(eventId, db)
	if err != nil {
		return nil, fmt.Errorf("failed to query event puljer for event %s: %w", eventId, err)
	}

	return puljer, nil
}

func GetYourBillettHolderInfo(userInfo requestctx.UserRequestInfo, ticketHolders []BillettHolder) BillettHolder {
	idx := slices.IndexFunc(ticketHolders, func(th BillettHolder) bool {
		return th.Email == userInfo.Email
	})

	if idx == -1 {
		return BillettHolder{
			Email: "unknown@example.com",
			Name:  "Unknown Ticket Holder",
		}
	}

	return ticketHolders[idx]
}

func GetInitials(s string) string {
	var initialsBuilder strings.Builder
	words := strings.Fields(s)
	if len(words) == 0 {
		return "TT"
	}

	firstWordRunes := []rune(words[0])
	lastWordRunes := []rune(words[len(words)-1])
	if len(firstWordRunes) == 0 || len(lastWordRunes) == 0 {
		return "TT"
	}
	firstChar := firstWordRunes[0]
	lastChar := lastWordRunes[0]

	if unicode.IsLetter(firstChar) {
		initialsBuilder.WriteRune(unicode.ToUpper(firstChar))
	}
	if !unicode.IsLetter(firstChar) {
		initialsBuilder.WriteRune('T')
	}

	if unicode.IsLetter(lastChar) {
		initialsBuilder.WriteRune(unicode.ToUpper(lastChar))
	}
	if !unicode.IsLetter(lastChar) {
		initialsBuilder.WriteRune('T')
	}

	return initialsBuilder.String()
}

func GetFirstNameAndLastInitial(fullName string) string {
	nameParts := strings.Fields(fullName)
	if len(nameParts) == 0 {
		return ""
	}

	if len(nameParts) == 1 {
		return nameParts[0]
	}

	lastNameRunes := []rune(nameParts[len(nameParts)-1])
	if len(lastNameRunes) == 0 {
		return nameParts[0]
	}

	return fmt.Sprintf("%s %c", nameParts[0], unicode.ToUpper(lastNameRunes[0]))
}
