package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/chip8"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
)

func engineWorker(rayrender *raygraphics.RayRender, keystateChan chan [16]bool, step chan bool) {
	// TODO: abort context
	romBase := "../../chip8-test-suite/bin/"
	// engine := chip8.NewEngine[*chip8.DebugRender](chip8.NewDebugRender(64, 32))
	// rayrender handles both rendering and input grabbing
	engine := chip8.NewEngine(rayrender, rayrender)
	//engine.LoadGame(romBase + "1-chip8-logo.ch8")
	//engine.LoadGame(romBase + "2-ibm-logo.ch8")
	//engine.LoadGame(romBase + "3-corax+.ch8")
	engine.LoadGame(romBase + "6-keypad.ch8")
	//engine.LoadGame("../../chip8-test-rom/test_opcode.ch8")
	engine.Start(keystateChan, step)
}

func main() {
	fmt.Println("Emulating...")

	bufChan := make(chan []byte)
	keystateChan := make(chan [16]bool)
	//stepChan := make(chan bool)
	rayrender := raygraphics.NewRaylibRender(64, 32, bufChan)
	go engineWorker(rayrender, keystateChan, nil)
	raygraphics.RenderLoop(64, 32, 1024, 512, rayrender,keystateChan, nil)
}

