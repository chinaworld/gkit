package gkit

import (
	"image"
)

type DrawingContext interface {
	BeginPaint() Painter
	EndPaint(painter Painter)
}

type Painter interface {
	SubPainter(x, y, width, height uint32) Painter
	SetColor(c Color)
	DrawRect(x, y, width, height uint32)
	SetFont(f *Font)
	SetFontSize(size uint32)
	DrawText(x, y uint32, text string)
	DrawImage(x, y, width, height uint32, image image.Image)
}
