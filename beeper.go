package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

type Beeper struct {
	audioBuffer *beep.Buffer
}

func NewBeeper() (*Beeper, error) {
	beepFilePath := "assets/audio/beep.mp3"

	f, err := os.Open(beepFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create beeper: %w", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to create beeper: %w", err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	bufferSize := format.SampleRate.N(time.Second / 10)
	speaker.Init(format.SampleRate, bufferSize)

	streamer.Close()
	f.Close()

	return &Beeper{
		audioBuffer: buffer,
	}, nil
}

func (b *Beeper) Beep() {
	sound := b.audioBuffer.Streamer(0, b.audioBuffer.Len())
	speaker.Play(sound)
}
