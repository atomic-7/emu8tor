package chip8

import (
	"fmt"
)

// 00E0 clear screen
func (c *Chip8) ClearScreen() {
	c.Display = make([]byte, 64*32)
}

// 0NNN exec routine at NNN
func (c *Chip8) ExecRoutine(addr int16) {
	println("Executing subroutine")
}

// 1NNN jump
func (c *Chip8) Jump(addr int16) {
	c.PC = addr
}

// 6XNN set register VX
func (c *Chip8) SetRegister(idx int8, val int8) {
	c.Registers[idx] = val
}

// 7XNN add value NN to register VX, does not set overflow reg VF
func (c *Chip8) AddValue(idx int8, val int8) {
	c.Registers[idx] = c.Registers[idx] + val
}

// ANNN set index register I
func (c *Chip8) SetIndex(idx int16) {
	c.I = idx
}

// DXYN display/draw into the display buffer, acutal rendering is handled by the engine
func (c *Chip8) Draw(xIDX int8, yIDX int8, nSize int8) {

	// coords wrap, but the drawn sprite is cut off at the edge of the screen
	xcoord := int16(c.Registers[xIDX]) % c.Width
	ycoord := int16(c.Registers[yIDX]) % c.Height
	c.Registers[len(c.Registers)-1] = 0
	fmt.Printf("Drawing sprite %d at x:%d, y:%d\n", c.I, xcoord, ycoord)
	
	// dbg
	c.Registers[5] = 1
	c.Display[1] = 0xf

	spriteData := c.Memory[c.I+int16(nSize)] // nSize???

	var idx, line int16 // index of the sprite, offset for the xcoordinate

	for line = 0; ycoord+line < int16(nSize) && ycoord+line < c.Height; line++ {
		for idx = 0; idx < 8 && xcoord+idx < c.Width; idx++ {
			if FontBitSet(spriteData, idx) {
				if c.Display[(ycoord+line)*c.Width+xcoord+idx] == 0x1 { // display was already set at this coordinate
					c.Registers[len(c.Registers)-1] = 1
					c.Display[(ycoord+line)*c.Width+xcoord+idx] = 0x0
				} else {
					c.Display[(ycoord+line)*c.Width+xcoord+idx] = 0x1
				}
			}
		}
	}
}
