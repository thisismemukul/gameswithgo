package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

// Constants and Game Configuration
type GameConfig struct {
	WindowWidth       int
	WindowHeight      int
	PaddleHeight      int
	PaddleWidth       int
	BallSize          int
	DefaultBallSpeed  int
	PaddleMoveSpeed   int
	DifficultyScaling float64
}

var config = GameConfig{
	WindowWidth:       640,
	WindowHeight:      480,
	PaddleHeight:      100,
	PaddleWidth:       18,
	BallSize:          25,
	DefaultBallSpeed:  3,
	PaddleMoveSpeed:   6,
	DifficultyScaling: 0.1,
}

type Rectangle struct {
	PosX, PosY, Width, Height int
}

type Paddle struct {
	Rect Rectangle
}

type Ball struct {
	Rect           Rectangle
	SpeedX, SpeedY int
}

type Game struct {
	Player1Paddle       Paddle
	Player2Paddle       Paddle
	GameBall            Ball
	Player1CurrentScore int
	Player2CurrentScore int
	FontFace            text.Face
	GameOver            bool
	Winner              string
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.WindowWidth, config.WindowHeight
}

func (game *Game) Draw(screen *ebiten.Image) {
	if game.GameOver {
		game.drawGameOverScreen(screen)
		return
	}
	game.drawPaddles(screen)
	game.drawBall(screen)
	game.drawScores(screen)

	if game.FontFace == nil {
		faceSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			log.Fatal(err)
		}
		game.FontFace = &text.GoTextFace{
			Source: faceSrc,
			Size:   21,
		}
	}
}

func (game *Game) Update() error {
	if game.GameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			game.resetGame()
		}
		return nil
	}

	game.Player1Paddle.HandleInput(ebiten.KeyW, ebiten.KeyS)
	game.Player2Paddle.HandleInput(ebiten.KeyArrowUp, ebiten.KeyArrowDown)
	game.GameBall.Move()
	game.checkCollisions()
	return nil
}

// Drawing Helpers
func (game *Game) drawTextCentered(screen *ebiten.Image, str string, yOffset float64) {
	x := float64(config.WindowWidth)/2 - float64(config.WindowWidth)/2
	y := float64(config.WindowHeight)/2 + yOffset
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)
	text.Draw(screen, str, game.FontFace, opts)
}

func (game *Game) drawGameOverScreen(screen *ebiten.Image) {
	game.drawTextCentered(screen, "Game Over!", -60)
	game.drawTextCentered(screen, fmt.Sprintf("Winner: %s", game.Winner), -20)
	score := game.Player1CurrentScore
	if game.Winner == "Player 2" {
		score = game.Player2CurrentScore
	}
	game.drawTextCentered(screen, fmt.Sprintf("Score: %d", score), 20)
	game.drawTextCentered(screen, "Press Enter to Restart", 60)
}

func (g *Game) drawPaddles(screen *ebiten.Image) {
	drawRect(screen, g.Player1Paddle.Rect, color.RGBA{63, 195, 128, 255})
	drawRect(screen, g.Player2Paddle.Rect, color.RGBA{63, 128, 195, 255})
}

// Utility Functions
func drawRect(screen *ebiten.Image, rect Rectangle, clr color.Color) {
	vector.DrawFilledRect(screen, float32(rect.PosX), float32(rect.PosY), float32(rect.Width), float32(rect.Height), clr, true)
}

func (g *Game) drawBall(screen *ebiten.Image) {
	drawCircle(screen, g.GameBall.Rect, color.RGBA{255, 50, 50, 255})
}

func drawCircle(screen *ebiten.Image, rect Rectangle, clr color.Color) {
	vector.DrawFilledCircle(screen, float32(rect.PosX+rect.Width/2), float32(rect.PosY+rect.Height/2), float32(rect.Width/2), clr, true)
}

func (g *Game) drawScores(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("Player1 CurrentScore: %d", g.Player1CurrentScore)
	scoreOptions := &text.DrawOptions{}
	scoreOptions.GeoM.Translate(10, 20)
	text.Draw(screen, scoreText, g.FontFace, scoreOptions)

	highScoreText := fmt.Sprintf("Player2 CurrentScore: %d", g.Player2CurrentScore)
	highScoreOptions := &text.DrawOptions{}
	highScoreOptions.GeoM.Translate(10, 40)
	text.Draw(screen, highScoreText, g.FontFace, highScoreOptions)
}

// Paddle and Ball Movement
func (p *Paddle) HandleInput(upKey, downKey ebiten.Key) {
	if ebiten.IsKeyPressed(upKey) {
		p.Rect.PosY = max(0, p.Rect.PosY-config.PaddleMoveSpeed)
	}
	if ebiten.IsKeyPressed(downKey) {
		p.Rect.PosY = min(config.WindowHeight-p.Rect.Height, p.Rect.PosY+config.PaddleMoveSpeed)
	}
}

func (b *Ball) Move() {
	b.Rect.PosX += b.SpeedX
	b.Rect.PosY += b.SpeedY
}

// Collision Detection
func (g *Game) checkCollisions() {
	// Ball Out of Bounds
	if g.GameBall.Rect.PosX < 0 {
		g.GameOver = true
		g.Winner = "Player 2"
	} else if g.GameBall.Rect.PosX+g.GameBall.Rect.Width > config.WindowWidth {
		g.GameOver = true
		g.Winner = "Player 1"
	}

	// Ball Hits Walls
	if g.GameBall.Rect.PosY <= 0 || g.GameBall.Rect.PosY+g.GameBall.Rect.Height >= config.WindowHeight {
		g.GameBall.SpeedY = -g.GameBall.SpeedY
	}

	// Ball Hits Paddles
	if g.isBallCollidingWithPaddle(g.Player1Paddle.Rect) {
		g.GameBall.SpeedX = -g.GameBall.SpeedX
		g.Player1CurrentScore++
		// Optional: Speed scaling logic
		g.GameBall.SpeedX += int(config.DifficultyScaling * float64(g.Player1CurrentScore+g.Player2CurrentScore))
	}
	if g.isBallCollidingWithPaddle(g.Player2Paddle.Rect) {
		g.GameBall.SpeedX = -g.GameBall.SpeedX
		g.Player2CurrentScore++
		// Optional: Speed scaling logic
		g.GameBall.SpeedX += int(config.DifficultyScaling * float64(g.Player1CurrentScore+g.Player2CurrentScore))
	}
}
func (g *Game) isBallCollidingWithPaddle(paddle Rectangle) bool {
	ball := g.GameBall.Rect
	return ball.PosX < paddle.PosX+paddle.Width &&
		ball.PosX+ball.Width > paddle.PosX &&
		ball.PosY < paddle.PosY+paddle.Height &&
		ball.PosY+ball.Height > paddle.PosY
}

// Game Reset
func (g *Game) resetGame() {
	g.GameOver = false
	g.GameBall = Ball{
		Rect: Rectangle{
			PosX:   config.WindowWidth / 2,
			PosY:   config.WindowHeight / 2,
			Width:  config.BallSize,
			Height: config.BallSize,
		},
		SpeedX: config.DefaultBallSpeed,
		SpeedY: config.DefaultBallSpeed,
	}
	g.Player1CurrentScore = 0
	g.Player2CurrentScore = 0
	g.Winner = ""
}

func loadFontFace() text.Face {
	faceSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	return &text.GoTextFace{Source: faceSrc, Size: 21}
}

func loadIconImage() image.Image {
	iconFile, err := os.Open("../assets/icon_black.png")
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

// Main Function
func main() {
	player1Paddle := Paddle{
		Rect: Rectangle{
			PosX:   20,
			PosY:   (config.WindowHeight - config.PaddleHeight) / 2,
			Width:  config.PaddleWidth,
			Height: config.PaddleHeight,
		},
	}

	player2Paddle := Paddle{
		Rect: Rectangle{
			PosX:   config.WindowWidth - 40,
			PosY:   (config.WindowHeight - config.PaddleHeight) / 2,
			Width:  config.PaddleWidth,
			Height: config.PaddleHeight,
		},
	}

	gameBall := Ball{
		Rect: Rectangle{
			PosX:   config.WindowWidth / 2,
			PosY:   config.WindowHeight / 2,
			Width:  config.BallSize,
			Height: config.BallSize,
		},
		SpeedX: config.DefaultBallSpeed,
		SpeedY: config.DefaultBallSpeed,
	}

	game := &Game{
		Player1Paddle: player1Paddle,
		Player2Paddle: player2Paddle,
		GameBall:      gameBall,
		FontFace:      loadFontFace(),
	}

	ebiten.SetWindowTitle("Ultimate Pong!")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowIcon([]image.Image{loadIconImage()})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
