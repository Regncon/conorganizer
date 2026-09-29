// Package qrtest scans QR codes drawn by components/qrcode in tests.
package qrtest

import (
	"image"
	"image/color"
	"regexp"
	"strconv"
	"testing"

	"github.com/makiuchi-d/gozxing"
	zxingqr "github.com/makiuchi-d/gozxing/qrcode"
)

var moduleRe = regexp.MustCompile(`M(\d+) (\d+)h1v1h-1z`)
var viewBoxRe = regexp.MustCompile(`viewBox="0 0 (\d+) (\d+)"`)

// ScanSVG paints every dark module of a QR SVG into an image (10px per
// module) and decodes it, so a test exercises the drawing, not only the
// encoder.
func ScanSVG(t *testing.T, svg string) string {
	t.Helper()
	box := viewBoxRe.FindStringSubmatch(svg)
	if box == nil {
		t.Fatalf("no viewBox in QR SVG: %.200s", svg)
	}
	size, _ := strconv.Atoi(box[1])
	const scale = 10
	img := image.NewGray(image.Rect(0, 0, size*scale, size*scale))
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	modules := moduleRe.FindAllStringSubmatch(svg, -1)
	if len(modules) == 0 {
		t.Fatalf("no modules in QR SVG: %.200s", svg)
	}
	for _, m := range modules {
		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		for dy := range scale {
			for dx := range scale {
				img.SetGray(x*scale+dx, y*scale+dy, color.Gray{Y: 0})
			}
		}
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		t.Fatalf("bitmap: %v", err)
	}
	result, err := zxingqr.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		t.Fatalf("decode QR from SVG: %v", err)
	}
	return result.GetText()
}
