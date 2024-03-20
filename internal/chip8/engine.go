package chip8

import (
	"fmt"
	"log"
	"os"
)

// Drive execution, manage timers and load games
// as well as manage rendering
type Engine struct {
	Chip *Chip8
}

func  NewEngine() *Engine {
	var e Engine
	e.Chip = NewChip8()
	return &e
}

func (e *Engine) LoadGame(path string) {
	// open fh, read bytes to  chip memory

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	const memOffset = 0x200
	for idx, b := range data {
		e.Chip.Memory[memOffset + idx] = b
	}

	fmt.Printf("Read %d bytes from %s\n", len(data), path)
}

func (e *Engine) Start() {	
	e.Chip.PC = 0x200
	running := true
	for count:=0; count <= 10 && running; count++ {
		ins, err := e.Chip.ReadInstruction()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(ins.String())
	}
}
