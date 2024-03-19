package chip8

import(
	"testing"
)

func TestReadInstruction(t *testing.T) {
	var memory [4096]byte
	LoadFont(&memory)

	// OP:f|0|9|0|&>144
	ins, err := ReadInstruction(&memory, 0x050)

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	t.Log(ins.String())
}
