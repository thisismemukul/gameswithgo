package models

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Paddle struct {
	Rect Rectangle
}

// func (p *Paddle) HandleInput(upKey, downKey ebiten.Key, speed int) {
// 	if ebiten.IsKeyPressed(upKey) && p.Rect.PosY > 0 {
// 		p.Rect.PosY -= speed
// 	}
// 	if ebiten.IsKeyPressed(downKey) && p.Rect.PosY < 480-p.Rect.Height {
// 		p.Rect.PosY += speed
// 	}
// }

// Paddle and Ball Movement
func (p *Paddle) HandleInput(upKey, downKey ebiten.Key, paddleSpeed int, ballSpeed int) {

	// paddleSpeed := config.PaddleMoveSpeed + int(math.Abs(float64(ballSpeed))/2)
	newPaddleSpeed := Config.PaddleMoveSpeed + int(math.Min(float64(ballSpeed)/3, 10)) // limit paddle speed increase

	paddleSpeed = max(paddleSpeed, newPaddleSpeed)

	if ebiten.IsKeyPressed(upKey) {
		p.Rect.PosY = max(0, p.Rect.PosY-paddleSpeed)
	}
	if ebiten.IsKeyPressed(downKey) {
		p.Rect.PosY = min(Config.WindowHeight-p.Rect.Height, p.Rect.PosY+paddleSpeed)
	}
}
