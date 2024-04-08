package main

import (
	"fmt"
	"github.com/atomic-7/emu8tor/internal/raylibgraphics"
	"time"
)

func main() {
	rayrender := raygraphics.CreateRayApp(64, 32, 1024, 512)

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
	fmt.Println("Done")

}
