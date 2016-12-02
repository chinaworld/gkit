package gkit

import (
	"image"
)

type DrawingContext interface {
	BeginPaint() Painter
	EndPaint(painter Painter)
}

type Painter interface {
	SubPainter(r Rect) Painter
	SetColor(c Color)
	DrawRect(r Rect)
	SetFont(f *Font)
	SetFontSize(size uint32)
	DrawText(x, y uint32, text string)
	DrawImage(r Rect, image image.Image)
}
