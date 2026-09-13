package rooms

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoomMaps_ReferenceExistingSVGAssets(t *testing.T) {
	// Given
	assetRoot := filepath.Join("..", "..", "static", "rooms")

	for roomNumber, url := range roomMaps {
		t.Run(roomNumber, func(t *testing.T) {
			// When
			contents, err := os.ReadFile(filepath.Join(assetRoot, filepath.Base(url)))

			// Then
			if err != nil {
				t.Fatalf("map asset %s is missing: %v", url, err)
			}
			if !strings.Contains(string(contents), `id="room-`+roomNumber+`-destination-fill"`) {
				t.Fatalf("map %s does not highlight room %s", url, roomNumber)
			}
		})
	}
}
