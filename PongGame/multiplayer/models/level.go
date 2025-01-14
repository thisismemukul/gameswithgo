package models

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
