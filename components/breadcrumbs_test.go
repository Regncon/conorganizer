package components

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestBreadcrumbs_RendersParentLinkCurrentStepAndMobileBackLink(t *testing.T) {
	// Gitt at en side har en overordnet brødsmulesti og en gjeldende side,
	// når brødsmulestien vises,
	// så skal overordnet side være lenke og gjeldende side være tekst.

	// Given
	expectedParentTexts := []string{"Hjem"}
	expectedCurrentTexts := []string{"Arrangement"}
	expectedParentHref := "/"

	paths := []BreadcrumbPath{
		{Name: "Hjem", Url: expectedParentHref},
		{Name: "Arrangement", Url: ""},
	}

	// When
	doc := templtest.Render(t, Breadcrumbs(paths))
	actualParentTexts := templtest.CollectTexts(doc, `.breadcrumb-step[href="/"]`)
	actualCurrentTexts := templtest.CollectTexts(doc, ".breadcrumb-end")
	actualMobileBackLinkVisible := templtest.HasSelector(doc, `.breadcrumb-mobile-return[href="/"]`)
	actualSeparatorVisible := templtest.HasSelector(doc, ".breadcrumb-separator")

	// Then
	if !slices.Equal(expectedParentTexts, actualParentTexts) {
		t.Fatalf("parent breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedParentTexts, actualParentTexts)
	}
	if !slices.Equal(expectedCurrentTexts, actualCurrentTexts) {
		t.Fatalf("current breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedCurrentTexts, actualCurrentTexts)
	}
	if !actualMobileBackLinkVisible {
		t.Fatalf("expected mobile back link to be visible")
	}
	if !actualSeparatorVisible {
		t.Fatalf("expected breadcrumb separator to be visible")
	}
}

func TestBreadcrumbs_WhenOnlyCurrentStep_RendersCurrentStepWithoutNavigationLinks(t *testing.T) {
	// Gitt at brødsmulestien bare har gjeldende side,
	// når brødsmulestien vises,
	// så skal den vise gjeldende side uten navigasjonslenker.

	// Given
	expectedCurrentTexts := []string{"Hjem"}
	expectedNavigationLinkCount := 0

	paths := []BreadcrumbPath{
		{Name: "Hjem", Url: ""},
	}

	// When
	doc := templtest.Render(t, Breadcrumbs(paths))
	actualCurrentTexts := templtest.CollectTexts(doc, ".breadcrumb-end")
	actualNavigationLinkCount := doc.Find("a[href]").Length()
	actualSeparatorVisible := templtest.HasSelector(doc, ".breadcrumb-separator")

	// Then
	if !slices.Equal(expectedCurrentTexts, actualCurrentTexts) {
		t.Fatalf("current breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedCurrentTexts, actualCurrentTexts)
	}
	if actualNavigationLinkCount != expectedNavigationLinkCount {
		t.Fatalf("navigation link count mismatch\nexpected: %d\nactual:   %d", expectedNavigationLinkCount, actualNavigationLinkCount)
	}
	if actualSeparatorVisible {
		t.Fatalf("expected breadcrumb separator to be hidden")
	}
}
