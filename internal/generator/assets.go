package generator

import (
	_ "embed"
)

//go:embed assets/qrcode_ext.png
var templateImageData []byte

//go:embed assets/SF-Pro-Text-Regular.otf
var textFontData []byte

//go:embed assets/barcode39.ttf
var barcodeFontData []byte

