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
func (c *Chip8) Jump(addr int16) {
	c.PC = addr
}
// 2NNN Subroutine
func (c *Chip8) Subroutine(addr int16) {
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
func (c *Chip8) SetRegister(idx int8, val int8) {
	c.Registers[idx] = val
}

// 7XNN add value NN to register VX, does not set overflow reg VF
func (c *Chip8) AddValue(idx int8, val int8) {
	c.Registers[idx] = c.Registers[idx] + val
}

// 8XY4 add value of register Y to register X, save in X does set overflow reg VF
func (c *Chip8) AddRegOverflow(idx int8, idy int8) {
	fatNum := int16(c.Registers[idx]) + int16(c.Registers[idy])
	if fatNum > 255 {
		c.Registers[len(c.Registers) - 1] = 1
		c.Registers[idx] = int8(fatNum - 255)
	} else {
		c.Registers[len(c.Registers) - 1] = 0
		c.Registers[idx] += c.Registers[idy]
	}
}

// subtract valb from vala, return result and overflow
func subOverflow(vala int8, valb int8) (int8, int8) {
	var overflow int8
	if vala >= valb {
		overflow = 1
	} else {
		overflow = 0
	}
	return vala - valb, overflow
}

// 8XY5 VX = VX - VY, Overflow is 1 if VX > VY, reverse x, y for 8XY7
func (c *Chip8) SubXYRegOverflow(idx int8, idy int8) {
	res, overflow := subOverflow(c.Registers[idx], c.Registers[idy])
	c.Registers[idx] = res
	c.Registers[len(c.Registers) - 1] = overflow
}

// 8XY6 Shift in X register with carry bit, respects architecture semantics
func (c *Chip8) Shift(idx int8, idy int8, lshift bool) {
	// CHIP8: put vy in vx, then shift vx 1 bit
	// CHIP48/SCHIP: Ignore vy, just shift vx in place
	if c.Architecture == CHIP8 {
		c.Registers[idx] = c.Registers[idy]
		var mask int8 = 1
		if lshift {
			// think about mask size/signed/unsigned, could make a difference
			mask <<= 7
		} 
		if (c.Registers[idx] & mask) > 0 {
			c.Registers[len(c.Registers) - 1] =  1
		} else {
			c.Registers[len(c.Registers) - 1] =  0
		}
	}
	//var shiftedBit int8
	if lshift {
		c.Registers[idx] <<= 1
	} else {
		c.Registers[idx] >>= 1
	}
}

// ANNN set index register I
func (c *Chip8) SetIndex(idx int16) {
	c.I = idx
}

// BXNN Jump to NNN with offset V0 or VX
func (c *Chip8) JumpOffset(idx int8, higher int8) {
	var addr int16 = int16(idx)
	addr <<= 8
	addr |= int16(higher)
	if c.Architecture == CHIP8 {
		c.PC = addr + int16(c.Registers[0])
	} else {
		c.PC = addr + int16(c.Registers[idx])
	}
}

// CXNN Put random number & NN in vx
func (c *Chip8) Random(idx int8, higher int8) {
	c.Registers[idx] = int8(rand.Intn(255)) & higher
}

// DXYN display/draw into the display buffer, actual rendering is handled by the engine
func (c *Chip8) Draw(xIDX int8, yIDX int8, nSize int8) {

	// coords wrap, but the drawn sprite is cut off at the edge of the screen
	xcoord := int16(c.Registers[xIDX]) % c.Width
	ycoord := int16(c.Registers[yIDX]) % c.Height
	c.Registers[len(c.Registers)-1] = 0
	fmt.Printf("Drawing sprite %d at x:%d, y:%d\n", c.I, xcoord, ycoord)
	
	//line
	var line, idx int16	// line of the font, position in the byte
	for line = 0; line < int16(nSize) && ycoord + line < c.Height; line++ {
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
