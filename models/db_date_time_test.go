package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestDBDateTimeScan_WhenInputUsesSupportedFormat_ReturnsValidValue(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given timestamp inputs in every supported database format.",
		When:  "When they are scanned into DBDateTime.",
		Then:  "Then every value is accepted as valid.",
	})

	// Given
	expectedValid := true
	inputs := []string{
		"2026-05-17T18:42:13.123Z",
		"2026-05-17T18:42:13Z",
		"2026-05-17 18:42:13",
		"2026-05-17",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			// When
			var actual DBDateTime
			err := actual.Scan(input)

			// Then
			if err != nil {
				t.Fatalf("expected scan to succeed: %v", err)
			}
			if actual.Valid != expectedValid {
				t.Fatalf("valid flag mismatch\nexpected: %v\nactual:   %v", expectedValid, actual.Valid)
			}
		})
	}
}

func TestDBDateTimeScan_WhenInputIsNull_ReturnsInvalidValue(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a null database timestamp.",
		When:  "When it is scanned into DBDateTime.",
		Then:  "Then the value is marked invalid without error.",
	})

	// Given
	expectedValid := false

	// When
	var actual DBDateTime
	err := actual.Scan(nil)

	// Then
	if err != nil {
		t.Fatalf("expected null scan to succeed: %v", err)
	}
	if actual.Valid != expectedValid {
		t.Fatalf("valid flag mismatch\nexpected: %v\nactual:   %v", expectedValid, actual.Valid)
	}
}

func TestDBDateTimeScan_WhenInputIsInvalid_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an invalid timestamp string.",
		When:  "When it is scanned into DBDateTime.",
		Then:  "Then parsing fails.",
	})

	// Given
	expectedError := true

	// When
	var actual DBDateTime
	err := actual.Scan("not a timestamp")
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestDBDateTimeValue_WhenValid_ReturnsRFC3339NanoUTCText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a valid DBDateTime.",
		When:  "When it is converted to a driver value.",
		Then:  "Then it returns RFC3339Nano text in UTC.",
	})

	// Given
	expectedValue := "2026-05-17T18:42:13.123Z"
	parsed := mustParseTime(t, expectedValue)

	// When
	actualValue, err := NewDBDateTime(parsed).Value()

	// Then
	if err != nil {
		t.Fatalf("expected value conversion to succeed: %v", err)
	}
	if actualValue != expectedValue {
		t.Fatalf("value mismatch\nexpected: %v\nactual:   %v", expectedValue, actualValue)
	}
}

func TestDBDateTimeValue_WhenInvalid_ReturnsNil(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an invalid DBDateTime.",
		When:  "When it is converted to a driver value.",
		Then:  "Then it returns nil.",
	})

	// Given
	var expectedValue any

	// When
	actualValue, err := (DBDateTime{}).Value()

	// Then
	if err != nil {
		t.Fatalf("expected null value conversion to succeed: %v", err)
	}
	if actualValue != expectedValue {
		t.Fatalf("value mismatch\nexpected: %v\nactual:   %v", expectedValue, actualValue)
	}
}

func TestDBDateTimeJSON_WhenValid_RoundTripsTimestampText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given valid timestamp JSON.",
		When:  "When it is unmarshaled and marshaled again.",
		Then:  "Then the same timestamp text is emitted.",
	})

	// Given
	expectedJSON := `"2026-05-17T18:42:13Z"`

	// When
	var actual DBDateTime
	unmarshalErr := json.Unmarshal([]byte(expectedJSON), &actual)
	actualJSON, marshalErr := json.Marshal(actual)

	// Then
	if unmarshalErr != nil {
		t.Fatalf("expected JSON unmarshal to succeed: %v", unmarshalErr)
	}
	if !actual.Valid {
		t.Fatalf("expected unmarshaled DBDateTime to be valid")
	}
	if marshalErr != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", marshalErr)
	}
	if string(actualJSON) != expectedJSON {
		t.Fatalf("JSON mismatch\nexpected: %s\nactual:   %s", expectedJSON, actualJSON)
	}
}

func TestDBDateTimeJSON_WhenInputIsNull_ReturnsInvalidValue(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given null timestamp JSON.",
		When:  "When it is unmarshaled into DBDateTime.",
		Then:  "Then the value is marked invalid.",
	})

	// Given
	expectedValid := false

	// When
	var actual DBDateTime
	err := json.Unmarshal([]byte(`null`), &actual)

	// Then
	if err != nil {
		t.Fatalf("expected JSON null unmarshal to succeed: %v", err)
	}
	if actual.Valid != expectedValid {
		t.Fatalf("valid flag mismatch\nexpected: %v\nactual:   %v", expectedValid, actual.Valid)
	}
}

func mustParseTime(t testing.TB, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("failed to parse test time: %v", err)
	}
	return parsed
}
