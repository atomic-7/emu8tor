package chip8

import (
	"errors"
	"fmt"
)

// Chip8 should represent the memory as seen by an original chip8 game
type Chip8 struct {
	Memory    []byte
	Registers []uint8
	Display   []byte
	Stack     *Stack
	PC        uint16
	I         uint16
	Width     uint16
	Height    uint16
	Architecture int
}

const (
	CHIP8 = iota
	CHIP48
	SCHIP
)

// both pc and index counter I can only adress 12 bits = 4096 addresses

func NewChip8(instructionset int) *Chip8 {
	var chip Chip8
	chip.Width = 64
	chip.Height = 32
	chip.Memory = make([]byte, 4096, 4096)
	LoadFont(chip.Memory)
	chip.Registers = make([]uint8, 16, 16)
	chip.Display = make([]byte, chip.Width*chip.Height, chip.Width*chip.Height)
	chip.Stack = &Stack{
		sp: 0,
	}
	chip.PC = 0
	chip.I = 0x020
	chip.Architecture = instructionset
	return &chip
}

// Read an instruction and automatically increase the program counter
func (c *Chip8) ReadInstruction() (Instruction, error) {
	var instruction Instruction
	if c.PC >= 4096 {
		return instruction, errors.New("pc out of bounds")
	}
	instruction.Lower = c.Memory[c.PC]
	instruction.Higher = c.Memory[c.PC+1]
	// automatically increase after every read
	c.PC = c.PC + 2
	return instruction, nil
}

// toString
func (c Chip8) String() string {
	return fmt.Sprintf("PC: %d\nInstruction: %d\nRegisters: %v\nStack: %v", c.PC, c.I, c.Registers, c.Stack)
}
