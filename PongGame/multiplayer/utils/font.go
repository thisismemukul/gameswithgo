package utils

import (
	"bytes"
	_ "embed"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed spaceranger.ttf
var fontData []byte

func LoadFontFaceData() text.Face {
	faceSrc, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatalf("Error creating text face source: %v", err)
	}
	return &text.GoTextFace{Source: faceSrc, Size: 21}
}

func LoadFontFace() text.Face {
	fontFilePath := "assets/spaceranger.ttf"
	fontFile, err := os.Open(fontFilePath)
	if err != nil {
		log.Fatalf("Error opening font file %s: %v", fontFilePath, err)
	}
	defer fontFile.Close()
	faceSrc, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		log.Fatalf("Error creating text face source: %v", err)
	}
	return &text.GoTextFace{Source: faceSrc, Size: 21}
}
