package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type packageInfo struct {
	ImportPath   string
	Dir          string
	TestGoFiles  []string
	XTestGoFiles []string
}

type testEvent struct {
	Action      string  `json:"Action"`
	Package     string  `json:"Package"`
	ImportPath  string  `json:"ImportPath"`
	FailedBuild string  `json:"FailedBuild"`
	Test        string  `json:"Test"`
	Elapsed     float64 `json:"Elapsed"`
	Output      string  `json:"Output"`
}

type testResult struct {
	Name    string
	Status  string
	Elapsed float64
}

type testFailure struct {
	Package string
	Test    string
	Output  string
}

func main() {
	packages, err := listPackages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "list packages: %v\n", err)
		os.Exit(1)
	}

	comments, err := collectBDDComments(packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect BDD metadata: %v\n", err)
		os.Exit(1)
	}

	results, failures, testErr := runTests()
	printReport(results, comments, failures...)

	if testErr != nil {
		fmt.Fprintf(os.Stderr, "go test failed: %v\n", testErr)
		os.Exit(1)
	}
}

func listPackages() ([]packageInfo, error) {
	cmd := exec.Command("go", "list", "-json", "./...")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(output))
	packages := []packageInfo{}
	for {
		var pkg packageInfo
		if err := decoder.Decode(&pkg); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func collectBDDComments(packages []packageInfo) (map[string]string, error) {
	comments := make(map[string]string)
	for _, pkg := range packages {
		files := append([]string{}, pkg.TestGoFiles...)
		files = append(files, pkg.XTestGoFiles...)
		for _, file := range files {
			path := filepath.Join(pkg.Dir, file)
			fileComments, err := collectFileBDDComments(path)
			if err != nil {
				return nil, err
			}
			for testName, comment := range fileComments {
				comments[testKey(pkg.ImportPath, testName)] = comment
			}
		}
	}
	return comments, nil
}

func collectFileBDDComments(path string) (map[string]string, error) {
	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	comments := make(map[string]string)
	for _, declaration := range parsedFile.Decls {
		functionDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok || functionDeclaration.Body == nil || !strings.HasPrefix(functionDeclaration.Name.Name, "Test") {
			continue
		}

		comment := firstBDDTextInFunction(parsedFile.Comments, functionDeclaration)
		if comment != "" {
			comments[functionDeclaration.Name.Name] = comment
		}
	}
	return comments, nil
}

func firstBDDTextInFunction(groups []*ast.CommentGroup, functionDeclaration *ast.FuncDecl) string {
	if structuredBDD := firstStructuredBDDInFunction(functionDeclaration); structuredBDD != "" {
		return structuredBDD
	}

	return firstBDDCommentInFunction(groups, functionDeclaration)
}

func firstStructuredBDDInFunction(functionDeclaration *ast.FuncDecl) string {
	for _, statement := range functionDeclaration.Body.List {
		expressionStatement, ok := statement.(*ast.ExprStmt)
		if !ok {
			continue
		}

		callExpression, ok := expressionStatement.X.(*ast.CallExpr)
		if !ok || !isBehaviorCall(callExpression) {
			continue
		}

		if bdd := bddFromBehaviorCall(callExpression); bdd != "" {
			return bdd
		}
	}

	return ""
}

func isBehaviorCall(callExpression *ast.CallExpr) bool {
	switch function := callExpression.Fun.(type) {
	case *ast.SelectorExpr:
		return function.Sel.Name == "Behavior"
	case *ast.Ident:
		return function.Name == "Behavior"
	default:
		return false
	}
}

func bddFromBehaviorCall(callExpression *ast.CallExpr) string {
	for _, argument := range callExpression.Args {
		compositeLiteral, ok := argument.(*ast.CompositeLit)
		if !ok || !isBDDCompositeLiteral(compositeLiteral) {
			continue
		}

		if bdd := bddFromCompositeLiteral(compositeLiteral); bdd != "" {
			return bdd
		}
	}

	return ""
}

func isBDDCompositeLiteral(compositeLiteral *ast.CompositeLit) bool {
	switch literalType := compositeLiteral.Type.(type) {
	case *ast.SelectorExpr:
		return literalType.Sel.Name == "BDD"
	case *ast.Ident:
		return literalType.Name == "BDD"
	default:
		return false
	}
}

func bddFromCompositeLiteral(compositeLiteral *ast.CompositeLit) string {
	lines := map[string]string{}
	for _, element := range compositeLiteral.Elts {
		keyValue, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := keyValue.Key.(*ast.Ident)
		if !ok {
			continue
		}

		value, ok := keyValue.Value.(*ast.BasicLit)
		if !ok || value.Kind != token.STRING {
			continue
		}

		text, err := strconv.Unquote(value.Value)
		if err != nil {
			continue
		}
		lines[key.Name] = strings.Join(strings.Fields(text), " ")
	}

	given := strings.TrimSpace(lines["Given"])
	when := strings.TrimSpace(lines["When"])
	then := strings.TrimSpace(lines["Then"])
	if given == "" || when == "" || then == "" {
		return ""
	}

	return strings.Join([]string{given, when, then}, "\n")
}

func firstBDDCommentInFunction(groups []*ast.CommentGroup, functionDeclaration *ast.FuncDecl) string {
	commentEnd := functionDeclaration.Body.Rbrace
	if len(functionDeclaration.Body.List) > 0 {
		commentEnd = functionDeclaration.Body.List[0].Pos()
	}

	for _, group := range groups {
		if group.Pos() <= functionDeclaration.Body.Lbrace || group.End() >= commentEnd {
			continue
		}

		lines := normalizeCommentLines(group.Text())
		text := strings.Join(lines, " ")
		if text == "" || isSectionComment(text) {
			continue
		}
		if !looksLikeBDDComment(text) {
			return ""
		}

		return strings.Join(lines, "\n")
	}

	return ""
}

func normalizeCommentLines(text string) []string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line != "" {
			normalized = append(normalized, line)
		}
	}
	return normalized
}

func isSectionComment(text string) bool {
	switch strings.TrimSpace(text) {
	case "Given", "When", "Then":
		return true
	default:
		return false
	}
}

func looksLikeBDDComment(text string) bool {
	lowerText := strings.ToLower(text)
	hasGiven := strings.Contains(lowerText, "given") || strings.Contains(lowerText, "gitt")
	hasWhen := strings.Contains(lowerText, "when") || strings.Contains(lowerText, "når")
	hasThen := strings.Contains(lowerText, "then") || strings.Contains(lowerText, "så")

	return hasGiven && hasWhen && hasThen
}

func runTests() (map[string][]testResult, []testFailure, error) {
	cmd := exec.Command("go", "test", "-json", "./...")
	output, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		output = append(output, exitErr.Stderr...)
	}

	results, failures, parseErr := parseTestOutput(output)
	if parseErr != nil {
		return results, failures, parseErr
	}

	return results, failures, err
}

func parseTestOutput(output []byte) (map[string][]testResult, []testFailure, error) {
	results := make(map[string][]testResult)
	failures := []testFailure{}
	buildOutput := make(map[string]string)
	testOutput := make(map[string]string)
	packageOutput := make(map[string]string)
	packageHasTestFailure := make(map[string]bool)

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue
		}

		switch event.Action {
		case "build-output":
			buildOutput[event.ImportPath] += event.Output
		case "build-fail":
			failures = append(failures, testFailure{
				Package: event.ImportPath,
				Output:  buildOutput[event.ImportPath],
			})
		case "output":
			if event.Test != "" {
				testOutput[testKey(event.Package, event.Test)] += event.Output
			} else if event.Package != "" {
				packageOutput[event.Package] += event.Output
			}
		case "pass", "fail", "skip":
			if event.Test == "" {
				if event.Action == "fail" && event.Package != "" && event.FailedBuild == "" && !packageHasTestFailure[event.Package] {
					failures = append(failures, testFailure{
						Package: event.Package,
						Output:  packageOutput[event.Package],
					})
				}
				continue
			}
			if event.Action == "fail" {
				packageHasTestFailure[event.Package] = true
				failures = append(failures, testFailure{
					Package: event.Package,
					Test:    event.Test,
					Output:  testOutput[testKey(event.Package, event.Test)],
				})
			}
			if strings.Contains(event.Test, "/") {
				continue
			}
			results[event.Package] = append(results[event.Package], testResult{
				Name:    event.Test,
				Status:  event.Action,
				Elapsed: event.Elapsed,
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return results, failures, err
	}

	return results, failures, nil
}

func printReport(results map[string][]testResult, comments map[string]string, failures ...testFailure) {
	writeReport(os.Stdout, results, comments, failures...)
}

func writeReport(writer io.Writer, results map[string][]testResult, comments map[string]string, failures ...testFailure) {
	packages := make([]string, 0, len(results))
	totalTests := 0
	failedTests := 0
	skippedTests := 0
	missingBDDMetadata := 0
	failedPackages := make(map[string]struct{})
	for _, failure := range failures {
		failedPackages[failure.Package] = struct{}{}
	}

	for packagePath, packageResults := range results {
		if len(packageResults) == 0 {
			continue
		}
		packages = append(packages, packagePath)
		totalTests += len(packageResults)
		for _, result := range packageResults {
			if result.Status == "fail" {
				failedTests++
			}
			if result.Status == "skip" {
				skippedTests++
			}
			if comments[testKey(packagePath, result.Name)] == "" {
				missingBDDMetadata++
			}
		}
	}
	sort.Strings(packages)

	fmt.Fprintln(writer, "# Automated Behavior Test Report")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Generated from `go test -json ./...` and the first structured BDD metadata or BDD comment in each `Test...` function.")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "## Summary")
	fmt.Fprintf(writer, "- Packages with tests: %d\n", len(packages))
	fmt.Fprintf(writer, "- Tests run: %d\n", totalTests)
	fmt.Fprintf(writer, "- Failed tests: %d\n", failedTests)
	fmt.Fprintf(writer, "- Failed packages: %d\n", len(failedPackages))
	fmt.Fprintf(writer, "- Skipped: %d\n", skippedTests)
	fmt.Fprintf(writer, "- Tests missing BDD metadata: %d\n", missingBDDMetadata)
	fmt.Fprintln(writer)

	for _, packagePath := range packages {
		packageResults := results[packagePath]
		sort.Slice(packageResults, func(i, j int) bool {
			return packageResults[i].Name < packageResults[j].Name
		})

		fmt.Fprintf(writer, "## %s\n\n", packagePath)
		for _, result := range packageResults {
			fmt.Fprintf(writer, "- `%s` %s", result.Name, strings.ToUpper(result.Status))
			if result.Elapsed > 0 {
				fmt.Fprintf(writer, " (%.2fs)", result.Elapsed)
			}
			fmt.Fprintln(writer)
			fmt.Fprintln(writer)

			comment := comments[testKey(packagePath, result.Name)]
			if comment == "" {
				comment = "BDD-metadata mangler."
			}
			for _, line := range formatBDDComment(comment) {
				fmt.Fprintf(writer, "  %s\n", line)
			}
			fmt.Fprintln(writer)
		}
	}

	if len(failures) == 0 {
		return
	}

	sortedFailures := append([]testFailure(nil), failures...)
	sort.Slice(sortedFailures, func(i, j int) bool {
		if sortedFailures[i].Package == sortedFailures[j].Package {
			return sortedFailures[i].Test < sortedFailures[j].Test
		}
		return sortedFailures[i].Package < sortedFailures[j].Package
	})

	fmt.Fprintln(writer, "## Failures")
	fmt.Fprintln(writer)
	for _, failure := range sortedFailures {
		if failure.Test == "" {
			fmt.Fprintf(writer, "### %s\n\n", failure.Package)
		} else {
			fmt.Fprintf(writer, "### %s / %s\n\n", failure.Package, failure.Test)
		}
		fmt.Fprintln(writer, "```text")
		fmt.Fprintln(writer, strings.TrimSpace(failure.Output))
		fmt.Fprintln(writer, "```")
		fmt.Fprintln(writer)
	}
}

func testKey(packagePath string, testName string) string {
	return packagePath + "." + testName
}

var bddSplitPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^(.+?),\s*(when\s+.+?),\s*(then\s+.+)$`),
	regexp.MustCompile(`(?i)^(.+?),\s*(når\s+.+?),\s*(så\s+.+)$`),
}

func formatBDDComment(comment string) []string {
	lines := normalizeCommentLines(comment)
	if len(lines) == 1 {
		if splitLines, ok := splitSingleLineBDDComment(lines[0]); ok {
			lines = splitLines
		}
	}

	for i, line := range lines {
		lines[i] = cleanBDDReportLine(line)
	}

	return lines
}

func splitSingleLineBDDComment(comment string) ([]string, bool) {
	for _, pattern := range bddSplitPatterns {
		matches := pattern.FindStringSubmatch(comment)
		if len(matches) == 4 {
			return []string{matches[1], matches[2], matches[3]}, true
		}
	}

	return nil, false
}

func cleanBDDReportLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ",")
	line = strings.TrimSpace(line)
	line = capitalizeFirstRune(line)
	if line == "" || strings.HasSuffix(line, ".") || strings.HasSuffix(line, "!") || strings.HasSuffix(line, "?") {
		return line
	}
	return line + "."
}

func capitalizeFirstRune(text string) string {
	if text == "" {
		return text
	}

	runes := []rune(text)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
