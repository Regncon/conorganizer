package formsubmission

import "strings"

func shouldShowStringValue(value string) string {
	if value != "" {
		return value
	}
	return ""
}

func normalizeTextareaSubmission(value string) string {
	return strings.TrimSpace(strings.ReplaceAll(value, "\r\n", "\n"))
}
