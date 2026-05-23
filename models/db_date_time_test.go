package models

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"
)

func TestDBDateTimeScanFormats(t *testing.T) {
	tests := []string{
		"2026-05-17T18:42:13.123Z",
		"2026-05-17T18:42:13Z",
		"2026-05-17 18:42:13",
		"2026-05-17",
	}

	for _, input := range tests {
		var got DBDateTime
		if err := got.Scan(input); err != nil {
			t.Fatalf("Scan(%q) returned error: %v", input, err)
		}
		if !got.Valid {
			t.Fatalf("Scan(%q) produced invalid DBDateTime", input)
		}
	}
}

func TestDBDateTimeScanNull(t *testing.T) {
	var got DBDateTime
	if err := got.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) returned error: %v", err)
	}
	if got.Valid {
		t.Fatalf("Scan(nil) produced valid DBDateTime")
	}
}

func TestDBDateTimeScanInvalid(t *testing.T) {
	var got DBDateTime
	if err := got.Scan("not a timestamp"); err == nil {
		t.Fatalf("expected invalid timestamp error")
	}
}

func TestDBDateTimeValue(t *testing.T) {
	parsed, err := time.Parse(time.RFC3339Nano, "2026-05-17T18:42:13.123Z")
	if err != nil {
		t.Fatalf("parse test timestamp: %v", err)
	}

	value, err := NewDBDateTime(parsed).Value()
	if err != nil {
		t.Fatalf("Value returned error: %v", err)
	}
	if value != driver.Value("2026-05-17T18:42:13.123Z") {
		t.Fatalf("Value = %v, want RFC3339Nano UTC text", value)
	}

	nullValue, err := (DBDateTime{}).Value()
	if err != nil {
		t.Fatalf("null Value returned error: %v", err)
	}
	if nullValue != nil {
		t.Fatalf("null Value = %v, want nil", nullValue)
	}
}

func TestDBDateTimeJSON(t *testing.T) {
	var got DBDateTime
	if err := json.Unmarshal([]byte(`"2026-05-17T18:42:13Z"`), &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if !got.Valid {
		t.Fatalf("Unmarshal produced invalid DBDateTime")
	}

	out, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(out) != `"2026-05-17T18:42:13Z"` {
		t.Fatalf("Marshal = %s", out)
	}

	if err := json.Unmarshal([]byte(`null`), &got); err != nil {
		t.Fatalf("Unmarshal null returned error: %v", err)
	}
	if got.Valid {
		t.Fatalf("Unmarshal null produced valid DBDateTime")
	}
}
