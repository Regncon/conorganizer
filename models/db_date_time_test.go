package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDBDateTimeScan_WhenInputUsesSupportedFormat_ReturnsValidValue(t *testing.T) {
	// Given timestamp inputs in every supported database format,
	// when they are scanned into DBDateTime,
	// then every value is accepted as valid.

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
	// Given a null database timestamp,
	// when it is scanned into DBDateTime,
	// then the value is marked invalid without error.

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
	// Given an invalid timestamp string,
	// when it is scanned into DBDateTime,
	// then parsing fails.

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
	// Given a valid DBDateTime,
	// when it is converted to a driver value,
	// then it returns RFC3339Nano text in UTC.

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
	// Given an invalid DBDateTime,
	// when it is converted to a driver value,
	// then it returns nil.

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
	// Given valid timestamp JSON,
	// when it is unmarshaled and marshaled again,
	// then the same timestamp text is emitted.

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
	// Given null timestamp JSON,
	// when it is unmarshaled into DBDateTime,
	// then the value is marked invalid.

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
