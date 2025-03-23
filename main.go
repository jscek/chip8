package main

import (
	"os"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	vm := NewVM()

	keymap := map[pixelgl.Button]uint8{
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
	}

	// TODO: Load ROM from CLI arg
	rom, err := os.ReadFile("roms/pong.ch8")
	if err != nil {
		panic(err)
	}
	vm.Load(rom)

	disp, err := NewDisplay(&vm.GFX)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-vm.cpuClock.C:
			if !disp.win.Closed() {
				for keyboardKey, pixelKey := range keymap {
					if disp.win.Pressed(pixelgl.Button(keyboardKey)) {
						vm.KeyInputs[pixelKey] = 1
					} else {
						vm.KeyInputs[pixelKey] = 0
					}
				}

				err := vm.Cycle()
				if err != nil {
					panic(err)
				}

				if vm.draw {
					disp.Draw()
				} else {
					disp.win.UpdateInput()
				}

				continue
			} else {
				return
			}
		}
	}
}

func main() {
	pixelgl.Run(run)
}
