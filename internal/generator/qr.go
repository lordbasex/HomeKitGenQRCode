package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	qrcode "github.com/skip2/go-qrcode"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// scaleCoords scales coordinates based on the scale factor.
// This is used to maintain proper positioning when the template image is scaled.
func scaleCoords(x, y float64, scale float64) (int, int) {
	return int(x * scale), int(y * scale)
}

// loadFontFace loads a font from byte data (OTF or TTF) and returns a font.Face.
// Uses opentype package to properly support OpenType fonts.
func loadFontFace(fontData []byte, size float64) (font.Face, error) {
	ttf, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %w", err)
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating font face: %w", err)
	}

	return face, nil
}

// drawTextWithFace draws text using a font.Face directly on the image.
// The y parameter represents the top of the text (matching gg's behavior).
// This function adjusts the Y position to account for font metrics since
// font.Drawer uses baseline Y, but we want top Y positioning.
func drawTextWithFace(img *image.RGBA, face font.Face, text string, x, y int, clr color.Color) {
	// Get font metrics to adjust Y position
	// font.Drawer uses baseline Y, but we want top Y (like gg)
	metrics := face.Metrics()
	ascent := metrics.Ascent
	// Adjust Y: add ascent to convert from top to baseline
	baselineY := y + int(ascent.Ceil())

	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(baselineY * 64),
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(clr),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
}

// drawScaledText draws text at scaled coordinates using gg (for TTF fonts).
// This is a fallback function when opentype cannot load a font.
func drawScaledText(dc *gg.Context, text string, x, y float64, scale float64) {
	sx, sy := scaleCoords(x, y, scale)
	dc.DrawString(text, float64(sx), float64(sy))
}

// drawScaledTextOTF draws text using OTF font directly on the image.
// This function uses opentype fonts for better quality rendering.
func drawScaledTextOTF(img *image.RGBA, face font.Face, text string, x, y float64, scale float64) {
	sx, sy := scaleCoords(x, y, scale)
	drawTextWithFace(img, face, text, sx, sy, color.Black)
}

// GenerateHomeKitLabel generates a HomeKit QR code label matching the Python implementation.
// This is the main function that creates the complete HomeKit setup label image.
//
// Parameters:
//   - category: HomeKit device category ID
//   - password: Setup password in format XXX-XX-XXX
//   - setupID: Setup ID (4 alphanumeric characters)
//   - mac: MAC address (12 hexadecimal characters)
//   - output: Output image file path (PNG format)
//
// The function:
//  1. Generates the HomeKit setup URI
//  2. Generates device codes, serial numbers, and CSN
//  3. Loads the template image and fonts
//  4. Draws all text elements and barcodes
//  5. Generates and positions the QR code
//  6. Saves the final image
func GenerateHomeKitLabel(category int, password, setupID, mac, output string) error {
	// Generate URI and codes
	uri := GenHomeKitSetupURI(category, password, setupID)
	device := GenerateDeviceCode(category)
	serial := GenerateSerial()
	csn := GenerateCSN()

	// Load base template image from embedded data
	templateImg, _, err := image.Decode(bytes.NewReader(templateImageData))
	if err != nil {
		return fmt.Errorf("error decoding template: %w", err)
	}

	// Calculate scale factor based on template width
	// Base template width is 842 pixels, scale is calculated as templateWidth / 842
	templateWidth := float64(templateImg.Bounds().Dx())
	scale := templateWidth / 842.0

	W := templateImg.Bounds().Dx()
	H := templateImg.Bounds().Dy()

	// Create new image context using gg library
	dc := gg.NewContext(W, H)
	dc.DrawImage(templateImg, 0, 0)

	// Calculate font sizes based on scale factor
	textFontSize := 18 * scale
	barcodeFontSize := 36 * scale
	superscriptFontSize := 8 * scale
	codeFontSize := 28 * scale

	// Load OTF font using opentype for proper OpenType support (from embedded data)
	textFace, err := loadFontFace(textFontData, textFontSize)
	if err != nil {
		return fmt.Errorf("error loading text font: %w", err)
	}

	// Load barcode font (TTF) from embedded data
	barcodeFace, err := loadFontFace(barcodeFontData, barcodeFontSize)
	if err != nil {
		return fmt.Errorf("error loading barcode font: %w", err)
	}

	// Load superscript font for trademark symbol (from embedded data)
	superscriptFace, err := loadFontFace(textFontData, superscriptFontSize)
	if err != nil {
		// Fallback to regular text face if superscript fails
		superscriptFace = textFace
	}

	// Load code font for setup code digits (from embedded data)
	codeFace, err := loadFontFace(textFontData, codeFontSize)
	if err != nil {
		return fmt.Errorf("error loading code font: %w", err)
	}

	// Generate QR code first (before converting to RGBA)
	// Use Medium error correction (matching Python's ERROR_CORRECT_Q)
	qr, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("error generating QR code: %w", err)
	}

	// Calculate final QR code size
	// Increased to 653x653 pixels to better fill the available space while maintaining centering
	qrFinalSize := 653

	// Generate QR code with high resolution first
	// Use a multiple of the final size to ensure we can scale down cleanly
	// QR codes have 25x25 modules (version 2), so we want box_size to be a multiple
	// that when scaled gives us sharp squares
	// Generate at 4x resolution for better quality, then scale down
	boxSizeMultiplier := 4
	boxSize := (qrFinalSize / 25) * boxSizeMultiplier
	if boxSize < boxSizeMultiplier {
		boxSize = boxSizeMultiplier
	}

	// Generate QR image at high resolution
	qrImgHighRes := qr.Image(boxSize)

	// Make QR code background transparent (white pixels become transparent)
	qrImgTransparent := makeQRTransparent(qrImgHighRes)

	// Scale down using NearestNeighbor to preserve sharp square edges
	// NearestNeighbor is perfect for QR codes as it doesn't smooth/blur the edges
	qrImg := resizeQRCode(qrImgTransparent, qrFinalSize, qrFinalSize)

	// Draw QR code directly on RGBA image to ensure it's properly included
	// First convert gg context to RGBA
	baseImg := dc.Image()
	rgbaImg, ok := baseImg.(*image.RGBA)
	if !ok {
		// Convert to RGBA if not already
		bounds := baseImg.Bounds()
		rgbaImg = image.NewRGBA(bounds)
		draw.Draw(rgbaImg, bounds, baseImg, bounds.Min, draw.Src)
	}

	// Now draw QR code on the RGBA image at the correct position
	// Calculate center position based on original QR code position and size
	// Original QR was at (19, 77) with size 136*scale
	// We want to center the new larger QR code (qrFinalSize) at the same center point
	originalQRSize := 136 * scale
	originalQRX := 19 * scale
	originalQRY := 77 * scale
	originalCenterX := originalQRX + originalQRSize/2
	originalCenterY := originalQRY + originalQRSize/2

	// Calculate new position to center the larger QR code
	qrX := int(originalCenterX - float64(qrFinalSize)/2)
	qrY := int(originalCenterY - float64(qrFinalSize)/2)

	qrBounds := qrImg.Bounds()
	// Ensure QR code is drawn at exact size
	qrDestRect := image.Rect(qrX, qrY, qrX+qrFinalSize, qrY+qrFinalSize)
	draw.Draw(rgbaImg, qrDestRect, qrImg, qrBounds.Min, draw.Over)

	// Get category name from reference map
	categoryName := CategoryReference[category]
	if categoryName == "" {
		categoryName = "Unknown"
	}

	// Positioning variables (matching Python code)
	// These values define the layout and spacing of text elements
	y := 6.0
	x := 200.0
	spacingTop := 20.0
	spacingBody := 18.0
	spacingExtra := 6.0
	barcodeSpacing := 30.0

	// Draw header text using OTF font
	headerText := fmt.Sprintf("HomeKit %s | %s | WIFI", categoryName, device)
	drawScaledTextOTF(rgbaImg, textFace, headerText, x, y, scale)
	y += spacingTop

	// Draw brand with superscript trademark symbol
	brand := "Designed by StudioPeters"
	drawScaledTextOTF(rgbaImg, textFace, brand, x, y, scale)

	// Get brand width for superscript positioning
	brandWidth := measureStringWidth(textFace, brand)
	supRX := x*scale + brandWidth
	supRY := y*scale + 3*scale

	// Draw ® with superscript font
	drawTextWithFace(rgbaImg, superscriptFace, "®", int(supRX), int(supRY), color.Black)

	y += spacingTop
	drawScaledTextOTF(rgbaImg, textFace, "Assembled in the Netherlands", x, y, scale)
	y += spacingTop

	// Draw device code
	drawScaledTextOTF(rgbaImg, textFace, fmt.Sprintf("(1P)%s", device), x, y, scale)

	// Draw MAC address if provided
	if mac != "" {
		// Format MAC address (add colons every 2 characters)
		formattedMAC := formatMAC(mac)
		drawScaledTextOTF(rgbaImg, textFace, fmt.Sprintf("MAC: %s", formattedMAC), 560, y, scale)

		// Draw MAC barcode
		drawScaledTextOTF(rgbaImg, barcodeFace, fmt.Sprintf("*%s*", strings.ToUpper(mac)), 560, y+spacingBody+spacingExtra, scale)
	}

	y += spacingBody + spacingExtra

	// Draw device code barcode
	drawScaledTextOTF(rgbaImg, barcodeFace, fmt.Sprintf("*%s*", device), x, y, scale)

	y += barcodeSpacing

	// Draw serial number
	drawScaledTextOTF(rgbaImg, textFace, fmt.Sprintf("(S) Serial No. %s", serial), x, y, scale)
	y += spacingBody + spacingExtra

	// Draw serial barcode
	drawScaledTextOTF(rgbaImg, barcodeFace, fmt.Sprintf("*%s*", serial), x, y, scale)

	y += barcodeSpacing

	// Draw CSN (Customer Serial Number)
	drawScaledTextOTF(rgbaImg, textFace, fmt.Sprintf("CSN %s", csn), x, y, scale)
	y += spacingBody + spacingExtra

	// Draw CSN barcode
	drawScaledTextOTF(rgbaImg, barcodeFace, fmt.Sprintf("*%s*", csn), x, y, scale)

	// Draw setup code digits (password without dashes) using OTF font
	code := strings.ReplaceAll(password, "-", "")

	// Draw code digits in two rows (4 digits per row)
	for i := 0; i < 4; i++ {
		// Top row
		cx, cy := scaleCoords(76+float64(i)*20, 12, scale)
		drawTextWithFace(rgbaImg, codeFace, string(code[i]), cx, cy, color.Black)

		// Bottom row
		cx, cy = scaleCoords(76+float64(i)*20, 39, scale)
		drawTextWithFace(rgbaImg, codeFace, string(code[i+4]), cx, cy, color.Black)
	}

	// Create output directory if needed
	outputDir := filepath.Dir(output)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}
	}

	// Save image (use rgbaImg which has all the text drawn)
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	// Note: Go's png.Encode doesn't support DPI directly
	// For 300 DPI, you would need to use a library that supports it
	// For now, we'll save as PNG
	return png.Encode(out, rgbaImg)
}

// measureStringWidth measures the width of a string using a font face.
// Returns the width in pixels as a float64.
func measureStringWidth(face font.Face, text string) float64 {
	var width fixed.Int26_6
	for _, r := range text {
		advance, ok := face.GlyphAdvance(r)
		if ok {
			width += advance
		}
	}
	return float64(width) / 64.0
}

// formatMAC formats a MAC address string by adding colons every 2 characters.
// Example: "AABBCCDDEEFF" -> "AA:BB:CC:DD:EE:FF"
// If the MAC address is not 12 characters, returns it unchanged.
func formatMAC(mac string) string {
	if len(mac) != 12 {
		return mac
	}
	var parts []string
	for i := 0; i < 12; i += 2 {
		parts = append(parts, mac[i:i+2])
	}
	return strings.ToUpper(strings.Join(parts, ":"))
}

// makeQRTransparent converts white pixels to transparent in QR code image.
// This allows the QR code to blend seamlessly with the template background.
// Pixels with RGB values greater than 200 are considered white and made transparent.
func makeQRTransparent(img image.Image) image.Image {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()

			// If pixel is mostly white (RGB > 200), make it transparent
			if r>>8 > 200 && g>>8 > 200 && b>>8 > 200 {
				rgba.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else {
				rgba.Set(x, y, c)
			}
		}
	}

	return rgba
}

// resizeQRCode resizes a QR code image preserving sharp square edges.
// Uses NearestNeighbor algorithm to avoid smoothing/rounding the square modules.
// This is essential for QR codes to maintain their scannable appearance.
func resizeQRCode(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	// Use NearestNeighbor for QR codes to preserve sharp square edges
	// This prevents the "rounded" appearance that happens with interpolation
	xdraw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}
