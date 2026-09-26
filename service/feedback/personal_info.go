package feedback

import (
	"regexp"
	"slices"
	"strings"
)

// EmailMarker is treated as an email address wherever it appears. Being strict
// on purpose: any "@" at all is rejected.
const EmailMarker = "@"

// DigitRunPattern matches a run of digits that may be separated by single
// spaces, dots or dashes. It must stay valid in both Go RE2 and JavaScript,
// since the frontend uses the same pattern.
const DigitRunPattern = `\d(?:[ .\-]?\d)*`

// BannedDigitCounts are the digit counts of a DigitRunPattern match that count
// as personal info: Norwegian phone numbers (8), with country code 47 (10) and
// fødselsnummer (11).
var BannedDigitCounts = []int{8, 10, 11}

// DateOrTimeRangePattern matches dates (DD.MM.YYYY or DD-MM-YYYY) and time
// ranges (18.30-19.45 or 18.30 19.45). DigitRunPattern would otherwise merge
// their digits into one 8-digit run and flag them as a phone number, so these
// shapes are blanked out before digits are counted. Like DigitRunPattern it
// must stay valid in both Go RE2 and JavaScript.
const DateOrTimeRangePattern = `\b(?:(?:0?[1-9]|[12]\d|3[01])[.\-](?:0?[1-9]|1[0-2])[.\-]\d{4}|(?:[01]?\d|2[0-3])\.[0-5]\d(?: ?- ?| )(?:[01]?\d|2[0-3])\.[0-5]\d)\b`

// DateOrTimeRangeReplacement replaces every DateOrTimeRangePattern match. It
// is neither a digit nor a DigitRunPattern separator, so it also stops a
// digit run from continuing across the blanked-out date.
const DateOrTimeRangeReplacement = "_"

var (
	digitRun        = regexp.MustCompile(DigitRunPattern)
	dateOrTimeRange = regexp.MustCompile(DateOrTimeRangePattern)
)

// ContainsPersonalInfo reports whether text looks like it contains an email
// address, a phone number or a fødselsnummer.
func ContainsPersonalInfo(text string) bool {
	if strings.Contains(text, EmailMarker) {
		return true
	}
	withoutDates := dateOrTimeRange.ReplaceAllString(text, DateOrTimeRangeReplacement)
	for _, match := range digitRun.FindAllString(withoutDates, -1) {
		if slices.Contains(BannedDigitCounts, countDigits(match)) {
			return true
		}
	}
	return false
}

func countDigits(s string) int {
	count := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			count++
		}
	}
	return count
}
