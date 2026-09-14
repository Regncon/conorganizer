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
			contentsString := string(contents)
			if strings.HasPrefix(filepath.Base(url), "terminus-7-etasje-") {
				if !strings.Contains(contentsString, `id="room-`+roomNumber+`-destination-fill"`) {
					t.Fatalf("map %s does not highlight room %s", url, roomNumber)
				}
				return
			}
			if !strings.Contains(contentsString, `data-room-number="`+roomNumber+`"`) {
				t.Fatalf("map %s does not contain room %s", url, roomNumber)
			}
		})
	}
}

func TestFloorMaps_ReferenceExistingSVGAssets(t *testing.T) {
	assetRoot := filepath.Join("..", "..", "static", "rooms")

	for floor, url := range floorMaps {
		t.Run(filepath.Base(url), func(t *testing.T) {
			contents, err := os.ReadFile(filepath.Join(assetRoot, filepath.Base(url)))
			if err != nil {
				t.Fatalf("map asset %s is missing: %v", url, err)
			}
			if !strings.Contains(string(contents), `id="room-targets"`) {
				t.Fatalf("floor %d map %s has no interactive room targets", floor, url)
			}
		})
	}
}
