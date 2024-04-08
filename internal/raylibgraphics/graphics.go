package raygraphics

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Implements Drawable

type RayRender struct {
	Width   int
	Height  int
	rayChan chan []byte
}

func (rr *RayRender) SetWidth(w int) {
	rr.Width = w
}

func (rr *RayRender) SetHeight(h int) {
	rr.Height = h
}

func (rr *RayRender) Draw(dbuf []byte) {
	rr.rayChan <- dbuf
}
func (rr *RayRender) DrawBuf(dbuf []byte) {
	//color := rl.NewColor(dbuf[0], dbuf[1], dbuf[2], 255)
	pxlLen := int32(rl.GetScreenWidth() / rr.Height)

	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	var x,y int32
	for y = 0; y < int32(rr.Height); y++ {
		for x = 0; x < int32(rr.Width); x++ {
			if dbuf[y * int32(rr.Width) + x] > 0 {
				rl.DrawRectangle(x*pxlLen, y*pxlLen, pxlLen, pxlLen, rl.Black)
			}
		}
	}
	rl.EndDrawing()
}

func NewRaylibRender(w int, h int, bufChan chan []byte) *RayRender {

	var rr RayRender
	rr.Width = w
	rr.Height = h
	rr.rayChan = bufChan

	return &rr
}

// Executing the renderloop in a seperate goroutine has led to crashes
// presumably because of conflicts with the garbage collector not restoring the stack as expected
func RenderLoop(dWidth int, dHeight int, wWidth int, wHeight int, render *RayRender) {
	// abort context here
	lastBuf :=make([]byte, dWidth*dHeight)
	rl.InitWindow(int32(wWidth), int32(wHeight), "Emu8tor - Raylib Render")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		select {
		case buf := <- render.rayChan:
			copy(lastBuf, buf)	
			render.DrawBuf(lastBuf)
		default:
			render.DrawBuf(lastBuf)
		}
	}
	fmt.Println("raylib: should close")

}
/*
	void EnableEventWaiting(void)
	float GetFrameTime(void)
	SetShapesTexture(Texture2d texture, Rectangle source)
*/
