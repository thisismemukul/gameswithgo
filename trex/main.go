package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 200
	dinoWidth    = 30.0 // Dino width
	dinoHeight   = 50.0 // Dino height
	cactusWidth  = 20.0 // Cactus width
	cactusHeight = 40.0 // Cactus height
	birdWidth    = 40.0 // Bird width
	birdHeight   = 30.0 // Bird height
)

var (
	dinoX        = 100.0 // Dino X position
	dinoY        = screenHeight - dinoHeight
	dinoVelocity = 0.0
	jumpStrength = -12.0
	gravity      = 0.6
	isJumping    = false
	score        int
	cactuses     []*Cactus
	birds        []*Bird

	// Use the new random generator
	randomGen = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Cactus struct {
	x, y float64
}

type Bird struct {
	x, y, velocity float64
}

type Game struct{}

func (g *Game) Update() error {
	// Dino jump logic
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !isJumping {
		dinoVelocity = jumpStrength
		isJumping = true
	}

	// Apply gravity
	if isJumping {
		dinoY += dinoVelocity
		dinoVelocity += gravity
	}

	// If dino hits the ground, stop jumping
	if dinoY >= screenHeight-dinoHeight {
		dinoY = screenHeight - dinoHeight
		isJumping = false
	}

	// Move the obstacles (cacti) and remove those off-screen
	for i := len(cactuses) - 1; i >= 0; i-- {
		cactuses[i].x -= 5
		// Remove cactus if it's off-screen
		if cactuses[i].x < -cactusWidth {
			cactuses = append(cactuses[:i], cactuses[i+1:]...)
			score++
		}
	}

	// Move the birds and remove those off-screen
	for i := len(birds) - 1; i >= 0; i-- {
		birds[i].x -= birds[i].velocity
		// Remove bird if it's off-screen
		if birds[i].x < -birdWidth {
			birds = append(birds[:i], birds[i+1:]...)
			score++
		}
	}

	// Spawn new cactus
	if randomGen.Float64() < 0.02 {
		cactuses = append(cactuses, &Cactus{
			x: screenWidth,
			y: screenHeight - cactusHeight,
		})
	}

	// Spawn new bird with random velocity
	if randomGen.Float64() < 0.01 {
		birds = append(birds, &Bird{
			x:        screenWidth,
			y:        randomGen.Float64() * (screenHeight - birdHeight),
			velocity: randomGen.Float64()*3.0 + 2.0, // Random velocity for bird
		})
	}

	// Check for collision with cactuses or birds
	for _, cactus := range cactuses {
		if dinoX+dinoWidth > cactus.x && dinoX < cactus.x+cactusWidth && dinoY+dinoHeight > cactus.y {
			// Game Over (collision with cactus)
			return fmt.Errorf("Game Over! Score: %d", score)
		}
	}

	for _, bird := range birds {
		if dinoX+dinoWidth > bird.x && dinoX < bird.x+birdWidth && dinoY+dinoHeight > bird.y && dinoY < bird.y+birdHeight {
			// Game Over (collision with bird)
			return fmt.Errorf("Game Over! Score: %d", score)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with white color
	screen.Fill(color.RGBA{234, 156, 255, 255})

	// Draw Dino (using a simple rectangle for now)
	vector.DrawFilledRect(screen, float32(dinoX), float32(dinoY), dinoWidth, dinoHeight, color.RGBA{63, 128, 195, 255}, true)

	// Draw Cactuses
	for _, cactus := range cactuses {
		vector.DrawFilledRect(screen, float32(cactus.x), float32(cactus.y), cactusWidth, cactusHeight, color.RGBA{63, 128, 195, 255}, true)
	}

	// Draw Birds
	for _, bird := range birds {
		vector.DrawFilledRect(screen, float32(bird.x), float32(bird.y), birdWidth, birdHeight, color.RGBA{255, 0, 0, 255}, true) // Red for birds
	}

	// Draw Score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", score))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Initialize Ebiten window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chrome Dino Game - Go Version")

	// Initialize the game loop
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
