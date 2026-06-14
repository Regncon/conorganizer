package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatBDDComment_PreservesSourceLineBreaksAsReadableSentences(t *testing.T) {
	// Given a BDD comment already split across source comment lines,
	// when it is formatted for the report,
	// then the report keeps three readable BDD lines.

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
	// Given a legacy Norwegian BDD comment stored as one sentence,
	// when it is formatted for the report,
	// then the BDD clauses are split into separate readable lines.

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
	// Given a legacy English BDD comment stored as one sentence,
	// when it is formatted for the report,
	// then the BDD clauses are split into separate readable lines.

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
	// Given multiple test results with BDD comments,
	// when the automated behavior report is written,
	// then each test block is separated by a blank line.

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
