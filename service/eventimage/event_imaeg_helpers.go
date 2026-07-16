package eventimage

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetEventImageUrl(eventID, kind string, eventImageDir *string) string {
	if eventImageDir == nil || *eventImageDir == "" {
		return fmt.Sprintf("/static/placeholder_%s.svg", kind)
	}

	filename := fmt.Sprintf("%s_%s.webp", eventID, kind)
	imagePath := filepath.Join(*eventImageDir, filename)

	if info, err := os.Stat(imagePath); err == nil {
		return fmt.Sprintf("/event-images/%s?v=%d", filename, info.ModTime().UnixNano())
	}
	return fmt.Sprintf("/static/placeholder_%s.svg", kind)
}
