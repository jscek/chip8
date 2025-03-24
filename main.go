package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <rom>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Provide a path to the ROM to load.\n")
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	romPath := flag.Arg(0)
	rom, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File '%s' doesn't exist\n", romPath)
		os.Exit(1)
	}

	vm := NewVM()
	vm.Load(rom)

	pixelgl.Run(vm.Run)
}
