package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type CPU struct {
	G             Graphics
	S             Speaker
	Memory        [4096]uint8
	V             [16]uint8
	MemoryAddress int
	DelayTimer    int
	SoundTimer    int
	PC            int
	Stack         []int
	Speed         int
	Paused        bool
	I             uint16
}

var STEP int = 2
var sprites = []uint8{
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
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func (c *CPU) loadSpritesToMemory() {
	copy(c.Memory[:len(sprites)], sprites)
	// for i, sprite := range sprites {
	// 	c.Memory[i] = sprite
	// }
}

func (c *CPU) loadProgramIntoMemory(program []uint8) {
	for i, programVal := range program {
		c.Memory[i+0x200] = programVal
	}
}

func (c *CPU) loadRom(romPath string) error {
	f, err := os.Open(romPath)
	if err != nil {
		return err
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return err
	}
	program := make([]byte, stats.Size())
	reader := bufio.NewReader(f)
	reader.Read(program)

	c.loadProgramIntoMemory(program)
	return nil

}

func (c *CPU) Execute(opcode uint16) error {
	c.PC += STEP
	// nnn, the lowest 12 bits
	nnn := 0x0fff & opcode
	// n, lowest 4 bits
	// n := 0x000f & opcode

	// x, lower 4 bits of the first half
	y := (0x00f0 & opcode) >> 4
	// y, upper 4 bits of the first half
	x := (0x0f00 & opcode) >> 8

	// kk, the lowest 8 bits
	kk := uint8(0x00ff & opcode)
	vF := uint8(0)
	// fmt.Printf("opcode is %04x, lowest 12 bits is %03x, lower 4 bits of first half is %x, upper 4 bits of first half is %x, kk is %02x\n", opcode, nnn, x, y, kk)
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			c.G.Clear()
		case 0x00EE:
			c.PC = c.Stack[len(c.Stack)-1]
			c.Stack = c.Stack[:len(c.Stack)-1]
		}

	case 0x1000:
		// JUMP addr
		c.PC = int(nnn)
	case 0x2000:
		// CALL addr
		c.Stack = append(c.Stack, c.PC)
		c.PC = int(nnn)
	case 0x3000:
		// SKIP IF Register[x] = kk
		if c.V[x] == kk {
			c.PC += STEP
		}
	case 0x4000:
		// SKIP IF Register[x] != kk
		if c.V[x] != kk {
			c.PC += STEP
		}

	case 0x5000:
		// SKIP if Register[x] == Register[y]
		if c.V[x] == c.V[y] {
			c.PC += STEP
		}

	case 0x6000:
		// Set Register[x] to kk
		// fmt.Printf("Setting register %d to value %d\n", x, kk)
		c.V[x] = kk

	case 0x7000:
		// Increment Vx by kk
		c.V[x] += kk

	case 0x8000:
		switch opcode & 0xF {
		case 0x0:
			// Set Vx to Vy
			c.V[x] = c.V[y]

		case 0x1:
			// Set Vx to Vx OR Vy
			c.V[x] = c.V[x] | c.V[y]

		case 0x2:
			// Set Vx to Vx AND Vy
			c.V[x] = c.V[x] & c.V[y]

		case 0x3:
			// Set Vx to Vx XOR Vy
			c.V[x] = c.V[x] ^ c.V[y]

		case 0x4:
			// Add Vx and Vy. Store 1 in VF if the result is greater than 8 bits
			sum := uint16(c.V[x]) + uint16(c.V[y])
			if sum > 0xff {
				vF = 1
			}
			c.V[x] += c.V[y]
			c.V[0xf] = vF

		case 0x5:
			// Perform Vx -= Vy. Store 1 in Vf if Vx > Vy
			if c.V[x] > c.V[y] {
				vF = 1
			}
			c.V[x] -= c.V[y]
			c.V[0xf] = vF

		case 0x6:
			// Shift right. If the least significant bit of Vx is 1, then VF is set to 1
			if 0x1&c.V[x] != 0 {
				vF = 1
			}
			c.V[x] = c.V[x] >> 1
			c.V[0xf] = vF

		case 0x7:
			// Inverse subtract. Set Vx = Vy - Vx
			if c.V[y] > c.V[x] {
				vF = 1
			}
			c.V[x] = c.V[y] - c.V[x]
			c.V[0xf] = vF

		case 0xE:
			// Shift left. Set Vx to Vx << 1. Set VF to 1 if the MSB of Vx is 1
			if c.V[x]&0b10000000 != 0 {
				vF = 1
			}
			c.V[x] = c.V[x] << 1
			c.V[0xf] = vF

		}

	case 0x9000:
		// Skip next instruction if Vx != Vy
		if c.V[x] != c.V[y] {
			c.PC += 2
		}

	case 0xA000:
		// Set I to nnn
		c.I = nnn

	case 0xB000:
		// Set PC to nnn + V0
		c.PC = int(nnn) + int(c.V[0])

	case 0xC000:
		// AND a random number from 0-255 inclusive with the lowest byte of the opcode
		value := rand.Intn(256)
		c.V[x] = kk & uint8(value)

	case 0xD000:
		// DRAW
		// always width of 8
		width := 8
		// height is the last nibble
		height := int(opcode & 0xf)

		for row := 0; row < height; row++ {
			sprite := c.Memory[int(c.I)+row]
			for col := 0; col < width; col++ {
				if sprite&0x80 > 0 {
					// we should toggle that value
					// fmt.Printf("y is %d, cY is %d\n", y, c.V[y])
					erased, err := c.G.Toggle(int(c.V[x])+col, int(c.V[y])+row)
					if err != nil {
						return err
					}
					if erased {
						vF = 1
					}
				}
				sprite = sprite << 1
			}
		}
		c.V[0xf] = vF

	case 0xE000:
		switch opcode & 0xFF {
		case 0x9E:

		case 0xA1:

		}

	case 0xF000:
		switch opcode & 0xFF {
		case 0x07:
			c.V[x] = uint8(c.DelayTimer)

		case 0x0A:

		case 0x15:
			c.DelayTimer = int(c.V[x])

		case 0x18:
			c.SoundTimer = int(c.V[x])

		case 0x1E:
			c.I += uint16(c.V[x])

		case 0x29:
			// Set I to the location of the sprite stored in Vx. Each sprite is 5 bytes long, so mult by 5.
			c.I = uint16(c.V[x] * 5)

		case 0x33:
			// for Vx set the hundred digits in I, tens digit in I+1, unit digit in I+2
			unit := c.V[x] % 10
			tens := ((c.V[x] - unit) / 10) % 10
			hundreds := ((c.V[x] - unit - tens*10) / 100) % 10
			c.Memory[c.I+2] = unit
			c.Memory[c.I+1] = tens
			c.Memory[c.I] = hundreds

		case 0x55:
			// set c.Memory[c.I:c.I+x] = c.V[0:x] (inclusive to x)
			for i := uint16(0); i <= x; i++ {
				c.Memory[c.I+i] = c.V[i]
			}

		case 0x65:
			// set c.V[0:x] to c.Memory[c.I:c.I+x] (inclusive to x)
			for i := uint16(0); i <= x; i++ {
				c.V[i] = c.Memory[c.I+i]
			}

		}

	default:
		return fmt.Errorf("unknown opcode %d", opcode)
	}

	return nil
}

func (c *CPU) updateTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}
}

func (c *CPU) playSounds() {
	FREQUENCY := 440
	if c.SoundTimer > 0 {
		c.S.Play(FREQUENCY)
	} else {
		c.S.Stop()
	}
}

func (c *CPU) PrintProgramMemory() {
	for i := 0; i < len(c.Memory); i += 2 {
		fmt.Printf("%#04x %04X\n", i, binary.BigEndian.Uint16(c.Memory[i:]))
	}
}
func (c *CPU) Cycle() error {
	for i := 0; i < c.Speed; i++ {
		if !c.Paused {
			// we load two bytes as the opcode
			opcode := binary.BigEndian.Uint16(c.Memory[c.PC : c.PC+2])
			// secondHalf := uint16(c.Memory[c.PC+1])
			// firstHalf := uint16(c.Memory[c.PC])
			// opcode := firstHalf<<8 | secondHalf
			// fmt.Printf("Opcode: %x\n", opcode)

			err := c.Execute(opcode)
			if err != nil {
				return err
			}
		}
	}
	c.updateTimers()
	c.playSounds()
	c.G.Render()
	return nil
}

func (c *CPU) Init(filepath string) error {
	c.loadSpritesToMemory()
	err := c.loadRom(filepath)
	c.PC = 0x200
	if err != nil {
		return err
	}
	return nil
}

func (c *CPU) Loop() {
	fps := 10
	fpsInterval := time.Duration(1000 / fps)
	tick := time.Tick(time.Millisecond * fpsInterval)
	for range tick {
		c.Cycle()
	}
}
