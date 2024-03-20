package chip8

import(
	"testing"
)

func TestReadInstruction(t *testing.T) {
	ch := NewChip8()

	// OP:f|0|9|0|&>144
	ch.PC = 0x050
	ins, err := ch.ReadInstruction()

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	t.Log(ins.String())
}
