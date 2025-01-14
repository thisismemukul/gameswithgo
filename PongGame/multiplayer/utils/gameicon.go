package utils

import (
	"bytes"
	_ "embed"
	"image"
	"log"
	"os"
)

//go:embed icon_black.png
var iconImageData []byte

func LoadIconImageData() image.Image {
	img, _, err := image.Decode(bytes.NewReader(iconImageData))
	if err != nil {
		log.Fatalf("failed to decode icon image: %v", err)
	}
	return img
}

func LoadIconImage() image.Image {
	iconFile, err := os.Open("assets/icon_black.png")
	if err != nil {
		log.Fatalf("failed to load icon image: %v", err)
	}
	defer iconFile.Close()

	iconImage, _, err := image.Decode(iconFile)
	if err != nil {
		log.Fatalf("failed to decode icon image: %v", err)
	}
	return iconImage
}
