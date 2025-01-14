package models

type Ball struct {
	Rect           Rectangle
	SpeedX, SpeedY int
}

func (b *Ball) Move() {
	b.Rect.PosX += b.SpeedX
	b.Rect.PosY += b.SpeedY
}
