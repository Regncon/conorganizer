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
	"sort"
	"strings"
)

type packageInfo struct {
	ImportPath   string
	Dir          string
	TestGoFiles  []string
	XTestGoFiles []string
}

type testEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

type testResult struct {
	Name    string
	Status  string
	Elapsed float64
}

func main() {
	packages, err := listPackages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "list packages: %v\n", err)
		os.Exit(1)
	}

	comments, err := collectBDDComments(packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect BDD comments: %v\n", err)
		os.Exit(1)
	}

	results, testErr := runTests()
	printReport(results, comments)

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

		comment := firstBDDCommentInFunction(parsedFile.Comments, functionDeclaration)
		if comment != "" {
			comments[functionDeclaration.Name.Name] = comment
		}
	}
	return comments, nil
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

		text := normalizeCommentText(group.Text())
		if text == "" || isSectionComment(text) {
			continue
		}
		if !looksLikeBDDComment(text) {
			return ""
		}

		return text
	}

	return ""
}

func normalizeCommentText(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line != "" {
			normalized = append(normalized, line)
		}
	}
	return strings.Join(normalized, " ")
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

func runTests() (map[string][]testResult, error) {
	cmd := exec.Command("go", "test", "-json", "./...")
	output, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		output = append(output, exitErr.Stderr...)
	}

	results := make(map[string][]testResult)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		var event testEvent
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &event); unmarshalErr != nil {
			continue
		}
		if event.Test == "" || strings.Contains(event.Test, "/") {
			continue
		}
		if event.Action != "pass" && event.Action != "fail" && event.Action != "skip" {
			continue
		}
		results[event.Package] = append(results[event.Package], testResult{
			Name:    event.Test,
			Status:  event.Action,
			Elapsed: event.Elapsed,
		})
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return results, scanErr
	}

	return results, err
}

func printReport(results map[string][]testResult, comments map[string]string) {
	packages := make([]string, 0, len(results))
	totalTests := 0
	failedTests := 0
	skippedTests := 0
	missingBDDComments := 0

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
				missingBDDComments++
			}
		}
	}
	sort.Strings(packages)

	fmt.Println("# Automated Behavior Test Report")
	fmt.Println()
	fmt.Println("Generated from `go test -json ./...` and the first BDD comment in each `Test...` function.")
	fmt.Println()
	fmt.Println("## Summary")
	fmt.Printf("- Packages with tests: %d\n", len(packages))
	fmt.Printf("- Tests run: %d\n", totalTests)
	fmt.Printf("- Failed: %d\n", failedTests)
	fmt.Printf("- Skipped: %d\n", skippedTests)
	fmt.Printf("- Tests missing BDD comments: %d\n", missingBDDComments)
	fmt.Println()

	for _, packagePath := range packages {
		packageResults := results[packagePath]
		sort.Slice(packageResults, func(i, j int) bool {
			return packageResults[i].Name < packageResults[j].Name
		})

		fmt.Printf("## %s\n", packagePath)
		for _, result := range packageResults {
			fmt.Printf("- `%s` %s", result.Name, strings.ToUpper(result.Status))
			if result.Elapsed > 0 {
				fmt.Printf(" (%.2fs)", result.Elapsed)
			}
			fmt.Println()

			comment := comments[testKey(packagePath, result.Name)]
			if comment == "" {
				comment = "BDD-kommentar mangler."
			}
			fmt.Printf("  %s\n", comment)
		}
		fmt.Println()
	}
}

func testKey(packagePath string, testName string) string {
	return packagePath + "." + testName
}
