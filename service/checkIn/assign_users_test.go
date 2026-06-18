package checkIn

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestAssociateUserWithBillettholder_WhenEmailsMatchCaseInsensitively_CreatesUserAssociations(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a user whose email matches multiple billettholder emails with different casing.",
		When:  "When the user is associated with billettholdere.",
		Then:  "Then the user is linked to every matching billettholder.",
	})

	// Given
	expectedAssociations := []models.BillettholderUsers{
		{BillettholderID: 8888, UserID: 1},
		{BillettholderID: 9999, UserID: 1},
	}
	matchingEmail := "test@regncon.no"

	db, logger := createCheckInTestDB(t)
	insertUser(t, db, 1, "test-user", matchingEmail)
	insertBillettholder(t, db, 9999)
	insertBillettholder(t, db, 8888)
	insertBillettholder(t, db, 7777)
	insertManualBillettholderEmail(t, db, 9999, matchingEmail)
	insertManualBillettholderEmail(t, db, 8888, "TEST@REGNCON.NO")
	insertManualBillettholderEmail(t, db, 7777, "other@regncon.no")

	// When
	err := AssociateUserWithBillettholder("test-user", db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected user association to succeed: %v", err)
	}
	assertBillettholderUserAssociations(t, db, expectedAssociations)
}

func TestAssociateUsersWithBillettholderEmail_CreatesAssociationForMatchingUserEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at ein billettholder har fått lagt til ei manuell e-postadresse, og ein eksisterande brukar har same e-postadresse med annan casing.",
		When:  "Når e-postadressa blir forsona mot brukarar.",
		Then:  "Så skal billettholderen få ei varig brukar-tilknyting.",
	})

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	manualEmail := "participant@example.com"
	userEmail := "Participant@Example.com"

	db, slogger := createCheckInTestDB(t)
	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", userEmail)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, manualEmail)

	// When
	err := AssociateUsersWithBillettholderEmail(expectedAssociation.BillettholderID, manualEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected association to succeed: %v", err)
	}
	assertOnlyBillettholderUserAssociation(t, db, expectedAssociation)
}

func TestAssociateUsersWithBillettholderEmail_DoesNotDuplicateExistingAssociation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at ein billettholder allereie er knytt til ein brukar via ei manuell e-postadresse.",
		When:  "Når same e-postforsoning køyrer på nytt.",
		Then:  "Så skal det framleis berre finnast ei brukar-tilknyting.",
	})

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	expectedAssociationCount := 1
	manualEmail := "participant@example.com"

	db, slogger := createCheckInTestDB(t)
	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", manualEmail)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, manualEmail)
	insertBillettholderUserAssociation(t, db, expectedAssociation)

	// When
	err := AssociateUsersWithBillettholderEmail(expectedAssociation.BillettholderID, manualEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected repeated association to succeed: %v", err)
	}
	assertBillettholderUserAssociationCount(t, db, expectedAssociation, expectedAssociationCount)
}

func TestDisassociateUsersFromBillettholderEmail_RemovesAssociationWhenNoRemainingEmailMatchesUser(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at ei manuell e-postadresse er fjerna frå ein billettholder, og ingen attverande e-postadresser på billettholderen samsvarer med brukaren.",
		When:  "Når e-postadressa blir forsona mot brukar-tilknytingar.",
		Then:  "Så skal den varige brukar-tilknytinga fjernast.",
	})

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	expectedAssociationCount := 0
	manualEmail := "participant@example.com"
	userEmail := "Participant@Example.com"

	db, slogger := createCheckInTestDB(t)
	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", userEmail)
	removedEmailID := insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, manualEmail)
	insertBillettholderUserAssociation(t, db, expectedAssociation)
	deleteBillettholderEmailByID(t, db, removedEmailID)

	// When
	err := DisassociateUsersFromBillettholderEmail(expectedAssociation.BillettholderID, manualEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected disassociation to succeed: %v", err)
	}
	assertBillettholderUserAssociationCount(t, db, expectedAssociation, expectedAssociationCount)
}

func TestDisassociateUsersFromBillettholderEmail_KeepsAssociationWhenRemainingEmailStillMatchesUser(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at ei manuell e-postadresse er fjerna frå ein billettholder, men ei anna attverande e-postadresse på same billettholder framleis samsvarer med brukaren.",
		When:  "Når e-postadressa blir forsona mot brukar-tilknytingar.",
		Then:  "Så skal den varige brukar-tilknytinga behaldast.",
	})

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	expectedAssociationCount := 1
	removedEmail := "participant@example.com"
	remainingEmail := "PARTICIPANT@example.com"
	userEmail := "Participant@Example.com"

	db, slogger := createCheckInTestDB(t)
	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", userEmail)
	removedEmailID := insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, removedEmail)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, remainingEmail)
	insertBillettholderUserAssociation(t, db, expectedAssociation)
	deleteBillettholderEmailByID(t, db, removedEmailID)

	// When
	err := DisassociateUsersFromBillettholderEmail(expectedAssociation.BillettholderID, removedEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected disassociation cleanup to succeed: %v", err)
	}
	assertBillettholderUserAssociationCount(t, db, expectedAssociation, expectedAssociationCount)
}
