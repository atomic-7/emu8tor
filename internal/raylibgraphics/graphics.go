package raygraphics

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Implements Drawable

type RayRender struct {
	Width int
	Height int
}

func (rr *RayRender) SetWidth(w int) {
	rr.Width = w
}

func (rr *RayRender) SetHeight(h int) {
	rr.Height = h
}

func (rr *RayRender) Draw(dbuf []byte) {
	/*
		GetScreenWidth()
		void EnableEventWaiting(void)
		float GetFrameTime(void)
		SetShapesTexture(Texture2d texture, Rectangle source)
	*/
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Raylib Woohoo!", 190, 200, 20, rl.LightGray)
	rl.EndDrawing()

}

func NewRaylibRender(w int, h int) *RayRender {

	var rr RayRender
	rr.Width = w
	rr.Height = h

	return &rr
}
