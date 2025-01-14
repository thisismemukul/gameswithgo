package gameaudio

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/thisismemukul/pong/consts"
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

var (
	audioContext             = audio.NewContext(consts.SampleRate)
	BackgroundPlayer         *audio.Player
	HitPaddleBallSoundPlayer *audio.Player
	HitWallBallSoundPlayer   *audio.Player
	OutOfBoundarySoundPlayer *audio.Player
	GameOverSoundDataPlayer  *audio.Player
	// infiniteLoop             *audio.InfiniteLoop
)

func InitAudio() {
	// Load background music
	stream, err := mp3.DecodeWithSampleRate(consts.SampleRate, bytes.NewReader(backgroundMusicData))
	if err != nil {
		log.Fatalf("Failed to decode background music: %v", err)
	}
	BackgroundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create background music player: %v", err)
	}
	// backgroundPlayer.SetLoop(true)
	// if infiniteLoop != nil {
	// 	backgroundPlayer.Play()
	// }

	// Load hit paddle ball sound
	stream, err = mp3.DecodeWithSampleRate(consts.SampleRate, bytes.NewReader(hitPaddleBallSoundData))
	if err != nil {
		log.Fatalf("Failed to decode hit paddle ball sound: %v", err)
	}
	HitPaddleBallSoundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create hit paddle ball sound player: %v", err)
	}

	// Load kit wall ball sound
	stream, err = mp3.DecodeWithSampleRate(consts.SampleRate, bytes.NewReader(hitWallBallSoundData))
	if err != nil {
		log.Fatalf("Failed to decode hit wall sound: %v", err)
	}
	HitWallBallSoundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create hit wall sound player: %v", err)
	}

	// Load ball out of boundary sound
	stream, err = mp3.DecodeWithSampleRate(consts.SampleRate, bytes.NewReader(outOfBoundarySoundData))
	if err != nil {
		log.Fatalf("Failed to decode ball out of boundary sound: %v", err)
	}
	OutOfBoundarySoundPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create ball out of boundary sound player: %v", err)
	}

	// Load game over sound
	stream, err = mp3.DecodeWithSampleRate(consts.SampleRate, bytes.NewReader(gameOverSoundData))
	if err != nil {
		log.Fatalf("Failed to decode hit wall sound: %v", err)
	}
	GameOverSoundDataPlayer, err = audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create hit wall sound player: %v", err)
	}
}
