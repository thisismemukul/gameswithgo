package models

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

var Config = GameConfig{
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
