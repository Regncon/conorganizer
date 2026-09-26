package userctx

import (
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestLoginHrefWithNeste_EncodesPathAndQueryAsReturnTarget(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en forespørsel til en bestemt side med spørrestreng.",
		When:  "Når lenken til innlogging bygges for en ikke-innlogget bruker.",
		Then:  "Så skal lenken inneholde en url-enkodet 'neste'-parameter som peker tilbake dit.",
	})

	// Given
	expectedHref := "/auth?neste=%2Ftilbakemelding%3Fom%3Dfestivalen"
	request := httptest.NewRequest("GET", "/tilbakemelding?om=festivalen", nil)

	// When
	actualHref := loginHrefWithNeste(request)

	// Then
	if actualHref != expectedHref {
		t.Fatalf("expected href %q, got %q", expectedHref, actualHref)
	}
}

func TestLoginHrefWithNeste_PlainRootPathOmitsNesteParam(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en forespørsel til forsiden.",
		When:  "Når lenken til innlogging bygges.",
		Then:  "Så trengs det ingen 'neste'-parameter, ettersom forsiden er standard.",
	})

	// Given
	expectedHref := "/auth"
	request := httptest.NewRequest("GET", "/", nil)

	// When
	actualHref := loginHrefWithNeste(request)

	// Then
	if actualHref != expectedHref {
		t.Fatalf("expected href %q, got %q", expectedHref, actualHref)
	}
}
