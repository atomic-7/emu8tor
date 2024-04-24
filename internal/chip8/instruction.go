package chip8
import(
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

func (ins *Instruction) OpCode() byte {
	return ins.Lower >> 4
}

func (ins *Instruction) N2() uint8 {
	return uint8(ins.Lower & 0x0f)
}

func (ins *Instruction) N3() uint8 {
	return uint8(ins.Higher >> 4)
}

func (ins *Instruction) N4() uint8 {
	return uint8(ins.Higher & 0x0f)
}

func (ins *Instruction) MemAddr() uint16 {
	var addr uint16
	addr = uint16(ins.N2())
	addr = addr << 8
	return addr | uint16(ins.Higher)
}

func (ins *Instruction) String() string {
	return fmt.Sprintf("OP:%x|N1:%x|N2:%x|N3:%x|&>%x", ins.OpCode(), ins.N2(), ins.N3(),ins.N4(), ins.MemAddr())
}

