package chip8

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Drive execution, manage timers and load games
// as well as manage rendering
type Engine[T Drawable] struct {
	Chip *Chip8
	Graphics T
}

func NewEngine[T Drawable](Renderer T) *Engine[T] {
	var e Engine[T]
	e.Chip = NewChip8()
	e.Graphics = Renderer
	return &e
}

func (e *Engine[_]) LoadGame(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	const memOffset = 0x200
	for idx, b := range data {
		e.Chip.Memory[memOffset+idx] = b
	}

	fmt.Printf("Read %d bytes from %s\n", len(data), path)
}

func (e *Engine[_]) Start() {
	e.Chip.PC = 0x200
	running := true
	for running {
		ins, err := e.Chip.ReadInstruction()

		switch ins.OpCode() {
		case 0x0:
			if ins.N3() == 0xE {
				fmt.Println("C8|Clearing screen.")
				e.Chip.ClearScreen()
			}
		case 0x1:
			e.Chip.Jump(ins.MemAddr())
		case 0x6:
			e.Chip.SetRegister(ins.N2(), int8(ins.Higher))
		case 0x7:
			e.Chip.AddValue(ins.N2(), int8(ins.Higher))
		case 0xA:
			e.Chip.SetIndex(int16(ins.MemAddr()))
		case 0xD:
			//fmt.Printf("C8|Drawing: %d, %d, %d", ins.N2(), ins.N3(), ins.N4())
			e.Chip.Draw(ins.N2(), ins.N3(), ins.N4())
			e.Graphics.Draw(e.Chip.Display)
		}

		if err != nil {
			log.Fatal(err.Error())
		}
		// the longer the wait here is, the longer it takes to crash. hm
		time.Sleep(100 * time.Millisecond)
	}
}
