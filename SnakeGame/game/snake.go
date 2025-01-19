package snakegame

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thisismemukul/snake/consts"
)

type Game struct {
	snake    []Position
	dir      Position
	food     Position
	score    int
	gameOver bool
	tick     int
	randGen  *rand.Rand
}

func (g *Game) Restart() {
	g.snake = []Position{{X: 5, Y: 5}, {X: 4, Y: 5}, {X: 3, Y: 5}, {X: 2, Y: 5}, {X: 1, Y: 5}}
	g.dir = Position{X: 1, Y: 0}
	g.score = 0
	g.gameOver = false
	g.tick = 0
	g.spawnFood()
}

func (g *Game) spawnFood() {
	g.food = Position{
		X: g.randGen.Intn(consts.SCREEN_WIDTH / consts.GRID_SIZE),
		Y: g.randGen.Intn(consts.SCREEN_HEIGHT / consts.GRID_SIZE),
	}
}
func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func NewGame() *Game {
	fmt.Print("nhf")
	return &Game{}
}
