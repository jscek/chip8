package main

import (
	"log"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

type Beeper struct {
	audioBuffer *beep.Buffer
}

func NewBeeper(filePath string) *Beeper {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	bufferSize := format.SampleRate.N(time.Second / 10)
	speaker.Init(format.SampleRate, bufferSize)

	streamer.Close()
	f.Close()

	return &Beeper{
		audioBuffer: buffer,
	}
}

func (b *Beeper) Beep() {
	sound := b.audioBuffer.Streamer(0, b.audioBuffer.Len())
	speaker.Play(sound)
}
