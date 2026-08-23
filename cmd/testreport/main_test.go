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

func TestParseTestOutput_CapturesPackageBuildFailure(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given Go test JSON containing a package build failure.",
		When:  "When the test output is parsed.",
		Then:  "Then the compiler diagnostics are retained for that package.",
	})

	// Given
	expectedPackage := "example/package"
	expectedOutput := "example.go:12:3: undefined: missingName"
	rawOutput := strings.Join([]string{
		`{"ImportPath":"example/package","Action":"build-output","Output":"# example/package\n"}`,
		`{"ImportPath":"example/package","Action":"build-output","Output":"example.go:12:3: undefined: missingName\n"}`,
		`{"ImportPath":"example/package","Action":"build-fail"}`,
	}, "\n")

	// When
	_, failures, err := parseTestOutput([]byte(rawOutput))

	// Then
	if err != nil {
		t.Fatalf("parse test output: %v", err)
	}
	if len(failures) != 1 {
		t.Fatalf("expected one build failure, got %d: %#v", len(failures), failures)
	}
	if failures[0].Package != expectedPackage {
		t.Fatalf("package mismatch\nexpected: %q\nactual:   %q", expectedPackage, failures[0].Package)
	}
	if !strings.Contains(failures[0].Output, expectedOutput) {
		t.Fatalf("expected compiler output to contain %q\nactual: %q", expectedOutput, failures[0].Output)
	}
}

func TestParseTestOutput_CapturesFailedTestDiagnostics(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given Go test JSON containing a failed test and its assertion output.",
		When:  "When the test output is parsed.",
		Then:  "Then the failed test retains the assertion diagnostics.",
	})

	// Given
	expectedPackage := "example/package"
	expectedTest := "TestExample"
	expectedOutput := "expected true, got false"
	rawOutput := strings.Join([]string{
		`{"Action":"run","Package":"example/package","Test":"TestExample"}`,
		`{"Action":"output","Package":"example/package","Test":"TestExample","Output":"=== RUN   TestExample\n"}`,
		`{"Action":"output","Package":"example/package","Test":"TestExample","Output":"    example_test.go:12: expected true, got false\n"}`,
		`{"Action":"fail","Package":"example/package","Test":"TestExample","Elapsed":0.01}`,
		`{"Action":"fail","Package":"example/package","Elapsed":0.01}`,
	}, "\n")

	// When
	_, failures, err := parseTestOutput([]byte(rawOutput))

	// Then
	if err != nil {
		t.Fatalf("parse test output: %v", err)
	}
	if len(failures) != 1 {
		t.Fatalf("expected one test failure, got %d: %#v", len(failures), failures)
	}
	if failures[0].Package != expectedPackage {
		t.Fatalf("package mismatch\nexpected: %q\nactual:   %q", expectedPackage, failures[0].Package)
	}
	if failures[0].Test != expectedTest {
		t.Fatalf("test mismatch\nexpected: %q\nactual:   %q", expectedTest, failures[0].Test)
	}
	if !strings.Contains(failures[0].Output, expectedOutput) {
		t.Fatalf("expected test output to contain %q\nactual: %q", expectedOutput, failures[0].Output)
	}
}

func TestParseTestOutput_CapturesPackageFailureWithoutFailedTest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given Go test JSON containing a package-level failure outside a test function.",
		When:  "When the test output is parsed.",
		Then:  "Then the package-level diagnostics are retained.",
	})

	// Given
	expectedPackage := "example/package"
	expectedOutput := "panic: failure in TestMain"
	rawOutput := strings.Join([]string{
		`{"Action":"output","Package":"example/package","Output":"panic: failure in TestMain\n"}`,
		`{"Action":"fail","Package":"example/package","Elapsed":0.01}`,
	}, "\n")

	// When
	_, failures, err := parseTestOutput([]byte(rawOutput))

	// Then
	if err != nil {
		t.Fatalf("parse test output: %v", err)
	}
	if len(failures) != 1 {
		t.Fatalf("expected one package failure, got %d: %#v", len(failures), failures)
	}
	if failures[0].Package != expectedPackage {
		t.Fatalf("package mismatch\nexpected: %q\nactual:   %q", expectedPackage, failures[0].Package)
	}
	if !strings.Contains(failures[0].Output, expectedOutput) {
		t.Fatalf("expected package output to contain %q\nactual: %q", expectedOutput, failures[0].Output)
	}
}

func TestWriteReport_IncludesPackageFailureDiagnostics(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a test run with a package build failure and no failed test function.",
		When:  "When the automated behavior report is written.",
		Then:  "Then the summary and failure section expose the compiler diagnostics.",
	})

	// Given
	expectedSnippets := []string{
		"- Failed tests: 0",
		"- Failed packages: 1",
		"## Failures",
		"### example/package",
		"example.go:12:3: undefined: missingName",
	}
	failure := testFailure{
		Package: "example/package",
		Output: strings.Join([]string{
			"# example/package",
			"example.go:12:3: undefined: missingName",
		}, "\n"),
	}
	var report bytes.Buffer

	// When
	writeReport(&report, map[string][]testResult{}, map[string]string{}, failure)

	// Then
	for _, expectedSnippet := range expectedSnippets {
		if !strings.Contains(report.String(), expectedSnippet) {
			t.Fatalf("expected report to contain %q\nactual report:\n%s", expectedSnippet, report.String())
		}
	}
}

func TestRunTests_CapturesPackageBuildFailure(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a Go package that does not compile.",
		When:  "When the reporter runs the package tests.",
		Then:  "Then it returns the compiler diagnostics with the failed command status.",
	})

	// Given
	expectedPackage := "example.com/broken"
	expectedOutput := "undefined: missingName"
	moduleDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module example.com/broken\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write temporary go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "broken.go"), []byte("package broken\n\nfunc broken() { missingName() }\n"), 0o644); err != nil {
		t.Fatalf("write broken Go source: %v", err)
	}
	t.Chdir(moduleDir)

	// When
	_, failures, err := runTests()

	// Then
	if err == nil {
		t.Fatal("expected go test to fail")
	}
	if len(failures) != 1 {
		t.Fatalf("expected one build failure, got %d: %#v", len(failures), failures)
	}
	if failures[0].Package != expectedPackage {
		t.Fatalf("package mismatch\nexpected: %q\nactual:   %q", expectedPackage, failures[0].Package)
	}
	if !strings.Contains(failures[0].Output, expectedOutput) {
		t.Fatalf("expected compiler output to contain %q\nactual: %q", expectedOutput, failures[0].Output)
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
