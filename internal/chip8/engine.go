package chip8

import (
	"fmt"
	"io"
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
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	
	totalBytes := 0
	const maxChunkSize = 8
	const memOffset = 0x200
	tmpB := make([]byte, maxChunkSize)
	for {
		numBytes, err := file.Read(tmpB)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break	// eof reached
		}
		for idx := 0; idx < numBytes; idx ++ {
			e.Chip.Memory[memOffset + totalBytes + idx] = tmpB[idx]
		}
		totalBytes += numBytes
	}
	
	fmt.Printf("Read %d bytes from %s\n", totalBytes, path)
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
