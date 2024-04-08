package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/chip8"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
)

func engineWorker(rayrender *raygraphics.RayRender, step chan bool) {
	// TODO: abort context
	romBase := "../../chip8-test-suite/bin/"
	// engine := chip8.NewEngine[*chip8.DebugRender](chip8.NewDebugRender(64, 32))
	engine := chip8.NewEngine[*raygraphics.RayRender](rayrender)
	//engine.LoadGame(romBase + "1-chip8-logo.ch8")
	engine.LoadGame(romBase + "2-ibm-logo.ch8")
	engine.Start(step)
}

func main() {
	fmt.Println("Emulating...")

	bufChan := make(chan []byte)
	//stepChan := make(chan bool)
	rayrender := raygraphics.NewRaylibRender(64, 32, bufChan)
	go engineWorker(rayrender, nil)
	raygraphics.RenderLoop(64, 32, 1024, 512, rayrender, nil)
}

