package main

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thisismemukul/snake/exceptions"
	snakegame "github.com/thisismemukul/snake/game"
	"github.com/thisismemukul/snake/utils"
)

func main() {

	game := snakegame.NewGame()
	ebiten.SetWindowTitle("Snake Game!")
	ebiten.SetWindowIcon([]image.Image{utils.LoadIconImage()})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		exceptions.CheckErrors(err, "Error while running game.")
	}

}
