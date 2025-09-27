package utils

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	d := time.Since(t)
	suffix := "ago"
	if d < 0 {
		d = -d
		suffix = "from now"
	}

	switch {
	case d < time.Minute:
		sec := int(d.Seconds())
		if sec <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%d seconds %s", sec, suffix)
	case d < time.Hour:
		min := int(d.Minutes())
		if min == 1 {
			return fmt.Sprintf("1 minute %s", suffix)
		}
		return fmt.Sprintf("%d minutes %s", min, suffix)
	case d < 24*time.Hour:
		hrs := int(d.Hours())
		if hrs == 1 {
			return fmt.Sprintf("1 hour %s", suffix)
		}
		return fmt.Sprintf("%d hours %s", hrs, suffix)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return fmt.Sprintf("1 day %s", suffix)
		}
		return fmt.Sprintf("%d days %s", days, suffix)
	case d < 365*24*time.Hour:
		months := int(d.Hours() / (24 * 30))
		if months <= 1 {
			return fmt.Sprintf("1 month %s", suffix)
		}
		return fmt.Sprintf("%d months %s", months, suffix)
	default:
		years := int(d.Hours() / (24 * 365))
		if years <= 1 {
			return fmt.Sprintf("1 year %s", suffix)
		}
		return fmt.Sprintf("%d years %s", years, suffix)
	}
}
