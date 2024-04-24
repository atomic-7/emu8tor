package chip8

import (
	"fmt"
	"math/rand"
	"log"
)

// 00E0 clear screen
func (c *Chip8) ClearScreen() {
	c.Display = make([]byte, 64*32)
}

func (c *Chip8) SubReturn() {
	pc, err := c.Stack.pop()
	if err != nil {
		log.Printf("0x00EE: %v", c)
		log.Fatal(err.Error())
	}
	c.PC = pc
}

// 0NNN exec machine routine at NNN
func (c *Chip8) ExecRoutine(addr int16) {
	println("Executing subroutine")
}

// 1NNN jump
func (c *Chip8) Jump(addr uint16) {
	c.PC = addr
}
// 2NNN Subroutine
func (c *Chip8) Subroutine(addr uint16) {
	// running the corax test rom results in a stack overflow after the 2X check
	// stack seems to work as intended, so it seems there is an issue with popping values
	err := c.Stack.push(c.PC)
	if err != nil {
		log.Printf("0x2NNN: %v", c)
		log.Fatal(err.Error())
	}
	c.PC = addr
	/* 0x2NNN: PC: 90
	Instruction: 1125
	Registers: [32 0 0 0 0 42 42 0 18 22 27 1 0 0 0 0]
	Stack: &{[90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90 90] 32} */
}
// 6XNN set register VX
func (c *Chip8) SetRegister(idx uint8, val uint8) {
	c.Registers[idx] = val
}

// 7XNN add value NN to register VX, does not set overflow reg VF
func (c *Chip8) AddValue(idx uint8, val uint8) {
	c.Registers[idx] = c.Registers[idx] + val
}

// 8XY4 add value of register Y to register X, save in X does set overflow reg VF
func (c *Chip8) AddRegOverflow(idx uint8, idy uint8) {
	fatNum := int16(c.Registers[idx]) + int16(c.Registers[idy])
	if fatNum > 256 {
		c.Registers[len(c.Registers) - 1] = 1
		c.Registers[idx] = uint8(fatNum - 256)
	} else {
		c.Registers[len(c.Registers) - 1] = 0
		c.Registers[idx] += c.Registers[idy]
	}
}

// subtract valb from vala, return result and overflow
func subOverflow(vala uint8, valb uint8) (uint8, uint8) {
	var overflow uint8
	if vala >= valb {
		overflow = 1
	} else {
		overflow = 0
	}
	return vala - valb, overflow
}

// 8XY5 VX = VX - VY, Overflow is 1 if VX > VY, reverse x, y for 8XY7
func (c *Chip8) SubXYRegOverflow(idx uint8, idy uint8, yx bool) {
	var res, overflow uint8
	if yx {
		res, overflow = subOverflow(c.Registers[idy], c.Registers[idx])
	} else {
		res, overflow = subOverflow(c.Registers[idx], c.Registers[idy])
	}
	c.Registers[idx] = res
	c.Registers[len(c.Registers) - 1] = overflow
}

// 8XY6 Shift in X register with carry bit, respects architecture semantics
func (c *Chip8) Shift(idx uint8, idy uint8, lshift bool) {
	// CHIP8: put vy in vx, then shift vx 1 bit
	// CHIP48/SCHIP: Ignore vy, just shift vx in place
	if c.Architecture == CHIP8 {
		c.Registers[idx] = c.Registers[idy]
	}
	var mask uint8 = 1
	if lshift {
		mask <<= 7
	} 
	if (c.Registers[idx] & mask) > 0 {
		c.Registers[len(c.Registers) - 1] =  1
	} else {
		c.Registers[len(c.Registers) - 1] =  0
	}
	if lshift {
		c.Registers[idx] <<= 1
	} else {
		c.Registers[idx] >>= 1
	}
}

// ANNN set index register I
func (c *Chip8) SetIndex(idx uint16) {
	c.I = idx
}

// BXNN Jump to NNN with offset V0 or VX
func (c *Chip8) JumpOffset(idx uint8, higher uint8) {
	var addr uint16 = uint16(idx)
	addr <<= 8
	addr |= uint16(higher)
	if c.Architecture == CHIP8 {
		c.PC = addr + uint16(c.Registers[0])
	} else {
		c.PC = addr + uint16(c.Registers[idx])
	}
}

// CXNN Put random number & NN in vx
func (c *Chip8) Random(idx uint8, higher uint8) {
	c.Registers[idx] = uint8(rand.Intn(255)) & higher
}

// DXYN display/draw into the display buffer, actual rendering is handled by the engine
func (c *Chip8) Draw(xIDX uint8, yIDX uint8, nSize uint8) {

	// coords wrap, but the drawn sprite is cut off at the edge of the screen
	xcoord := uint16(c.Registers[xIDX]) % c.Width
	ycoord := uint16(c.Registers[yIDX]) % c.Height
	c.Registers[len(c.Registers)-1] = 0
	fmt.Printf("Drawing sprite %d at x:%d, y:%d\n", c.I, xcoord, ycoord)
	
	//line
	var line, idx uint16	// line of the font, position in the byte
	for line = 0; line < uint16(nSize) && ycoord + line < c.Height; line++ {
		spriteData := c.Memory[c.I + line]
		for idx = 0; idx < 8 && xcoord + idx < c.Width; idx++ {
			if FontBitSet(spriteData, idx) {
				// sprites seem to be stored flipped by chip8 programs
				pos := (ycoord + line) * c.Width + xcoord + 8 -idx
				if c.Display[pos] == 0x1 { // display was already set at this coordinate
					c.Registers[len(c.Registers)-1] = 1
					c.Display[pos] = 0x0
				} else {
					c.Display[pos] = 0x1
				}
			}
		}
	}

}
