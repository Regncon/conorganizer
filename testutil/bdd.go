package testutil

import (
	"strings"
	"testing"
)

type BDD struct {
	Given string
	When  string
	Then  string
}

func Behavior(t testing.TB, behavior BDD) {
	t.Helper()

	missingFields := []string{}
	if strings.TrimSpace(behavior.Given) == "" {
		missingFields = append(missingFields, "Given")
	}
	if strings.TrimSpace(behavior.When) == "" {
		missingFields = append(missingFields, "When")
	}
	if strings.TrimSpace(behavior.Then) == "" {
		missingFields = append(missingFields, "Then")
	}
	if len(missingFields) > 0 {
		t.Fatalf("BDD behavior metadata is missing %s", strings.Join(missingFields, ", "))
	}
}
