package utils

import (
	"image"
	"log"
	"os"
)

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
