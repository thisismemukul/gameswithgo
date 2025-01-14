package models

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Paddle struct {
	Rect       Rectangle
	IsDragging bool
	DragOffset int
}

func (p *Paddle) HandleInput(upKey, downKey ebiten.Key, paddleSpeed int, touchIDs []ebiten.TouchID, ballSpeed int) {

	// paddleSpeed := config.PaddleMoveSpeed + int(math.Abs(float64(ballSpeed))/2)
	newPaddleSpeed := Config.PaddleMoveSpeed + int(math.Min(float64(ballSpeed)/3, 10)) // limit paddle speed increase

	paddleSpeed = max(paddleSpeed, newPaddleSpeed)

	if ebiten.IsKeyPressed(upKey) {
		p.Rect.PosY = max(0, p.Rect.PosY-paddleSpeed)
	}
	if ebiten.IsKeyPressed(downKey) {
		p.Rect.PosY = min(Config.WindowHeight-p.Rect.Height, p.Rect.PosY+paddleSpeed)
	}

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !p.IsDragging && y >= p.Rect.PosY && y <= p.Rect.PosY+p.Rect.Height && x >= p.Rect.PosX && x <= p.Rect.PosX+p.Rect.Width {
			p.IsDragging = true
			p.DragOffset = y - p.Rect.PosY
		}

		if p.IsDragging {
			p.Rect.PosY = y - p.DragOffset
			p.Rect.PosY = max(0, min(Config.WindowHeight-p.Rect.Height, p.Rect.PosY))
		}
	} else {
		p.IsDragging = false
	}

	for _, touchID := range touchIDs {
		tx, ty := ebiten.TouchPosition(touchID)

		if p.Rect.PosX < Config.WindowWidth/2 && tx < Config.WindowWidth/2 {
			paddleCenterY := p.Rect.PosY + p.Rect.Height/2
			if ty < paddleCenterY {
				p.Rect.PosY = max(0, p.Rect.PosY-paddleSpeed)
			} else if ty > paddleCenterY {
				p.Rect.PosY = min(Config.WindowHeight-p.Rect.Height, p.Rect.PosY+paddleSpeed)
			}
		} else if p.Rect.PosX > Config.WindowWidth/2 && tx > Config.WindowWidth/2 {
			paddleCenterY := p.Rect.PosY + p.Rect.Height/2
			if ty < paddleCenterY {
				p.Rect.PosY = max(0, p.Rect.PosY-paddleSpeed)
			} else if ty > paddleCenterY {
				p.Rect.PosY = min(Config.WindowHeight-p.Rect.Height, p.Rect.PosY+paddleSpeed)
			}
		}
	}

}
