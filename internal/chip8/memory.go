package chip8

import (
	"fmt"
)

type Memory struct {
	Size int
	Buf []byte
}

func NewMemory(size int) *Memory {
	var mem Memory
	mem.Size = size
	mem.Buf = make([]byte, size, size)
	return &mem
}

func (mem *Memory) PrintInfo() {
	fmt.Printf("%d\n", mem.Size)
}
