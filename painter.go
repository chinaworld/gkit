package gkit

import (
	"image"
)

type DrawingContext interface {
	BeginPaint() Painter
	EndPaint(painter Painter)
}

type Layer interface {
	PropagateDraw(Painter)
	Draw(Painter)
	NeedsRedraw() bool
}

type Painter interface {
	DrawLayer(r Rect, l Layer)
	SetColor(c Color)
	DrawRect(r Rect)
	SetFont(f *Font)
	SetFontSize(size uint32)
	DrawText(p Point, text string)
	DrawImage(r Rect, image image.Image)
}
