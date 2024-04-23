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
	Chip     *Chip8
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

func (e *Engine[_]) Start(step chan bool) {
	// TODO: add next instruction channel as an argument to make a rudimentary debugger
	e.Chip.PC = 0x200
	running := true
	for running {
		if step != nil {
			_ = <-step // read blocking from the step channel
		}
		ins, err := e.Chip.ReadInstruction()

		if err != nil {
			log.Fatal(err.Error())
		}
		switch ins.OpCode() {
		case 0x0:
			if ins.N3() == 0xE {
				if ins.N4() == 0xE { // return from function
					pc, err := e.Chip.Stack.pop()
					if err != nil {
						log.Fatal(err.Error())
					}
					e.Chip.PC = pc
				} else {
					fmt.Println("C8|Clearing screen.")
					e.Chip.ClearScreen()
				}
			}
		case 0x1:
			e.Chip.Jump(ins.MemAddr())
		case 0x2:
			err := e.Chip.Stack.push(e.Chip.PC)
			if err != nil {
				log.Fatal(err.Error())
			}
			e.Chip.Jump(ins.MemAddr())
		case 0x3:
			if ins.N2() == int8(ins.Higher) {
				e.Chip.PC += 2
			}
		case 0x4:
			if ins.N2() != int8(ins.Higher) {
				e.Chip.PC += 2
			}
		case 0x5:
			if ins.N2() == ins.N3() {
				e.Chip.PC += 2
			}
		case 0x6:
			e.Chip.SetRegister(ins.N2(), int8(ins.Higher))
		case 0x7:
			e.Chip.AddValue(ins.N2(), int8(ins.Higher))
		case 0x8:
			switch ins.N4() {
			case 0x0: // set vx to vy
				e.Chip.SetRegister(ins.N2(), ins.N3())
			case 0x1: // set vx to vx | vy
				e.Chip.SetRegister(ins.N2(), ins.N2()|ins.N3())
			case 0x2: // set vx to vy & vx
				e.Chip.SetRegister(ins.N2(), ins.N2()&ins.N3())
			case 0x3: // set vx to vx xor vy
				e.Chip.SetRegister(ins.N2(), ins.N2()^ins.N3())
			case 0x4: // set vx to vx + vy, set carry flag if it overflows 255
			case 0x5:
			case 0x6:
			case 0x7:
			case 0xE:
			}
		case 0x9:
			if ins.N2() != ins.N3() {
				e.Chip.PC += 2
			}
		case 0xA:
			e.Chip.SetIndex(int16(ins.MemAddr()))
		case 0xD:
			//fmt.Printf("C8|Drawing: %d, %d, %d", ins.N2(), ins.N3(), ins.N4())
			e.Chip.Draw(ins.N2(), ins.N3(), ins.N4())
			e.Graphics.Draw(e.Chip.Display)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
