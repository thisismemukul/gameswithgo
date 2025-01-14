package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game struct implementing the ebiten.Game interface
type Game struct{}

// Update method, part of the ebiten.Game interface
func (g *Game) Update() error {
	// Your game logic goes here
	return nil
}

// Draw method, part of the ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with a black color
	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})
	// Display some text
	ebitenutil.DebugPrint(screen, "Hello, Ebiten in the Browser!")
}

// Layout method, part of the ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return the window size
	return 640, 480
}

func main() {
	// Create a new Game instance
	game := &Game{}

	// Set the window size and title
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Game")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
