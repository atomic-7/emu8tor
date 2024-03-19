package chip8
import(
	"errors"
	"fmt"
)
/*	X:n2: look up registers v0 through vf
	Y:n3: look up registers v0 through vf
	N:n4: 4bit number
	NN: 8bit number, higher byte
	NNN: 12bit immediate memory address
*/
type Instruction struct {
	Lower byte	// opcode, n2
	Higher byte	// n3, n4
}
	
func ReadInstruction(mem *[4096]byte, pc int16) (Instruction, error) {
	var instruction Instruction
	if pc >= 4096 {
		return instruction, errors.New("pc out of bounds")
	}
	instruction.Lower = mem[pc]
	instruction.Higher = mem[pc+1]
	return instruction, nil
}

func (ins *Instruction) OpCode() byte {
	return ins.Lower >> 4
}

func (ins *Instruction) N2() byte {
	return ins.Lower << 4
}

func (ins *Instruction) N3() byte {
	return ins.Higher >> 4
}

func (ins *Instruction) N4() byte {
	return ins.Higher << 4
}

func (ins *Instruction) MemAddr() int16 {
	var addr int16
	addr = int16(ins.N2())
	addr = addr << 8
	return addr | int16(ins.Higher)
}

func (ins *Instruction) String() string {
	return fmt.Sprintf("OP:%x|%x|%x|%x|&>%d", ins.OpCode(), ins.N2(), ins.N3(),ins.N4(), ins.MemAddr())
}

