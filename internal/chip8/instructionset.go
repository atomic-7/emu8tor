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
