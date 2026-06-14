package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestFormatBDDComment_PreservesSourceLineBreaksAsReadableSentences(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a BDD comment already split across source comment lines.",
		When:  "When it is formatted for the report.",
		Then:  "Then the report keeps three readable BDD lines.",
	})

	// Given
	expectedLines := []string{
		"Gitt at brukeren er admin.",
		"Når hovednavigasjonen vises.",
		"Så skal adminlenken vises.",
	}
	comment := strings.Join([]string{
		"Gitt at brukeren er admin,",
		"når hovednavigasjonen vises,",
		"så skal adminlenken vises.",
	}, "\n")

	// When
	actualLines := formatBDDComment(comment)

	// Then
	assertLines(t, expectedLines, actualLines)
}

func TestFormatBDDComment_SplitsNorwegianSingleLineComment(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a legacy Norwegian BDD comment stored as one sentence.",
		When:  "When it is formatted for the report.",
		Then:  "Then the BDD clauses are split into separate readable lines.",
	})

	// Given
	expectedLines := []string{
		"Gitt at brukeren ikke er innlogget.",
		"Når hovednavigasjonen vises.",
		"Så skal bare offentlige lenker vises.",
	}
	comment := "Gitt at brukeren ikke er innlogget, når hovednavigasjonen vises, så skal bare offentlige lenker vises."

	// When
	actualLines := formatBDDComment(comment)

	// Then
	assertLines(t, expectedLines, actualLines)
}

func TestFormatBDDComment_SplitsEnglishSingleLineComment(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a legacy English BDD comment stored as one sentence.",
		When:  "When it is formatted for the report.",
		Then:  "Then the BDD clauses are split into separate readable lines.",
	})

	// Given
	expectedLines := []string{
		"Given the process is serving HTTP.",
		"When the health endpoint is requested.",
		"Then it returns a generic OK response.",
	}
	comment := "Given the process is serving HTTP, when the health endpoint is requested, then it returns a generic OK response."

	// When
	actualLines := formatBDDComment(comment)

	// Then
	assertLines(t, expectedLines, actualLines)
}

func TestWriteReport_AddsBlankLineBetweenTests(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given multiple test results with BDD comments.",
		When:  "When the automated behavior report is written.",
		Then:  "Then each test block is separated by a blank line.",
	})

	// Given
	expectedSnippet := strings.Join([]string{
		"- `TestFirst` PASS",
		"",
		"  Given a first behavior.",
		"  When it runs.",
		"  Then it passes.",
		"",
		"- `TestSecond` FAIL (0.01s)",
	}, "\n")
	results := map[string][]testResult{
		"example/package": {
			{Name: "TestSecond", Status: "fail", Elapsed: 0.01},
			{Name: "TestFirst", Status: "pass"},
		},
	}
	comments := map[string]string{
		testKey("example/package", "TestFirst"):  "Given a first behavior,\nwhen it runs,\nthen it passes.",
		testKey("example/package", "TestSecond"): "Given a second behavior,\nwhen it runs,\nthen it fails.",
	}
	var report bytes.Buffer

	// When
	writeReport(&report, results, comments)

	// Then
	if !strings.Contains(report.String(), expectedSnippet) {
		t.Fatalf("expected report to contain formatted test blocks:\n%s\n\nactual report:\n%s", expectedSnippet, report.String())
	}
}

func TestCollectFileBDDComments_PrefersStructuredBehaviorMetadata(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a test with structured BDD metadata and an older fallback BDD comment.",
		When:  "When BDD comments are collected from the source file.",
		Then:  "Then the structured metadata is used for the generated report.",
	})

	// Given
	expectedLines := []string{
		"Gitt at strukturert metadata finnes.",
		"Når rapporten samles.",
		"Så skal strukturert metadata brukes.",
	}
	source := `package sample

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

func TestStructuredBehavior(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at strukturert metadata finnes.",
		When:  "Når rapporten samles.",
		Then:  "Så skal strukturert metadata brukes.",
	})

	// Gitt at gammel kommentar finnes,
	// når rapporten samles,
	// så skal denne ikke brukes.

	// Given
	expected := true

	// When
	actual := true

	// Then
	if actual != expected {
		t.Fatal("expected true")
	}
}
`
	path := filepath.Join(t.TempDir(), "structured_test.go")
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatalf("write test source: %v", err)
	}

	// When
	comments, err := collectFileBDDComments(path)

	// Then
	if err != nil {
		t.Fatalf("collect BDD comments: %v", err)
	}
	actualLines := formatBDDComment(comments["TestStructuredBehavior"])
	assertLines(t, expectedLines, actualLines)
}

func assertLines(t *testing.T, expectedLines []string, actualLines []string) {
	t.Helper()

	if len(actualLines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d: %#v", len(expectedLines), len(actualLines), actualLines)
	}
	for i, expectedLine := range expectedLines {
		if actualLines[i] != expectedLine {
			t.Fatalf("line %d mismatch:\nexpected: %q\nactual:   %q", i, expectedLine, actualLines[i])
		}
	}
}
