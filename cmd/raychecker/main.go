package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
	"time"
)

func patternGenerator(rayrender *raygraphics.RayRender) {
	
	pattern1 := make([]byte, 64 * 32)
	pattern2 := make([]byte, 64 * 32)

	for idx := 0; idx < 64 * 32; idx++ {
		if idx % 2 == 0 {
			pattern1[idx] = 0xf
			pattern2[idx] = 0
		} else {
			pattern1[idx] = 0
			pattern2[idx] = 0xf
		}
	}
	
	for count := 0; count < 200; count ++{
		if count % 2 == 0 {
			rayrender.Draw(pattern1)
		} else {
			rayrender.Draw(pattern2)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	//rayrender := raygraphics.CreateRayApp(64, 32, 1024, 512)
	bufChan := make(chan []byte)
	rayrender := raygraphics.NewRaylibRender(64, 32, bufChan)
	go patternGenerator(rayrender)
	raygraphics.RenderLoop(64, 32, 1024, 512, rayrender, nil)

	fmt.Println("Done")
}
