package models

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thisismemukul/pong/consts"
	"github.com/thisismemukul/pong/gameaudio"
	"github.com/thisismemukul/pong/utils"
)

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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Config.WindowWidth, Config.WindowHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0}) // background color
	if g.GameOver {
		g.drawGameOverScreen(screen)
		return
	}
	if g.SelectedLevel == "" {
		g.drawSelectGameLevelScreen(screen)
		return
	}
	g.drawMidLine(screen)
	g.drawPaddles(screen)
	g.drawBall(screen)
	g.drawScores(screen)
}
func isWithinRect(px, py, x, y, width, height int) bool {
	return px >= x && px <= x+width && py >= y && py <= y+height
}

func (g *Game) Update() error {

	if !gameaudio.BackgroundPlayer.IsPlaying() {
		gameaudio.BackgroundPlayer.Play()
	}

	if g.SelectedLevel == "" {
		// Handle keyboard input
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.SelectedLevel = "Easy"
			g.applyLevelSettings()
		} else if ebiten.IsKeyPressed(ebiten.Key2) {
			g.SelectedLevel = "Medium"
			g.applyLevelSettings()
		} else if ebiten.IsKeyPressed(ebiten.Key3) {
			g.SelectedLevel = "Hard"
			g.applyLevelSettings()
		}

		// Handle mouse clicks
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if isWithinRect(x, y, Config.WindowWidth/2-50, Config.WindowHeight/2-30, 100, 30) {
				g.SelectedLevel = "Easy"
				g.applyLevelSettings()
			} else if isWithinRect(x, y, Config.WindowWidth/2-50, Config.WindowHeight/2+10, 100, 30) {
				g.SelectedLevel = "Medium"
				g.applyLevelSettings()
			} else if isWithinRect(x, y, Config.WindowWidth/2-50, Config.WindowHeight/2+50, 100, 30) {
				g.SelectedLevel = "Hard"
				g.applyLevelSettings()
			}
		}

		// Handle touch input
		touchIDs := make([]ebiten.TouchID, 0)
		touchIDs = ebiten.AppendTouchIDs(touchIDs)
		for _, touchID := range touchIDs {
			tx, ty := ebiten.TouchPosition(touchID)
			if isWithinRect(tx, ty, Config.WindowWidth/2-50, Config.WindowHeight/2-30, 100, 30) {
				g.SelectedLevel = "Easy"
				g.applyLevelSettings()
			} else if isWithinRect(tx, ty, Config.WindowWidth/2-50, Config.WindowHeight/2+10, 100, 30) {
				g.SelectedLevel = "Medium"
				g.applyLevelSettings()
			} else if isWithinRect(tx, ty, Config.WindowWidth/2-50, Config.WindowHeight/2+50, 100, 30) {
				g.SelectedLevel = "Hard"
				g.applyLevelSettings()
			}
		}
		return nil
	}
	if g.GameOver {
		if gameaudio.BackgroundPlayer.IsPlaying() {
			gameaudio.BackgroundPlayer.Pause()
		}
		if gameaudio.GameOverSoundDataPlayer != nil && !gameaudio.GameOverSoundDataPlayer.IsPlaying() {
			gameaudio.GameOverSoundDataPlayer.Rewind()
			gameaudio.GameOverSoundDataPlayer.Play()
		}
		newIndex := (g.WordIndex + consts.WordsPerSec/60.0)
		g.WordIndex = math.Mod(newIndex, float64(len(Config.GameOverWords)))
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.resetGame()
		}
		return nil
	}
	touchIDs := ebiten.AppendTouchIDs(nil)
	g.Player1Paddle.HandleInput(ebiten.KeyW, ebiten.KeyS, Levels[g.SelectedLevel].PaddleSpeed, touchIDs, g.GameBall.SpeedX)
	g.Player2Paddle.HandleInput(ebiten.KeyArrowUp, ebiten.KeyArrowDown, Levels[g.SelectedLevel].PaddleSpeed, touchIDs, g.GameBall.SpeedX)
	g.GameBall.Move()
	g.checkCollisions()
	return nil
}

func (game *Game) applyLevelSettings() {
	level := Levels[game.SelectedLevel]
	game.GameBall.SpeedX = level.BallSpeed
	game.GameBall.SpeedY = level.BallSpeed
}

// Drawing Helpers
func (g *Game) drawMidLine(screen *ebiten.Image) {
	midX := Config.WindowWidth / 2
	for y := 0; y < Config.WindowHeight; y += 20 {
		if y%40 == 0 {
			drawRect(screen, Rectangle{PosX: midX - 1, PosY: y, Width: 2, Height: 20}, color.RGBA{200, 200, 200, 255})
		}
	}
}

func drawButton(screen *ebiten.Image, x, y, width, height int, label string) {
	drawRect(screen, Rectangle{PosX: x, PosY: y, Width: width + width/3, Height: height}, color.RGBA{100, 100, 100, 255})
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(x+width/4), float64(y+height/3))
	text.Draw(screen, label, utils.LoadFontFaceData(), opts)
}

func (game *Game) drawSelectGameLevelScreen(screen *ebiten.Image) {
	game.drawTextCentered(screen, "Need 5 Scores to win", -100)
	game.drawTextCentered(screen, "Select level:", -60)
	drawButton(screen, Config.WindowWidth/2-50, Config.WindowHeight/2-30, 100, 30, "Easy")
	drawButton(screen, Config.WindowWidth/2-50, Config.WindowHeight/2+10, 100, 30, "Medium")
	drawButton(screen, Config.WindowWidth/2-50, Config.WindowHeight/2+50, 100, 30, "Hard")
}

func (game *Game) drawGameOverScreen(screen *ebiten.Image) {
	// word := Config.GameOverWords[int(game.WordIndex)]
	word := Config.GameOverWords[int(game.WordIndex)%len(Config.GameOverWords)]
	game.drawTextCentered(screen, word, -60)

	game.drawTextCentered(screen, fmt.Sprintf("Winner: %s", game.Winner), -20)

	scoreDifference := int(math.Abs(float64(game.Player1CurrentScore - game.Player2CurrentScore)))
	game.drawTextCentered(screen, fmt.Sprintf("Won by Score: %d", scoreDifference), 20)

	game.drawTextCentered(screen, "Press Enter to Restart", 60)
}

func (game *Game) drawTextCentered(screen *ebiten.Image, str string, yOffset float64) {
	y := float64(Config.WindowHeight)/2 + yOffset
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
	highScoreOptions.GeoM.Translate(float64((Config.WindowWidth/2)+10), 20)
	text.Draw(screen, highScoreText, g.FontFace, highScoreOptions)
}

// Collision Detection purana
func (g *Game) checkCollisions() {
	// Ball Out of Bounds
	if g.GameBall.Rect.PosX < 0 {
		g.Player2CurrentScore++
		if gameaudio.OutOfBoundarySoundPlayer != nil {
			gameaudio.OutOfBoundarySoundPlayer.Rewind()
			gameaudio.OutOfBoundarySoundPlayer.Play()
		}
		g.resetBallPosition()
	} else if g.GameBall.Rect.PosX+g.GameBall.Rect.Width > Config.WindowWidth {
		g.Player1CurrentScore++
		if gameaudio.OutOfBoundarySoundPlayer != nil {
			gameaudio.OutOfBoundarySoundPlayer.Rewind()
			gameaudio.OutOfBoundarySoundPlayer.Play()
		}
		g.resetBallPosition()
	}

	if g.Player1CurrentScore >= 5 {
		g.declareWinner("Player 1")
	} else if g.Player2CurrentScore >= 5 {
		g.declareWinner("Player 2")
	}

	// Ball Hits Walls
	if g.GameBall.Rect.PosY <= 0 || g.GameBall.Rect.PosY+g.GameBall.Rect.Height >= Config.WindowHeight {
		g.GameBall.SpeedY = -g.GameBall.SpeedY
		if gameaudio.HitWallBallSoundPlayer != nil {
			gameaudio.HitWallBallSoundPlayer.Rewind()
			gameaudio.HitWallBallSoundPlayer.Play()
		}
	}

	// Ball hits paddles
	if g.isBallCollidingWithPaddle(g.Player1Paddle.Rect) || g.isBallCollidingWithPaddle(g.Player2Paddle.Rect) {
		if gameaudio.HitPaddleBallSoundPlayer != nil {
			gameaudio.HitPaddleBallSoundPlayer.Rewind()
			gameaudio.HitPaddleBallSoundPlayer.Play()
		}
		g.GameBall.SpeedX = -g.GameBall.SpeedX
		g.increaseSpeed()
	}

	if g.GameBall.SpeedX > consts.MaxBallSpeed {
		g.GameBall.SpeedX = consts.MaxBallSpeed
	}
	if g.GameBall.SpeedX < -consts.MaxBallSpeed {
		g.GameBall.SpeedX = -consts.MaxBallSpeed
	}
	if g.GameBall.SpeedY > consts.MaxBallSpeed {
		g.GameBall.SpeedY = consts.MaxBallSpeed
	}
	if g.GameBall.SpeedY < -consts.MaxBallSpeed {
		g.GameBall.SpeedY = -consts.MaxBallSpeed
	}
}

func (g *Game) increaseSpeed() {
	g.GameBall.SpeedX = int(float64(g.GameBall.SpeedX) * consts.SpeedIncreaseFactor)
	g.GameBall.SpeedY = int(float64(g.GameBall.SpeedY) * consts.SpeedIncreaseFactor)
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
			PosX:   Config.WindowWidth / 2,
			PosY:   Config.WindowHeight / 2,
			Width:  Config.BallSize,
			Height: Config.BallSize,
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

func NewGame() *Game {

	player1Paddle := Paddle{
		Rect: Rectangle{
			PosX:   20,
			PosY:   (Config.WindowHeight - Config.PaddleHeight) / 2,
			Width:  Config.PaddleWidth,
			Height: Config.PaddleHeight,
		},
	}
	player2Paddle := Paddle{
		Rect: Rectangle{
			PosX:   Config.WindowWidth - 40,
			PosY:   (Config.WindowHeight - Config.PaddleHeight) / 2,
			Width:  Config.PaddleWidth,
			Height: Config.PaddleHeight,
		},
	}

	gameBall := Ball{
		Rect: Rectangle{
			PosX:   Config.WindowWidth / 2,
			PosY:   Config.WindowHeight / 2,
			Width:  Config.BallSize,
			Height: Config.BallSize,
		},
		SpeedX: Config.DefaultBallSpeed,
		SpeedY: Config.DefaultBallSpeed,
	}

	return &Game{
		Player1Paddle: player1Paddle,
		Player2Paddle: player2Paddle,
		GameBall:      gameBall,
		FontFace:      utils.LoadFontFaceData(),
		SelectedLevel: "",
	}
}
