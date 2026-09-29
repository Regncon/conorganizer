package qrcode

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/qrtest"
)

func TestSVG_ScansBackToTheEncodedURL(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en absolutt URL til tilbakemeldingsskjemaet.",
		When:  "Når den tegnes som QR-kode i SVG og SVG-en rastreres og skannes.",
		Then:  "Så leser skanneren nøyaktig den samme URL-en.",
	})

	// Given
	expectedURL := "https://regncon.no/tilbakemelding?om=festivalen"

	// When
	svg, err := SVG(expectedURL)

	// Then
	if err != nil {
		t.Fatalf("expected SVG to succeed: %v", err)
	}
	if got := qrtest.ScanSVG(t, svg); got != expectedURL {
		t.Fatalf("scanned %q, want %q", got, expectedURL)
	}
}

func TestSVG_EmptyContentIsAnError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt tomt innhold.",
		When:  "Når det tegnes som QR-kode.",
		Then:  "Så returneres en feil i stedet for en tom kode.",
	})

	// Given
	expectedErr := true

	// When
	_, err := SVG("")

	// Then
	if (err != nil) != expectedErr {
		t.Fatalf("expected an error for empty content, got %v", err)
	}
}
