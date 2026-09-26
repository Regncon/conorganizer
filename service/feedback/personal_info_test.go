package feedback

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

var personalInfoExamples = []string{
	"kari@example.no",
	"@",
	"12345678",
	"123 45 678",
	"+47 123 45 678",
	"4712345678",
	"010190 12345",
	"12.34.56.78",
}

var harmlessExamples = []string{
	"Rom 101",
	"2026",
	"kl 18:30",
	"1234567",
	"123456789",
	"26.09.2026",
	"1.9.2026",
	"26-09-2026",
	"kl 18.30-19.45",
	"18.30 - 19.45",
	"18.30 19.45",
}

var personalInfoNextToDateExamples = []string{
	"26.09.2026 12345678",
	"12345678 26.09.2026",
	"kl 18.30-19.45 ring 123 45 678",
}

func TestContainsPersonalInfo_FlagsEmailsPhoneNumbersAndIDNumbers(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given texts with an @, a phone number or a fødselsnummer, alone or inside a sentence.",
		When:  "When they are checked for personal info.",
		Then:  "Then every one of them is flagged.",
	})

	for _, example := range personalInfoExamples {
		for _, text := range []string{example, "Kontakt meg: " + example + ", takk"} {
			// Given
			expectedFlagged := true

			// When
			flagged := ContainsPersonalInfo(text)

			// Then
			if flagged != expectedFlagged {
				t.Fatalf("expected %q to be flagged as personal info", text)
			}
		}
	}
}

func TestContainsPersonalInfo_AllowsOtherNumbers(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given texts with room numbers, years, times, dates, time ranges and digit runs of other lengths.",
		When:  "When they are checked for personal info.",
		Then:  "Then none of them are flagged.",
	})

	for _, example := range harmlessExamples {
		for _, text := range []string{example, "Vi møttes i " + example + " og det var fint."} {
			// Given
			expectedFlagged := false

			// When
			flagged := ContainsPersonalInfo(text)

			// Then
			if flagged != expectedFlagged {
				t.Fatalf("expected %q not to be flagged as personal info", text)
			}
		}
	}
}

func TestContainsPersonalInfo_FlagsPhoneNumbersNextToDates(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given texts where a phone number sits right next to a date or a time range.",
		When:  "When they are checked for personal info.",
		Then:  "Then the phone number is still flagged even though the date itself is allowed.",
	})

	for _, text := range personalInfoNextToDateExamples {
		// Given
		expectedFlagged := true

		// When
		flagged := ContainsPersonalInfo(text)

		// Then
		if flagged != expectedFlagged {
			t.Fatalf("expected %q to be flagged as personal info", text)
		}
	}
}
