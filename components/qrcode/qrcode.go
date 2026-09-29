// Package qrcode draws QR codes as SVG, so they stay sharp when printed.
package qrcode

import (
	"errors"
	"fmt"
	"strings"

	goqrcode "github.com/skip2/go-qrcode"
)

// SVG returns content as a QR code: black modules on a white background,
// including the quiet zone scanners need around the code. One viewBox unit is
// one module, so the SVG scales to any print size without blurring.
func SVG(content string) (string, error) {
	if content == "" {
		return "", errors.New("qr code content is empty")
	}
	code, err := goqrcode.New(content, goqrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("encode qr code: %w", err)
	}
	bitmap := code.Bitmap()
	size := len(bitmap)

	var path strings.Builder
	for y, row := range bitmap {
		for x, dark := range row {
			if dark {
				fmt.Fprintf(&path, "M%d %dh1v1h-1z", x, y)
			}
		}
	}

	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" shape-rendering="crispEdges"><rect width="%d" height="%d" fill="#fff"/><path fill="#000" d="%s"/></svg>`,
		size, size, size, size, path.String(),
	), nil
}
