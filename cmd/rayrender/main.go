package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
	)

func colorGenerator(bufChan chan []byte) {
	localBuf := make([]byte, 64 * 32)
	for count := 0; count < 200; count++ {
		localBuf[0] = byte(100 + 20 * (count % 3))
		localBuf[1] = byte(100 + 10 * (count % 5))
		localBuf[2] = byte(100 + 30 * (count % 7))
		
		bufChan <- localBuf
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("Done generating")
}

func main() {
	fmt.Println("RAYGRAPHICS")

	render := raygraphics.NewRaylibRender(800, 450)
	bufChan := make(chan []byte)
	lastBuf := make([]byte, 64 * 32)

	rl.InitWindow(800, 450, "RayRender tests")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	go colorGenerator(bufChan)

	for !rl.WindowShouldClose() {
	
		select {
		case buf := <- bufChan:
			copy(lastBuf, buf)	
			render.DrawBuf(lastBuf)
		default:
			render.DrawBuf(lastBuf)
		}
	}
}
