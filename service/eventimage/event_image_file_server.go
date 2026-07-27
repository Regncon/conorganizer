package eventimage

import (
	"fmt"
	"net/http"
	"time"
)

func FileServer(root string) http.Handler {
	files := http.FileServer(http.Dir(root))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setEventImageETag(w, root, r.URL.Path)
		files.ServeHTTP(w, r)
	})
}

func setEventImageETag(w http.ResponseWriter, root string, requestPath string) {
	file, err := http.Dir(root).Open(requestPath)
	if err != nil {
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		return
	}

	w.Header().Set("ETag", eventImageETag(info.ModTime(), info.Size()))
}

func eventImageETag(modTime time.Time, size int64) string {
	return fmt.Sprintf(`W/"%x-%x"`, modTime.UnixNano(), size)
}
