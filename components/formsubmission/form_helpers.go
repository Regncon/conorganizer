package formsubmission

func shouldShowStringValue(value string) string {
	if value != "" {
		return value
	}
	return ""
}
