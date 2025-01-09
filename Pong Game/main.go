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

var (
	windowWidth      = 640
	windowHeight     = 480
	paddleHeight     = 100
	paddleWidth      = 18
	ballDimension    = 25
	defaultBallSpeed = 2
	paddleMoveSpeed  = 6
)

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
	PlayerPaddle Paddle
	GameBall     Ball
	CurrentScore int
	HighScore    int
	FontFace     text.Face
	Fullscreen   bool
	gameOver     bool
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func (game *Game) drawTextCentered(screen *ebiten.Image, str string, yOffset float64) {
	x := float64(windowWidth)/2 - float64(windowWidth)/2
	y := float64(windowHeight)/2 + yOffset

	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)
	text.Draw(screen, str, game.FontFace, opts)
}

func (game *Game) DrawGameOverScreen(screen *ebiten.Image) {
	game.drawTextCentered(screen, "Game Over!", -40)
	game.drawTextCentered(screen, fmt.Sprintf("Your Score: %d", game.CurrentScore), 0)
	game.drawTextCentered(screen, "Press enter key to continue", 40)
}

func (game *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(
		screen,
		float32(game.GameBall.Rect.PosX+game.GameBall.Rect.Width/2),  // Center X position
		float32(game.GameBall.Rect.PosY+game.GameBall.Rect.Height/2), // Center Y position
		float32(game.GameBall.Rect.Width/2),                          // Radius (half of width)
		color.RGBA{255, 50, 50, 255},                                 // Fill color
		true,
	)

	vector.DrawFilledRect(
		screen,
		float32(game.PlayerPaddle.Rect.PosX), float32(game.PlayerPaddle.Rect.PosY),
		float32(game.PlayerPaddle.Rect.Width), float32(game.PlayerPaddle.Rect.Height),
		color.RGBA{63, 195, 128, 255}, true,
	)

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

	if game.gameOver {
		game.DrawGameOverScreen(screen)
	} else {
		scoreText := fmt.Sprintf("Score: %d", game.CurrentScore)
		scoreOptions := &text.DrawOptions{}
		scoreOptions.GeoM.Translate(10, 20)
		text.Draw(screen, scoreText, game.FontFace, scoreOptions)

		highScoreText := fmt.Sprintf("High Score: %d", game.HighScore)
		highScoreOptions := &text.DrawOptions{}
		highScoreOptions.GeoM.Translate(10, 40)
		text.Draw(screen, highScoreText, game.FontFace, highScoreOptions)
	}
}

func (game *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	}

	if game.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			game.ResetGame()
		}
		return nil
	}

	game.PlayerPaddle.HandleInput()
	game.GameBall.Move()
	game.CheckCollisions()
	return nil
}

func (p *Paddle) HandleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Rect.PosY -= paddleMoveSpeed
		if p.Rect.PosY < 0 {
			p.Rect.PosY = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Rect.PosY += paddleMoveSpeed
		if p.Rect.PosY+p.Rect.Height > windowHeight {
			p.Rect.PosY = windowHeight - p.Rect.Height
		}
	}
}

func (b *Ball) Move() {
	b.Rect.PosX += b.SpeedX
	b.Rect.PosY += b.SpeedY
}

func (game *Game) ResetGame() {
	game.gameOver = false
	game.GameBall.Rect.PosX = windowWidth / 2
	game.GameBall.Rect.PosY = windowHeight / 2
	game.CurrentScore = 0
}

func (game *Game) CheckCollisions() {
	if game.GameBall.Rect.PosX >= windowWidth {
		game.gameOver = true
	} else if game.GameBall.Rect.PosX <= 0 {
		game.GameBall.SpeedX = defaultBallSpeed
	} else if game.GameBall.Rect.PosY <= 0 {
		game.GameBall.SpeedY = defaultBallSpeed
	} else if game.GameBall.Rect.PosY >= windowHeight {
		game.GameBall.SpeedY = -defaultBallSpeed
	}

	// Ball collision with paddle
	if game.isBallCollidingWithPaddle() {
		game.GameBall.SpeedX = -game.GameBall.SpeedX
		game.CurrentScore++
		if game.CurrentScore > game.HighScore {
			game.HighScore = game.CurrentScore
		}
	}
}

func (game *Game) isBallCollidingWithPaddle() bool {
	ball := game.GameBall.Rect
	paddle := game.PlayerPaddle.Rect
	return ball.PosX >= paddle.PosX && ball.PosY >= paddle.PosY && ball.PosY <= paddle.PosY+paddle.Height
}

func loadIconImage() image.Image {
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

func main() {
	// Initialize game components
	playerPaddle := Paddle{
		Rect: Rectangle{
			PosX:   windowWidth - 40,
			PosY:   (windowHeight - paddleHeight) / 2,
			Width:  paddleWidth,
			Height: paddleHeight,
		},
	}

	gameBall := Ball{
		Rect: Rectangle{
			PosX:   windowWidth / 2,
			PosY:   windowHeight / 2,
			Width:  ballDimension,
			Height: ballDimension,
		},
		SpeedX: defaultBallSpeed,
		SpeedY: defaultBallSpeed,
	}

	game := &Game{
		PlayerPaddle: playerPaddle,
		GameBall:     gameBall,
	}

	ebiten.SetWindowTitle("Pong Game")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowIcon([]image.Image{loadIconImage()})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
