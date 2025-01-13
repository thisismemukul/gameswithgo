package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Embed audio files
//
//go:embed assets/background.mp3
var backgroundMusicData []byte

//go:embed assets/hitpaddleball.mp3
var hitPaddleBallSoundData []byte

//go:embed assets/hitwallball.mp3
var hitWallBallSoundData []byte

//go:embed assets/outofboundary.mp3
var outOfBoundarySoundData []byte

//go:embed assets/gameover.mp3
var gameOverSoundData []byte

const (
	SampleRate          = 44100
	WordsPerSec         = 2.71828
	SpeedIncreaseFactor = 1.1
	MaxBallSpeed        = 20
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
	GameOverWords     []string
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
	GameOverWords: []string{
		"G", "G\tA", "G\tA\tM", "G\tA\tM\tE\t\tO", "o\te\t", "\t\t\tr",
		"G A M E   O V", "GAME OVER", "GAME OVER", "GAME OVER", "GAME OVER",
	},
}

var (
	audioContext             = audio.NewContext(SampleRate)
	backgroundPlayer         *audio.Player
	hitPaddleBallSoundPlayer *audio.Player
	hitWallBallSoundPlayer   *audio.Player
	outOfBoundarySoundPlayer *audio.Player
	gameOverSoundDataPlayer  *audio.Player
	// infiniteLoop             *audio.InfiniteLoop
)

// Game levels
type DifficultyLevel struct {
	BallSpeed   int
	PaddleSpeed int
}

var Levels = map[string]DifficultyLevel{
	"Easy":   {BallSpeed: 3, PaddleSpeed: 6},
	"Medium": {BallSpeed: 5, PaddleSpeed: 8},
	"Hard":   {BallSpeed: 6, PaddleSpeed: 9},
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
	WordIndex           float64
	SelectedLevel       string
}

func init() {
	// Load background music
	stream, err := mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(backgroundMusicData))
	if err != nil {
		log.Fatalf("Failed to decode background music: %v", err)
	}
	backgroundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create background music player: %v", err)
	}
	// backgroundPlayer.SetLoop(true)
	// if infiniteLoop != nil {
	// 	backgroundPlayer.Play()
	// }

	// Load hit paddle ball sound
	stream, err = mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(hitPaddleBallSoundData))
	if err != nil {
		log.Fatalf("Failed to decode hit paddle ball sound: %v", err)
	}
	hitPaddleBallSoundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create hit paddle ball sound player: %v", err)
	}

	// Load kit wall ball sound
	stream, err = mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(hitWallBallSoundData))
	if err != nil {
		log.Fatalf("Failed to decode hit wall sound: %v", err)
	}
	hitWallBallSoundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create hit wall sound player: %v", err)
	}

	// Load ball out of boundary sound
	stream, err = mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(outOfBoundarySoundData))
	if err != nil {
		log.Fatalf("Failed to decode ball out of boundary sound: %v", err)
	}
	outOfBoundarySoundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create ball out of boundary sound player: %v", err)
	}

	// Load game over sound
	stream, err = mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(gameOverSoundData))
	if err != nil {
		log.Fatalf("Failed to decode hit wall sound: %v", err)
	}
	gameOverSoundDataPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create hit wall sound player: %v", err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.WindowWidth, config.WindowHeight
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0}) // background color
	if game.GameOver {
		if game.GameOver {
			if backgroundPlayer.IsPlaying() {
				backgroundPlayer.Close()
			}
			if gameOverSoundDataPlayer != nil && !gameOverSoundDataPlayer.IsPlaying() {
				gameOverSoundDataPlayer.Rewind()
				gameOverSoundDataPlayer.Play()
			}
			game.drawGameOverScreen(screen)
			return
		}

		game.drawGameOverScreen(screen)
		return
	}
	if game.SelectedLevel == "" {
		game.drawSelectGameLevelScreen(screen)
		return
	}
	game.drawMidLine(screen)
	game.drawPaddles(screen)
	game.drawBall(screen)
	game.drawScores(screen)
}

func (game *Game) Update() error {

	if !backgroundPlayer.IsPlaying() {
		backgroundPlayer.Play()
	}

	if game.SelectedLevel == "" {
		if ebiten.IsKeyPressed(ebiten.Key1) {
			game.SelectedLevel = "Easy"
			game.applyLevelSettings()
		} else if ebiten.IsKeyPressed(ebiten.Key2) {
			game.SelectedLevel = "Medium"
			game.applyLevelSettings()
		} else if ebiten.IsKeyPressed(ebiten.Key3) {
			game.SelectedLevel = "Hard"
			game.applyLevelSettings()
		}
		return nil
	}
	if game.GameOver {
		newIndex := (game.WordIndex + WordsPerSec/60.0)
		game.WordIndex = math.Mod(newIndex, float64(len(config.GameOverWords)))
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			game.resetGame()
		}
		return nil
	}

	game.Player1Paddle.HandleInput(ebiten.KeyW, ebiten.KeyS, Levels[game.SelectedLevel].PaddleSpeed, game.GameBall.SpeedX)
	game.Player2Paddle.HandleInput(ebiten.KeyArrowUp, ebiten.KeyArrowDown, Levels[game.SelectedLevel].PaddleSpeed, game.GameBall.SpeedX)
	game.GameBall.Move()
	game.checkCollisions()
	return nil
}

func (game *Game) applyLevelSettings() {
	level := Levels[game.SelectedLevel]
	game.GameBall.SpeedX = level.BallSpeed
	game.GameBall.SpeedY = level.BallSpeed
}

// Drawing Helpers
func (g *Game) drawMidLine(screen *ebiten.Image) {
	midX := config.WindowWidth / 2
	for y := 0; y < config.WindowHeight; y += 20 {
		if y%40 == 0 {
			drawRect(screen, Rectangle{PosX: midX - 1, PosY: y, Width: 2, Height: 20}, color.RGBA{200, 200, 200, 255})
		}
	}
}

func (game *Game) drawSelectGameLevelScreen(screen *ebiten.Image) {
	game.drawTextCentered(screen, "Need 5 Scores to win", -100)
	game.drawTextCentered(screen, "Select level", -60)
	game.drawTextCentered(screen, "Easy:\tPress 1", -20)
	game.drawTextCentered(screen, "Medium:\tPress 2", 20)
	game.drawTextCentered(screen, "Hard:\tPress 3", 60)
}

func (game *Game) drawGameOverScreen(screen *ebiten.Image) {
	// word := config.GameOverWords[int(game.WordIndex)]
	word := config.GameOverWords[int(game.WordIndex)%len(config.GameOverWords)]
	game.drawTextCentered(screen, word, -60)

	game.drawTextCentered(screen, fmt.Sprintf("Winner: %s", game.Winner), -20)

	scoreDifference := int(math.Abs(float64(game.Player1CurrentScore - game.Player2CurrentScore)))
	game.drawTextCentered(screen, fmt.Sprintf("Won by Score: %d", scoreDifference), 20)

	game.drawTextCentered(screen, "Press Enter to Restart", 60)
}

func (game *Game) drawTextCentered(screen *ebiten.Image, str string, yOffset float64) {
	y := float64(config.WindowHeight)/2 + yOffset
	bounds := screen.Bounds()
	x := float64(bounds.Dx() / 2)
	opts := &text.DrawOptions{}
	opts.GeoM.Translate((x)-float64(len(str)*8/2), (y))
	text.Draw(screen, str, game.FontFace, opts)
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
	scoreText := fmt.Sprintf("Player1 Score: %d", g.Player1CurrentScore)
	scoreOptions := &text.DrawOptions{}
	scoreOptions.GeoM.Translate(10, 20)
	text.Draw(screen, scoreText, g.FontFace, scoreOptions)

	highScoreText := fmt.Sprintf("Player2 Score: %d", g.Player2CurrentScore)
	highScoreOptions := &text.DrawOptions{}
	highScoreOptions.GeoM.Translate(float64((config.WindowWidth/2)+10), 20)
	text.Draw(screen, highScoreText, g.FontFace, highScoreOptions)
}

// Paddle and Ball Movement
func (p *Paddle) HandleInput(upKey, downKey ebiten.Key, paddleSpeed int, ballSpeed int) {

	// paddleSpeed := config.PaddleMoveSpeed + int(math.Abs(float64(ballSpeed))/2)
	newPaddleSpeed := config.PaddleMoveSpeed + int(math.Min(float64(ballSpeed)/3, 10)) // limit paddle speed increase

	paddleSpeed = max(paddleSpeed, newPaddleSpeed)

	if ebiten.IsKeyPressed(upKey) {
		p.Rect.PosY = max(0, p.Rect.PosY-paddleSpeed)
	}
	if ebiten.IsKeyPressed(downKey) {
		p.Rect.PosY = min(config.WindowHeight-p.Rect.Height, p.Rect.PosY+paddleSpeed)
	}
}

func (b *Ball) Move() {
	b.Rect.PosX += b.SpeedX
	b.Rect.PosY += b.SpeedY
}

// Collision Detection purana
func (g *Game) checkCollisions() {
	// Ball Out of Bounds
	if g.GameBall.Rect.PosX < 0 {
		g.Player2CurrentScore++
		if outOfBoundarySoundPlayer != nil {
			outOfBoundarySoundPlayer.Rewind()
			outOfBoundarySoundPlayer.Play()
		}
		g.resetBallPosition()
	} else if g.GameBall.Rect.PosX+g.GameBall.Rect.Width > config.WindowWidth {
		g.Player1CurrentScore++
		if outOfBoundarySoundPlayer != nil {
			outOfBoundarySoundPlayer.Rewind()
			outOfBoundarySoundPlayer.Play()
		}
		g.resetBallPosition()
	}

	if g.Player1CurrentScore >= 5 {
		g.declareWinner("Player 1")
	} else if g.Player2CurrentScore >= 5 {
		g.declareWinner("Player 2")
	}

	// Ball Hits Walls
	if g.GameBall.Rect.PosY <= 0 || g.GameBall.Rect.PosY+g.GameBall.Rect.Height >= config.WindowHeight {
		g.GameBall.SpeedY = -g.GameBall.SpeedY
		if hitWallBallSoundPlayer != nil {
			hitWallBallSoundPlayer.Rewind()
			hitWallBallSoundPlayer.Play()
		}
	}

	// Ball hits paddles
	if g.isBallCollidingWithPaddle(g.Player1Paddle.Rect) || g.isBallCollidingWithPaddle(g.Player2Paddle.Rect) {
		if hitPaddleBallSoundPlayer != nil {
			hitPaddleBallSoundPlayer.Rewind()
			hitPaddleBallSoundPlayer.Play()
		}
		g.GameBall.SpeedX = -g.GameBall.SpeedX
		g.increaseSpeed()
	}

	if g.GameBall.SpeedX > MaxBallSpeed {
		g.GameBall.SpeedX = MaxBallSpeed
	}
	if g.GameBall.SpeedX < -MaxBallSpeed {
		g.GameBall.SpeedX = -MaxBallSpeed
	}
	if g.GameBall.SpeedY > MaxBallSpeed {
		g.GameBall.SpeedY = MaxBallSpeed
	}
	if g.GameBall.SpeedY < -MaxBallSpeed {
		g.GameBall.SpeedY = -MaxBallSpeed
	}
}

func (g *Game) increaseSpeed() {
	g.GameBall.SpeedX = int(float64(g.GameBall.SpeedX) * SpeedIncreaseFactor)
	g.GameBall.SpeedY = int(float64(g.GameBall.SpeedY) * SpeedIncreaseFactor)
	const minSpeed = 5
	if g.GameBall.SpeedX < minSpeed && g.GameBall.SpeedX > -minSpeed {
		if g.GameBall.SpeedX > 0 {
			g.GameBall.SpeedX = minSpeed
		} else {
			g.GameBall.SpeedX = -minSpeed
		}
	}

	if g.GameBall.SpeedY < minSpeed && g.GameBall.SpeedY > -minSpeed {
		if g.GameBall.SpeedY > 0 {
			g.GameBall.SpeedY = minSpeed
		} else {
			g.GameBall.SpeedY = -minSpeed
		}
	}
}

func (g *Game) isBallCollidingWithPaddle(paddle Rectangle) bool {
	ball := g.GameBall.Rect
	return ball.PosX < paddle.PosX+paddle.Width &&
		ball.PosX+ball.Width > paddle.PosX &&
		ball.PosY < paddle.PosY+paddle.Height &&
		ball.PosY+ball.Height > paddle.PosY
}

func (g *Game) resetBallPosition() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	directionX := rng.Intn(2)*2 - 1
	directionY := rng.Intn(2)*2 - 1
	g.GameOver = false
	g.GameBall = Ball{
		Rect: Rectangle{
			PosX:   config.WindowWidth / 2,
			PosY:   config.WindowHeight / 2,
			Width:  config.BallSize,
			Height: config.BallSize,
		},
		SpeedX: directionX * Levels[g.SelectedLevel].BallSpeed,
		SpeedY: directionY * Levels[g.SelectedLevel].BallSpeed,
	}
}

func (g *Game) declareWinner(winner string) {
	g.Winner = winner
	g.GameOver = true
}

// Game Reset
func (g *Game) resetGame() {
	g.resetBallPosition()
	g.Player1CurrentScore = 0
	g.Player2CurrentScore = 0
	g.Winner = ""
}

func loadFontFace() text.Face {
	fontFilePath := "../assets/spaceranger.ttf"
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
		SelectedLevel: "",
	}

	ebiten.SetWindowTitle("Ultimate Pong!")
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowIcon([]image.Image{loadIconImage()})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
