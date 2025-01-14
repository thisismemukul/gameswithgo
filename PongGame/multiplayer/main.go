package main

import (
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thisismemukul/pong/gameaudio"
	"github.com/thisismemukul/pong/models"
	"github.com/thisismemukul/pong/utils"
)

func main() {
	gameaudio.InitAudio()

	game := models.NewGame()

	ebiten.SetWindowTitle("Ultimate Pong!")
	ebiten.SetWindowSize(models.Config.WindowWidth, models.Config.WindowHeight)
	ebiten.SetWindowIcon([]image.Image{utils.LoadIconImage()})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
