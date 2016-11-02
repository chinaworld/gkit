package gl

import (
	"github.com/alex-ac/gkit"
)

type glPainterInternal interface {
	setColor(c gkit.Color)
	drawRect(x, y, z, width, height uint32)
	setFont(font *gkit.Font)
	setFontSize(size uint32)
	drawText(x, y, z uint32, text string)
}

func normalizeCoords(x, y, width, height, maxWidth, maxHeight uint32) (uint32, uint32, uint32, uint32) {
	if x > maxWidth {
		x = maxWidth
	}
	if y > maxHeight {
		y = maxHeight
	}
	if x+width > maxWidth {
		width = maxWidth - x
	}
	if y+height > maxHeight {
		height = maxHeight - y
	}

	return x, y, width, height
}

type painter struct {
	context *drawingContext
	width   uint32
	height  uint32

	mask *bitmap8bpp

	vertices     []float32
	currentColor [4]float32

	currentFont     *gkit.Font
	currentFontSize uint32
}

var _ gkit.Painter = &painter{}
var _ glPainterInternal = &painter{}

func (p *painter) SubPainter(x, y, width, height uint32) gkit.Painter {
	x, y, width, height = normalizeCoords(
		x, y, width, height, p.width, p.height)
	return &painterProxy{
		impl:   p,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (p *painter) SetColor(c gkit.Color) {
	p.setColor(c)
}

func (p *painter) setColor(c gkit.Color) {
	p.currentColor = [4]float32{
		float32(c.R()) / 255,
		float32(c.G()) / 255,
		float32(c.B()) / 255,
		float32(c.A()) / 255,
	}
}

func (p *painter) DrawRect(x, y, width, height uint32) {
	p.drawRect(x, y, 0, width, height)
}

func (p *painter) drawRect(x, y, z, width, height uint32) {
	left, top, right, bottom, Z := float32(x), float32(y), float32(x+width), float32(y+height), float32(z)
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	U, V := float32(0), float32(0)
	p.vertices = append(p.vertices,
		left, top, Z, R, G, B, A, U, V,
		right, top, Z, R, G, B, A, U, V,
		left, bottom, Z, R, G, B, A, U, V,
		right, top, Z, R, G, B, A, U, V,
		right, bottom, Z, R, G, B, A, U, V,
		left, bottom, Z, R, G, B, A, U, V,
	)
}

func (p *painter) SetFont(font *gkit.Font) {
	p.setFont(font)
}

func (p *painter) setFont(font *gkit.Font) {
	p.currentFont = font
}

func (p *painter) SetFontSize(size uint32) {
	p.setFontSize(size)
}

func (p *painter) setFontSize(size uint32) {
	p.currentFontSize = size
}

func (p *painter) DrawText(x, y uint32, text string) {
	p.drawText(x, y, 0, text)
}

func (p *painter) drawText(x, y, z uint32, text string) {
	size := p.currentFont.StringSize(p.currentFontSize, text)
	p.currentFont.DrawString(p.currentFontSize, text, x, y, p.mask)

	left, top, right, bottom, Z := float32(x), float32(y), float32(x+size.Width), float32(x+size.Height), float32(z)
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	p.vertices = append(p.vertices,
		left, top, Z, R, G, B, A, left, top,
		right, top, Z, R, G, B, A, right, top,
		left, bottom, Z, R, G, B, A, left, bottom,
		right, top, Z, R, G, B, A, right, top,
		right, bottom, Z, R, G, B, A, right, bottom,
		left, bottom, Z, R, G, B, A, left, bottom,
	)
}