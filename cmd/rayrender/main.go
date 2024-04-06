package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
	rl "github.com/gen2brain/raylib-go/raylib"
	)

func main() {
	fmt.Println("RAYGRAPHICS")

	render := raygraphics.NewRaylibRender(800, 450)
	dbuf := make([]byte, 64 * 32)

	rl.InitWindow(800, 450, "RayRender tests")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		
		render.Draw(dbuf)
	}
}
