package feedback

import (
	"regexp"
	"strings"
)

// EmailMarker is treated as an email address wherever it appears. Being strict
// on purpose: any "@" at all is rejected.
const EmailMarker = "@"

// DigitRunPattern matches a run of digits that may be separated by any number of
// spaces, dots, dashes, slashes or brackets, so "342  12 362", "342/12/362" and
// "(+47) 342 12 362" are one run. It must stay valid in both Go RE2 and
// JavaScript, since the frontend uses the same pattern.
const DigitRunPattern = `\d(?:[ \t.\-/()]*\d)*`

// PhoneNumberDigits is the length of a Norwegian phone number. A digit run of
// exactly this many digits counts as personal info.
const PhoneNumberDigits = 8

// MinLongNumberDigits is the length from which any digit run counts as personal
// info: a phone number with a country code ("+47", "0047"), a fødselsnummer
// (11), or a longer run with a number inside. Nine digits (for example an
// organisation number) are allowed.
const MinLongNumberDigits = 10

// IsBannedDigitCount reports whether a digit run with count digits counts as
// personal info.
func IsBannedDigitCount(count int) bool {
	return count == PhoneNumberDigits || count >= MinLongNumberDigits
}

// DateOrTimeRangePattern matches dates (DD.MM.YYYY, DD-MM-YYYY or DD/MM/YYYY) and time
// ranges (18.30-19.45 or 18.30 19.45). DigitRunPattern would otherwise merge
// their digits into one 8-digit run and flag them as a phone number, so these
// shapes are blanked out before digits are counted. Like DigitRunPattern it
// must stay valid in both Go RE2 and JavaScript.
const DateOrTimeRangePattern = `\b(?:(?:0?[1-9]|[12]\d|3[01])[.\-/](?:0?[1-9]|1[0-2])[.\-/]\d{4}|(?:[01]?\d|2[0-3])\.[0-5]\d(?: ?- ?| )(?:[01]?\d|2[0-3])\.[0-5]\d)\b`

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
		if IsBannedDigitCount(countDigits(match)) {
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
