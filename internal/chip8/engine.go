package chip8

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// Drive execution, manage timers and load games
// as well as manage rendering
type Engine struct {
	Chip     *Chip8
	Graphics Drawable
	Keypad   Keypad
	TickRate int64 // Microseconds between instructions
}

func NewEngine(renderer Drawable, keypad Keypad) *Engine {
	var e Engine
	e.Chip = NewChip8(CHIP48)
	e.Graphics = renderer
	e.Keypad = keypad
	e.TickRate = 1400
	return &e
}

func (e *Engine) LoadGame(path string) {

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

func (e *Engine) Start(keystateChan chan [16]bool, step chan bool) {

	e.Chip.PC = 0x200
	//	tickduration := int(time.Second) / e.TickRate
	var keyState [16]bool
	chiptimer := NewChipTimer()
	timerctx := context.Background()
	defer timerctx.Done()
	beeper := NewChipTimer()
	beeperctx := context.Background()
	defer beeperctx.Done()
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
				if ins.N4() == 0xE {
					e.Chip.SubReturn()
				} else if ins.N4() == 0x0 {
					e.Chip.ClearScreen()
				} else {
					log.Fatal(fmt.Sprintf("Not implemented:%v", ins))
				}
			}
		case 0x1:
			e.Chip.Jump(ins.MemAddr())
		case 0x2:
			e.Chip.Subroutine(ins.MemAddr())
		case 0x3:
			if e.regVal(ins.N2()) == uint8(ins.Higher) { // skip if r[vx] == nn
				e.Chip.PC += 2
			}
		case 0x4:
			if e.regVal(ins.N2()) != uint8(ins.Higher) {
				e.Chip.PC += 2
			}
		case 0x5:
			if e.regVal(ins.N2()) == e.regVal(ins.N3()) {
				e.Chip.PC += 2
			}

		case 0x6:
			e.Chip.SetRegister(ins.N2(), uint8(ins.Higher))
		case 0x7:
			e.Chip.AddValue(ins.N2(), uint8(ins.Higher))
		case 0x8:
			switch ins.N4() {
			case 0x0: // set vx to vy
				e.Chip.SetRegister(ins.N2(), e.regVal(ins.N3()))
			case 0x1: // set vx to vx | vy
				e.Chip.LogicalOr(ins.N2(), ins.N3())
			case 0x2: // set vx to vy & vx
				e.Chip.LogicalAnd(ins.N2(), ins.N3())
			case 0x3: // set vx to vx xor vy
				e.Chip.LogicalXor(ins.N2(), ins.N3())
			case 0x4: // set vx to vx + vy, set carry flag if it overflows 255
				e.Chip.AddRegOverflow(ins.N2(), ins.N3())
			case 0x5:
				e.Chip.SubXYRegOverflow(ins.N2(), ins.N3(), false)
			case 0x6: // test fail, it seems 0x6 and 0xE are switched. The first one tests right shift
				e.Chip.RShift(ins.N2(), ins.N3())
			case 0x7:
				e.Chip.SubXYRegOverflow(ins.N2(), ins.N3(), true)
			case 0xE:
				e.Chip.LShift(ins.N2(), ins.N3())
			}
		case 0x9:
			if e.regVal(ins.N2()) != e.regVal(ins.N3()) {
				e.Chip.PC += 2
			}
		case 0xA:
			e.Chip.SetIndex(ins.MemAddr())
		case 0xB:
			e.Chip.JumpOffset(ins.N2(), uint8(ins.Higher))
		case 0xC:
			e.Chip.Random(ins.N2(), uint8(ins.Higher))
		case 0xD:
			e.Chip.Draw(ins.N2(), ins.N3(), ins.N4())
			e.Graphics.Draw(e.Chip.Display)

		case 0xE:
			keyState = e.Keypad.GetKeyStates()
			// Skip, now waiting on input
			if ins.N3() == 0x9 {
				// skip next if key in vx is pressed right now
				if keyState[e.regVal(ins.N2())] {
					e.Chip.PC += 2
				}
			} else if ins.N3() == 0xA {
				// skip next if key in vx is not pressed right now
				if !keyState[e.regVal(ins.N2())] {
					e.Chip.PC += 2
				}

			} else {
				log.Fatalf("Unkown instruction:\n%v\n%v", ins, e.Chip)
			}

		case 0xF:
			switch ins.Higher {
			case 0x7:
				//set vx to current value of delay timer
				e.Chip.Registers[ins.N2()] = chiptimer.Count
			case 0xA:
				// block execution until a key is pressed
				keyState = e.Keypad.GetKeyStates()
				pressed := false
				for _, key := range keyState {
					if key {
						pressed = true
						break
					}
				}
				if !pressed {
					e.Chip.PC -= 2
				}

			case 0x15:
				// set delay timer to value in vx
				chiptimer.SetTimer(timerctx, ins.N2())
			case 0x18:
				// set sound timer to value in vx
				beeper.SetBeep(beeperctx, ins.N2(), SilentBeeper)
			case 0x1E:
				// add value in vx to index register I, C8 for Amiga did overflow for 0FFF to 1000
				e.Chip.I += uint16(e.regVal(ins.N2()))
			case 0x29:
				e.Chip.I = 0x50 + uint16(e.regVal(ins.N2())) // may need to use the last nibble of the value in the register
			case 0x33:
				e.Chip.BinaryDecimalConversion(ins.N2())
			case 0x55:
				e.Chip.StoreRegisters(ins.N2())
			case 0x65:
				e.Chip.LoadRegisters(ins.N2())
			}
		default:
			log.Fatal(fmt.Sprintf("Not implemented:%v", ins))
		}

		time.Sleep(time.Duration(e.TickRate) * time.Microsecond)
		//time.Sleep(1 * time.Millisecond)

	}
}

// Get the value in the register idx, hopefully this gets inlined
func (e *Engine) regVal(idx uint8) uint8 {
	return e.Chip.Registers[idx]
}

const (
	BTN0 = iota //x
	BTN1        //1
	BTN2        //2
	BTN3        //3
	BTN4        //q
	BTN5        //w
	BTN6        //e
	BTN7        //a
	BTN8        //s
	BTN9        //d
	BTNA        //y
	BTNB        //c
	BTNC        //4
	BTND        //r
	BTNE        //f
	BTNF        //v
)
