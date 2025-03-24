package main

import "github.com/faiface/pixel/pixelgl"

type Keypad struct {
	keys map[pixelgl.Button]uint8
}

func NewKeypad() *Keypad {
	return &Keypad{
		keys: map[pixelgl.Button]uint8{
			pixelgl.Key1: 0x1,
			pixelgl.Key2: 0x2,
			pixelgl.Key3: 0x3,
			pixelgl.Key4: 0xC,

			pixelgl.KeyQ: 0x4,
			pixelgl.KeyW: 0x5,
			pixelgl.KeyE: 0x6,
			pixelgl.KeyR: 0xD,

			pixelgl.KeyA: 0x7,
			pixelgl.KeyS: 0x8,
			pixelgl.KeyD: 0x9,
			pixelgl.KeyF: 0xE,

			pixelgl.KeyZ: 0xA,
			pixelgl.KeyX: 0x0,
			pixelgl.KeyC: 0xB,
			pixelgl.KeyV: 0xF,
		},
	}
}
