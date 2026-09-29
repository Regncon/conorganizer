package feedback

import (
	"errors"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestDelete_RemovesOnlyThatFeedback(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt to lagrede tilbakemeldinger.",
		When:  "Når en administrator sletter den ene.",
		Then:  "Så er bare den andre igjen.",
	})

	// Given
	expectedMessages := []string{"beholdes"}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_delete_one")
	insertFeedback(t, db, CategoryOther, "slettes", "2026-09-20T12:00:00.000Z")
	insertFeedback(t, db, CategoryOther, "beholdes", "2026-09-21T12:00:00.000Z")
	target := findEntry(t, mustList(t, db, ""), "slettes")

	// When
	err := Delete(db, target.ID)

	// Then
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	assertMessages(t, mustList(t, db, ""), expectedMessages)
}

func TestDelete_UnknownFeedbackReturnsErrNotFound(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding som ikke finnes, for eksempel fordi den allerede er slettet.",
		When:  "Når en administrator prøver å slette den.",
		Then:  "Så returneres ErrNotFound, og ingenting annet slettes.",
	})

	// Given
	expectedErr := ErrNotFound
	expectedMessages := []string{"beholdes"}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_delete_unknown")
	insertFeedback(t, db, CategoryOther, "beholdes", "2026-09-21T12:00:00.000Z")

	// When
	err := Delete(db, 999)

	// Then
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
	assertMessages(t, mustList(t, db, ""), expectedMessages)
}

func findEntry(t *testing.T, entries []Entry, message string) Entry {
	t.Helper()
	for _, entry := range entries {
		if entry.Message == message {
			return entry
		}
	}
	t.Fatalf("no feedback with message %q", message)
	return Entry{}
}
