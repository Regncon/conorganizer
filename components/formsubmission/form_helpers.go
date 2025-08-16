package formsubmission

import (
	"fmt"
)

func shouldShowStringValue(value string) string {
	if value != "" {
		return value
	}
	return ""
}

func shouldShowNumberValue(value int64) string {
	if value != 0 {
		return fmt.Sprintf("%d", value)
	}
	return ""
}
