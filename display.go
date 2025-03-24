package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Display struct {
	pixels *[32 * 8]uint8
	win    *pixelgl.Window
}

func NewDisplay(pixels *[32 * 8]uint8) (*Display, error) {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip-8",
		Bounds: pixel.R(0, 0, 1536, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}

	display := &Display{
		pixels: pixels,
		win:    win,
	}

	return display, nil
}

func (display *Display) Closed() bool {
	return display.win.Closed()
}

func (display *Display) UpdateInput() {
	display.win.UpdateInput()
}

func (display *Display) Draw() {
	display.win.Clear(colornames.Black)

	imd := imdraw.New(nil)

	for y := range 32 {
		for x := range 64 {
			xByte := display.pixels[(y*8)+(x/8)]

			val := uint8(xByte>>(8-1-(x%8))) & 1

			if val == 1 {
				imd.Color = colornames.Lightblue
			} else {
				imd.Color = colornames.Darkblue
			}

			xRes := x * 24
			yRes := display.win.Bounds().H() - float64((y+1)*24)

			imd.Push(pixel.V(float64(xRes), float64(yRes)))
			imd.Push(pixel.V(float64(xRes+24), float64(yRes+24)))
			imd.Rectangle(0)
		}
	}

	imd.Draw(display.win)
	display.win.Update()
}
