package utils

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thisismemukul/snake/exceptions"
)

//go:embed icon_black.png
var iconImageData []byte

//go:embed applefood.png
var apleImageData []byte

//go:embed snakehead.png
var snakeHeadImageData []byte

func LoadIconImage() image.Image {
	img, _, err := image.Decode(bytes.NewReader(iconImageData))
	exceptions.CheckErrors(err, "Error While reading icon image")
	return img
}

func LoadAppleImage() *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(apleImageData))
	if err != nil {
		log.Fatal("Error while reading apple image: ", err)
	}
	// Convert image.Image to *ebiten.Image
	ebitenImage := ebiten.NewImageFromImage(img)
	return ebitenImage
}

func LoadSnakeHeadImage() *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(snakeHeadImageData))
	if err != nil {
		log.Fatal("Error while reading apple image: ", err)
	}
	// Convert image.Image to *ebiten.Image
	ebitenImage := ebiten.NewImageFromImage(img)
	return ebitenImage
}
