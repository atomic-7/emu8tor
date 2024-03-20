package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/chip8"
)

func main() {
	fmt.Println("Emulating...")
	engine := chip8.NewEngine()
	engine.LoadGame()
	engine.Start()
}

