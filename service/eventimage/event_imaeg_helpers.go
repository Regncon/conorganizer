package eventimage

import (
	"fmt"
	"os"
)

func GetEventImageUrl(eventId string, kind string, eventImageDir *string) string {
	imageFile := fmt.Sprintf("/%s/%s_%s.webp", *eventImageDir, eventId, kind)
	if _, err := os.Stat(fmt.Sprintf(".%s", imageFile)); err == nil {
		imageUrl := fmt.Sprintf("/event-images/%s_%s.webp", eventId, kind)
		return imageUrl
	}
	return fmt.Sprintf("/static/placeholder_%s.svg", kind)
}
