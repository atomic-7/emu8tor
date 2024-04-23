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
	pxlLen := int32(rl.GetScreenWidth() / rr.Width)

	rl.ClearBackground(rl.Black)
	drawGrid(pxlLen)
	var x, y int32
	for y = 0; y < int32(rr.Height); y++ {
		for x = 0; x < int32(rr.Width); x++ {
			if dbuf[y*int32(rr.Width)+x] > 0 {
				rl.DrawRectangle(x*pxlLen, y*pxlLen, pxlLen, pxlLen, rl.RayWhite)
			}
		}
	}
}

func drawGrid(pxlLen int32) {

	var x, y, vlines, hlines int32

	w := int32(rl.GetScreenWidth())
	h := int32(rl.GetScreenHeight())
	hlines = h / pxlLen
	vlines = w / pxlLen

	for x = 0; x < hlines; x++ {
		rl.DrawLine(0, x*pxlLen, w, x*pxlLen, rl.Lime)
	}
	for y = 0; y < vlines; y++ {
		rl.DrawLine(y*pxlLen, 0, y*pxlLen, h, rl.Lime)
	}
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
func RenderLoop(dWidth int, dHeight int, wWidth int, wHeight int, render *RayRender, step chan bool) {
	// abort context here
	lastBuf := make([]byte, dWidth*dHeight)
	rl.InitWindow(int32(wWidth), int32(wHeight), "Emu8tor - Raylib Render")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		select {
		case buf := <-render.rayChan:
			copy(lastBuf, buf)
			rl.BeginDrawing()
			render.DrawBuf(lastBuf)
			rl.EndDrawing()
		default:
			rl.BeginDrawing()
			render.DrawBuf(lastBuf)
			rl.EndDrawing()
		}
		if rl.IsKeyPressed(rl.KeyS) {
			if step != nil {
				step <- true
			}
		}
	}
	fmt.Println("raylib: should close")

}

/*
	void EnableEventWaiting(void)
	float GetFrameTime(void)
	SetShapesTexture(Texture2d texture, Rectangle source)
*/
