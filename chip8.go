package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type VM struct {
	opcode uint16

	// RAM
	memory [4096]uint8

	// General purpose 8-bit registers
	v [16]uint8

	// Address register
	i uint16

	// Program counter
	pc uint16

	// Stack pointer
	sp uint8

	// Stack
	stack [16]uint16

	gfx [32 * 8]uint8

	// Delay Timer
	dt uint8

	// Sound Timer
	st uint8

	draw bool

	beep bool

	cpuClock *time.Ticker

	keyInputs [16]uint8
}

func (vm *VM) loadSpriteByte(x uint8, y uint8, spriteByte uint8) uint8 {
	yPos := y % 32
	xPos := x % 64
	collision := uint8(0)

	for bit := uint8(0); bit < 8; bit++ {
		if (spriteByte & (0x80 >> bit)) != 0 {
			// Calculate the position in the GFX array
			arrayPos := (yPos * 8) + ((xPos + bit) / 8)
			bitPos := 7 - ((xPos + bit) % 8)

			// Check if we're going to flip a set bit (collision)
			if (vm.gfx[arrayPos] & (1 << bitPos)) != 0 {
				collision = 1
			}

			// XOR the bit
			vm.gfx[arrayPos] ^= (1 << bitPos)
		}
	}

	return collision
}

func (vm *VM) fetchOpcode() {
	b1, b2 := vm.memory[vm.pc], vm.memory[vm.pc+1]
	vm.opcode = uint16(b1)<<8 | uint16(b2)
}

func (vm *VM) executeOpcode() error {
	nnn := uint16(vm.opcode & 0x0FFF)
	n := uint8(vm.opcode & 0x000F)
	x := uint8(vm.opcode & 0x0F00 >> 8)
	y := uint8(vm.opcode&0x00F0) >> 4
	kk := uint8(vm.opcode & 0x00FF)

	switch vm.opcode & 0xF000 {
	case 0x0000:
		switch vm.opcode {
		case 0x00E0:
			vm.gfx = [32 * 8]uint8{}
			vm.pc += 2

		case 0x00EE:
			vm.pc = vm.stack[vm.sp]
			vm.sp--
			vm.pc += 2

		default:
			return fmt.Errorf("invalid opcode: %04X", vm.opcode)
		}

	case 0x1000:
		vm.pc = nnn

	case 0x2000:
		vm.sp++
		vm.stack[vm.sp] = vm.pc
		vm.pc = nnn

	case 0x3000:
		if vm.v[x] == kk {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0x4000:
		if vm.v[x] != kk {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0x5000:
		if vm.v[x] == vm.v[y] {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0x6000:
		vm.v[x] = kk
		vm.pc += 2

	case 0x7000:
		vm.v[x] = vm.v[x] + kk
		vm.pc += 2

	case 0x8000:
		switch vm.opcode & 0x000F {
		case 0x0:
			vm.v[x] = vm.v[y]
			vm.pc += 2

		case 0x1:
			vm.v[x] = vm.v[x] | vm.v[y]
			vm.pc += 2

		case 0x2:
			vm.v[x] = vm.v[x] & vm.v[y]
			vm.pc += 2

		case 0x3:
			vm.v[x] = vm.v[x] ^ vm.v[y]
			vm.pc += 2

		case 0x4:
			sum := uint16(vm.v[x] + vm.v[y])
			carry := (sum & 0x0F00) >> 8
			vm.v[x] = uint8(sum & 0x00FF)
			vm.v[0xF] = uint8(carry)
			vm.pc += 2

		case 0x5:
			if vm.v[x] > vm.v[y] {
				vm.v[0xF] = 1
			} else {
				vm.v[0xF] = 0
			}

			vm.v[x] -= vm.v[y]
			vm.pc += 2

		case 0x6:
			// TODO: can be done better?
			vm.v[0xF] = vm.v[x] & 0b00000001
			vm.v[x] = vm.v[x] >> 1
			vm.pc += 2

		case 0x7:
			if vm.v[y] > vm.v[x] {
				vm.v[0xF] = 1
			} else {
				vm.v[0xF] = 0
			}

			vm.v[x] = vm.v[y] - vm.v[x]
			vm.pc += 2

		case 0xE:
			// TODO: can be done better?
			if vm.v[x]&0b10000000 == 0b10000000 {
				vm.v[0xF] = 1
			} else {
				vm.v[0xF] = 0
			}

			vm.v[x] = vm.v[x] << 1
			vm.pc += 2

		default:
			return fmt.Errorf("invalid opcode: %04X", vm.opcode)
		}

	case 0x9000:
		if vm.v[x] != vm.v[y] {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0xA000:
		vm.i = nnn
		vm.pc += 2

	case 0xB000:
		vm.pc = nnn + uint16(vm.v[0])

	case 0xC000:
		r := uint8(rand.Int31())
		vm.v[x] = kk & r
		vm.pc += 2

	case 0xD000:
		xPixel := vm.v[x]
		yPixel := vm.v[y]

		vm.v[0xF] = 0 // Reset collision flag

		for i := range n {
			spriteByte := vm.memory[vm.i+uint16(i)]
			collision := vm.loadSpriteByte(xPixel, yPixel+i, spriteByte)

			if collision == 1 {
				vm.v[0xF] = 1
			}
		}

		vm.draw = true
		vm.pc += 2

	case 0xE000:
		switch vm.opcode & 0x00FF {
		case 0x9E:
			if vm.keyInputs[vm.v[x]] != 0 {
				vm.pc += 4
			} else {
				vm.pc += 2
			}

		case 0xA1:
			if vm.keyInputs[vm.v[x]] == 0 {
				vm.pc += 4
			} else {
				vm.pc += 2
			}

		default:
			return fmt.Errorf("invalid opcode: %04X", vm.opcode)
		}

	case 0xF000:
		switch vm.opcode & 0x00FF {
		case 0x07:
			vm.v[x] = vm.dt
			vm.pc += 2

		case 0x0A:
			for key, val := range vm.keyInputs {
				if val != 0 {
					vm.v[x] = uint8(key)
					vm.pc += 2
					break
				}
			}
			vm.keyInputs[vm.v[x]] = 0

		case 0x15:
			vm.dt = vm.v[x]
			vm.pc += 2

		case 0x18:
			vm.st = vm.v[x]
			vm.pc += 2

		case 0x1E:
			vm.i = vm.i + uint16(vm.v[x])
			vm.pc += 2

		case 0x29:
			digit := vm.v[x]
			vm.i = uint16(digit * 5)
			vm.pc += 2

		case 0x33:
			val := vm.v[x]

			vm.memory[vm.i] = val / 100
			vm.memory[vm.i+1] = (val % 100) / 10
			vm.memory[vm.i+2] = val % 10

			vm.pc += 2

		case 0x55:
			for vi := uint16(0); vi <= uint16(x); vi++ {
				vm.memory[vm.i+vi] = vm.v[vi]
			}
			vm.pc += 2

		case 0x65:
			for vi := uint16(0); vi <= uint16(x); vi++ {
				vm.v[vi] = vm.memory[vm.i+vi]
			}
			vm.pc += 2

		default:
			return fmt.Errorf("invalid opcode: %04X", vm.opcode)
		}
	}

	return nil
}

func NewVM() *VM {
	cpuFreq := 500.0

	vm := &VM{
		pc:       0x200,
		draw:     true,
		cpuClock: time.NewTicker(time.Duration(float64(time.Second) / cpuFreq)),
	}

	fonts := [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	}

	copy(vm.memory[:], fonts[:])

	return vm
}

func (vm *VM) Load(rom []byte) {
	copy(vm.memory[0x200:], rom)
}

func (vm *VM) Cycle() error {
	vm.draw = false
	vm.beep = false
	vm.fetchOpcode()

	if vm.dt > 0 {
		vm.dt--
	}
	if vm.st > 0 {
		vm.st--

		if vm.st == 0 {
			vm.beep = true
		}
	}

	return vm.executeOpcode()
}

func (vm *VM) Run() {
	beeper, err := NewBeeper("assets/beep.mp3")
	if err != nil {
		panic(err)
	}

	keypad := NewKeypad()
	disp, err := NewDisplay(&vm.gfx)
	if err != nil {
		panic(err)
	}

	for range vm.cpuClock.C {
		if !disp.Closed() {
			for keyboardKey, pixelKey := range keypad.keys {
				if disp.win.Pressed(pixelgl.Button(keyboardKey)) {
					vm.keyInputs[pixelKey] = 1
				} else {
					vm.keyInputs[pixelKey] = 0
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
			}

			disp.UpdateInput()
		} else {
			return
		}
	}
}
