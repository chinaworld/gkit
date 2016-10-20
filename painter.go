package gkit

type Painter interface {
	SubPainter(x, y, width, height uint32) Painter
	SetColor(c Color)
	DrawRect(x, y, width, height uint32)
}
