package snakegame

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (g *Game) HandleInput(upKey, downKey, leftKey, rightKey ebiten.Key) {
	if ebiten.IsKeyPressed(upKey) && g.dir.Y == 0 {
		g.dir = Position{X: 0, Y: -1}
	} else if ebiten.IsKeyPressed(downKey) && g.dir.Y == 0 {
		g.dir = Position{X: 0, Y: 1}
	} else if ebiten.IsKeyPressed(leftKey) && g.dir.X == 0 {
		g.dir = Position{X: -1, Y: 0}
	} else if ebiten.IsKeyPressed(rightKey) && g.dir.X == 0 {
		g.dir = Position{X: 1, Y: 0}
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.Restart()
		}
		return nil
	}

	g.tick++
	if g.tick%consts.SPEED != 0 {
		return nil
	}
	g.HandleInput(ebiten.KeyUp, ebiten.KeyDown, ebiten.KeyLeft, ebiten.KeyRight)

	// Move the snake
	head := g.snake[0]
	newHead := Position{X: head.X + g.dir.X, Y: head.Y + g.dir.Y}
	g.snake = append([]Position{newHead}, g.snake[:len(g.snake)-1]...)

	// Check for collisions with food
	if newHead == g.food {
		g.snake = append(g.snake, Position{})
		g.score++
		g.spawnFood()
	}

	// Check for collisions with walls or self
	if newHead.X < 0 || newHead.X >= consts.SCREEN_WIDTH/consts.GRID_SIZE || newHead.Y < 0 || newHead.Y >= consts.SCREEN_HEIGHT/consts.GRID_SIZE || g.collidesWithSelf(newHead) {
		g.gameOver = true
	}

	return nil
}

// Check if the snake collides with itself
func (g *Game) collidesWithSelf(head Position) bool {
	for _, part := range g.snake[1:] {
		if head == part {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})

	// Draw game over message
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over! Press R to Restart", consts.SCREEN_WIDTH/2-60, consts.SCREEN_HEIGHT/2)
	}

	// Draw the snake
	for _, pos := range g.snake {
		x := pos.X * consts.GRID_SIZE
		y := pos.Y * consts.GRID_SIZE

		segment := ebiten.NewImage(consts.GRID_SIZE, consts.GRID_SIZE)
		segment.Fill(color.RGBA{0x00, 0x80, 0x00, 0xff})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(segment, op)
	}

	// Draw the food
	fx := g.food.X * consts.GRID_SIZE
	fy := g.food.Y * consts.GRID_SIZE
	food := ebiten.NewImage(consts.GRID_SIZE, consts.GRID_SIZE)
	food.Fill(color.RGBA{0x80, 0x00, 0x00, 0xff})
	foodOp := &ebiten.DrawImageOptions{}
	foodOp.GeoM.Translate(float64(fx), float64(fy))
	screen.DrawImage(food, foodOp)

	// Draw score
	g.drawScores(screen)
}

func (g *Game) drawScores(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Score: "+strconv.Itoa(g.score))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func NewGame() *Game {
	g := &Game{
		randGen: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	g.Restart()
	return g
}
