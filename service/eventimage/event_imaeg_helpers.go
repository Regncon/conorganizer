package eventimage

import (
	"fmt"
	"os"
)

func GetEventImaegeUrl(eventId string, kind string, eventImageDir *string) string {

	imageFile := fmt.Sprintf("/%s/%s_%s.webp", *eventImageDir, eventId, kind)
	if _, err := os.Stat(fmt.Sprintf(".%s", imageFile)); err == nil {
		return imageFile
	}
	return fmt.Sprintf("/static/placeholder_%s.svg", kind)
}
