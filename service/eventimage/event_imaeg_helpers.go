package eventimage

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func GetEventImageUrl(eventID, kind string, eventImageDir *string) string {
	if eventImageDir == nil || *eventImageDir == "" {
		return fmt.Sprintf("/static/placeholder_%s.svg", kind)
	}

	filename := fmt.Sprintf("%s_%s.webp", eventID, kind)
	imagePath := filepath.Join(*eventImageDir, filename)

	info, err := os.Stat(imagePath)
	if err == nil {
		stamp := strconv.FormatInt(info.ModTime().UnixMilli(), 36)
		return fmt.Sprintf("/event-images/%s_%s_%s.webp", eventID, kind, stamp)
	}
	return fmt.Sprintf("/static/placeholder_%s.svg", kind)
}

func FileServer(eventImageDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(eventImageDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
		if name == "" || strings.Contains(name, "/") {
			http.NotFound(w, r)
			return
		}

		if stableName, ok := stableEventImageName(name); ok {
			r = r.Clone(r.Context())
			r.URL.Path = "/" + stableName
		}

		fileServer.ServeHTTP(w, r)
	})
}

func stableEventImageName(name string) (string, bool) {
	if stableName, ok := stablePublicImageName(name); ok {
		return stableName, true
	}
	return stableSourceImageName(name)
}

func stablePublicImageName(name string) (string, bool) {
	if !strings.HasSuffix(name, ".webp") {
		return "", false
	}

	base := strings.TrimSuffix(name, ".webp")
	for _, kind := range []string{"card", "banner"} {
		eventID, stamp, ok := strings.Cut(base, "_"+kind+"_")
		if ok && eventID != "" && stamp != "" {
			return eventID + "_" + kind + ".webp", true
		}
	}

	return "", false
}

func stableSourceImageName(name string) (string, bool) {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	eventID, stamp, ok := strings.Cut(base, "_source_")
	if !ok || eventID == "" || stamp == "" {
		return "", false
	}

	return eventID + "_source" + ext, true
}
