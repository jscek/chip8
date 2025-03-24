package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <rom>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Provide a path to the ROM to load.\n")
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	vm := NewVM()

	beeper := NewBeeper("assets/beep.mp3")

	romPath := flag.Arg(0)
	rom, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File '%s' doesn't exist\n", romPath)
		os.Exit(1)
	}

	vm.Load(rom)

	keypad := NewKeypad()
	disp, err := NewDisplay(&vm.GFX)
	if err != nil {
		panic(err)
	}

	for range vm.cpuClock.C {
		if !disp.win.Closed() {
			for keyboardKey, pixelKey := range keypad.keys {
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

			if vm.beep {
				beeper.Beep()
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

func main() {
	pixelgl.Run(run)
}
