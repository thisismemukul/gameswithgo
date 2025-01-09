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
	defaultBallSpeed = 3
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
	Player1Paddle       Paddle
	Player2Paddle       Paddle
	GameBall            Ball
	Player1CurrentScore int
	Player2CurrentScore int
	FontFace            text.Face
	gameOver            bool
	winner              string
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
	game.drawTextCentered(screen, "Game Over!", -60)
	game.drawTextCentered(screen, fmt.Sprintf("Winner: %s", game.winner), -20)
	var winnerScore int
	if game.winner == "Player 1" {
		winnerScore = game.Player1CurrentScore
	} else if game.winner == "Player 2" {
		winnerScore = game.Player2CurrentScore
	}
	game.drawTextCentered(screen, fmt.Sprintf("Score: %d", winnerScore), 20)
	game.drawTextCentered(screen, "Press Enter to Restart", 60)
}

func (game *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(
		screen,
		float32(game.GameBall.Rect.PosX+game.GameBall.Rect.Width/2),
		float32(game.GameBall.Rect.PosY+game.GameBall.Rect.Height/2),
		float32(game.GameBall.Rect.Width/2),
		color.RGBA{255, 50, 50, 255},
		true,
	)

	vector.DrawFilledRect(
		screen,
		float32(game.Player1Paddle.Rect.PosX), float32(game.Player1Paddle.Rect.PosY),
		float32(game.Player1Paddle.Rect.Width), float32(game.Player1Paddle.Rect.Height),
		color.RGBA{63, 195, 128, 255}, true,
	)
	vector.DrawFilledRect(
		screen,
		float32(game.Player2Paddle.Rect.PosX), float32(game.Player2Paddle.Rect.PosY),
		float32(game.Player2Paddle.Rect.Width), float32(game.Player2Paddle.Rect.Height),
		color.RGBA{63, 128, 195, 255}, true,
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
		scoreText := fmt.Sprintf("Player1 CurrentScore: %d", game.Player1CurrentScore)
		scoreOptions := &text.DrawOptions{}
		scoreOptions.GeoM.Translate(10, 20)
		text.Draw(screen, scoreText, game.FontFace, scoreOptions)

		highScoreText := fmt.Sprintf("Player2 CurrentScore: %d", game.Player2CurrentScore)
		highScoreOptions := &text.DrawOptions{}
		highScoreOptions.GeoM.Translate(10, 40)
		text.Draw(screen, highScoreText, game.FontFace, highScoreOptions)
	}
}

func (game *Game) Update() error {
	if game.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			game.ResetGame()
		}
		return nil
	}

	game.Player1Paddle.HandleInput(ebiten.KeyW, ebiten.KeyS)
	game.Player2Paddle.HandleInput(ebiten.KeyArrowUp, ebiten.KeyArrowDown)
	game.GameBall.Move()
	game.CheckCollisions()
	return nil
}

func (p *Paddle) HandleInput(upKey, downKey ebiten.Key) {
	if ebiten.IsKeyPressed(upKey) {
		p.Rect.PosY -= paddleMoveSpeed
		if p.Rect.PosY < 0 {
			p.Rect.PosY = 0
		}
	}
	if ebiten.IsKeyPressed(downKey) {
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
	game.GameBall.SpeedX = defaultBallSpeed
	game.GameBall.SpeedY = defaultBallSpeed
	game.winner = ""
	game.Player1CurrentScore = 0
	game.Player2CurrentScore = 0
}

func (game *Game) CheckCollisions() {
	if game.GameBall.Rect.PosX >= windowWidth {
		game.gameOver = true
		game.winner = "Player 1"
	} else if game.GameBall.Rect.PosX <= 0 {
		game.gameOver = true
		game.winner = "Player 2"
	}

	if game.GameBall.Rect.PosY <= 0 || game.GameBall.Rect.PosY+game.GameBall.Rect.Height >= windowHeight {
		game.GameBall.SpeedY = -game.GameBall.SpeedY
	}

	if game.isBallCollidingWithPaddle(game.Player1Paddle.Rect) {
		game.GameBall.SpeedX = -game.GameBall.SpeedX
		game.Player1CurrentScore++
	}
	if game.isBallCollidingWithPaddle(game.Player2Paddle.Rect) {
		game.GameBall.SpeedX = -game.GameBall.SpeedX
		game.Player2CurrentScore++
	}
}

func (game *Game) isBallCollidingWithPaddle(paddle Rectangle) bool {
	ball := game.GameBall.Rect
	return ball.PosX+ball.Width >= paddle.PosX &&
		ball.PosX <= paddle.PosX+paddle.Width &&
		ball.PosY+ball.Height >= paddle.PosY &&
		ball.PosY <= paddle.PosY+paddle.Height
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

func main() {
	player1Paddle := Paddle{
		Rect: Rectangle{
			PosX:   20,
			PosY:   (windowHeight - paddleHeight) / 2,
			Width:  paddleWidth,
			Height: paddleHeight,
		},
	}

	player2Paddle := Paddle{
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
		Player1Paddle: player1Paddle,
		Player2Paddle: player2Paddle,
		GameBall:      gameBall,
		FontFace:      loadFontFace(),
	}

	ebiten.SetWindowTitle("2-Player Pong")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowIcon([]image.Image{loadIconImage()})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
