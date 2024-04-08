package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/chip8"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
)

func engineWorker(rayrender *raygraphics.RayRender) {
	// TODO: abort context
	romBase := "../../chip8-test-suite/bin/"
	// engine := chip8.NewEngine[*chip8.DebugRender](chip8.NewDebugRender(64, 32))
	engine := chip8.NewEngine[*raygraphics.RayRender](rayrender)
	engine.LoadGame(romBase + "1-chip8-logo.ch8")
	engine.Start()
}

func main() {
	fmt.Println("Emulating...")

	bufChan := make(chan []byte)
	rayrender := raygraphics.NewRaylibRender(64, 32, bufChan)
	go engineWorker(rayrender)
	raygraphics.RenderLoop(64, 32, 1024, 512, rayrender)
}

