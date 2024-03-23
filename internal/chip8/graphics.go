package chip8

import (
	"bytes"
	"strings"
	// "encoding/hex"
	"fmt"
)

type Drawable interface {
	SetWidth(int)
	SetHeight(int)
	Draw([]byte)
}

type DebugRender struct {
	Width  int
	Height int
}

func (db *DebugRender) SetWidth(w int) {
	db.Width = w
}

func (db *DebugRender) SetHeight(h int) {
	db.Height = h
}

func drawLine(line []byte) string {
	return string(
		bytes.Map(func(r rune) rune {
			if r > 0 {
				return '#'
			} else {
				return '.'
			}
		}, line))
}

func (db *DebugRender) Draw(dbuf []byte) {

	// fmt.Println(hex.Dump(dbuf))
	lineSeperator := strings.Repeat("-", db.Width)
	println(lineSeperator)
	for y := 0; y < db.Height; y++ {
		for x := 0; x < db.Width; x++ {
			//[64/896]0xc000025480
			//println(dbuf[y * db.Width:(y+1)*db.Width])
			//fmt.Printf("%x\n", dbuf[y*db.Width:(y+1)*db.Width])
			fmt.Printf("%s\n", drawLine(dbuf[y*db.Width:(y+1)*db.Width]))
		}
	}
	println(lineSeperator)
}

func NewDebugRender(w, h int) *DebugRender {
	var dbr DebugRender
	dbr.Width = w
	dbr.Height = h
	return &dbr
}
