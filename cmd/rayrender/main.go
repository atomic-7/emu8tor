package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
	)

func colorGenerator(bufChan chan []byte) {
	localBuf := make([]byte, 64 * 32)
	for count := 0; count < 32; count++ {

		if count % 2 == 0 {
			for idx := 0; idx < 64; idx++ {
				if idx % 2 == 0 {
					localBuf[count * 64 + idx] = byte(100)
				}
			}
		}

		localBuf[0] = byte(100 + 10 * (count % 3))
		localBuf[1] = byte(100 + 20 * (count % 5))
		localBuf[2] = byte(100 + 30 * (count % 7))
		
		bufChan <- localBuf
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("Done generating")
}

func main() {
	fmt.Println("RAYGRAPHICS")

	bufChan := make(chan []byte)
	lastBuf := make([]byte, 64 * 32)
	render := raygraphics.NewRaylibRender(64, 32, bufChan)

	rl.InitWindow(1024, 512, "RayRender tests")
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
