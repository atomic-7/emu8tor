package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/chip8"
)

func main() {
	romBase := "../../chip8-test-suite/bin/"
	fmt.Println("Emulating...")
	engine := chip8.NewEngine()
	engine.LoadGame(romBase + "1-chip8-logo.ch8")
	engine.Start()
}

