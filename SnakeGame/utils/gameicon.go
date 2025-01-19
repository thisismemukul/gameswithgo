package utils

import (
	"bytes"
	_ "embed"
	"image"

	"github.com/thisismemukul/snake/exceptions"
)

//go:embed icon_black.png
var iconImageData []byte

func LoadIconImage() image.Image {
	img, _, err := image.Decode(bytes.NewReader(iconImageData))
	exceptions.CheckErrors(err, "Error While reading icon image")
	return img
}
