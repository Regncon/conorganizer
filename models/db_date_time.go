package models

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	_ sql.Scanner   = (*DBDateTime)(nil)
	_ driver.Valuer = DBDateTime{}
)

const DBDateTimeNowSQL = "strftime('%Y-%m-%dT%H:%M:%fZ', 'now')"

type DBDateTime struct {
	Time  time.Time
	Valid bool
}

func NewDBDateTime(value time.Time) DBDateTime {
	if value.IsZero() {
		return DBDateTime{}
	}
	return DBDateTime{Time: value, Valid: true}
}

func (dbDateTime DBDateTime) Format(layout string) string {
	if !dbDateTime.Valid {
		return ""
	}
	return dbDateTime.Time.Format(layout)
}

func (dbDateTime DBDateTime) IsZero() bool {
	return !dbDateTime.Valid || dbDateTime.Time.IsZero()
}

func (dbDateTime DBDateTime) TimeOrZero() time.Time {
	if !dbDateTime.Valid {
		return time.Time{}
	}
	return dbDateTime.Time
}

func (dbDateTime DBDateTime) Value() (driver.Value, error) {
	if !dbDateTime.Valid {
		return nil, nil
	}
	return dbDateTime.Time.UTC().Format(time.RFC3339Nano), nil
}

func (dbDateTime *DBDateTime) Scan(value any) error {
	if value == nil {
		*dbDateTime = DBDateTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*dbDateTime = NewDBDateTime(v)
		return nil
	case string:
		return dbDateTime.scanString(v)
	case []byte:
		return dbDateTime.scanString(string(v))
	default:
		return fmt.Errorf("cannot scan %T into DBDateTime", value)
	}
}

func (dbDateTime DBDateTime) MarshalJSON() ([]byte, error) {
	if !dbDateTime.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(dbDateTime.Time.UTC().Format(time.RFC3339Nano))
}

func (dbDateTime *DBDateTime) UnmarshalJSON(value []byte) error {
	if bytes.Equal(value, []byte("null")) {
		*dbDateTime = DBDateTime{}
		return nil
	}

	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return fmt.Errorf("unmarshal DBDateTime: %w", err)
	}
	if strings.TrimSpace(s) == "" {
		*dbDateTime = DBDateTime{}
		return nil
	}
	return dbDateTime.scanString(s)
}

func (dbDateTime *DBDateTime) scanString(value string) error {
	parsed, err := parseDBDateTime(value)
	if err != nil {
		return err
	}
	*dbDateTime = NewDBDateTime(parsed)
	return nil
}

func parseDBDateTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty DBDateTime value")
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported DBDateTime format %q", value)
}
