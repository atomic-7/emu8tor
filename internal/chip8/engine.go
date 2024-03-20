package chip8

import(
	"fmt"
	"log"
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

func (e *Engine) LoadGame() {
	// open fh, read bytes to  chip memory
	fmt.Println("Loading Game")
}

func (e *Engine) Start() {	
	e.Chip.PC = 0x050
	running := true
	for count:=0; count <= 10 && running; count++ {
		ins, err := e.Chip.ReadInstruction()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(ins.String())
	}
}
