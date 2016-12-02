package gl

import (
	"image"

	"golang.org/x/image/draw"

	"github.com/alex-ac/gkit"
)

type glPainterInternal interface {
	setColor(c gkit.Color)
	drawRect(r gkit.Rect, z uint32)
	setFont(font *gkit.Font)
	setFontSize(size uint32)
	drawText(x, y, z uint32, text string)
	drawImage(r gkit.Rect, z uint32, image image.Image)
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func normalizeCoords(r gkit.Rect, maxWidth, maxHeight uint32) gkit.Rect {
	r.X = min(r.X, maxWidth)
	r.Y = min(r.Y, maxHeight)
	r.Width = min(r.Width, maxWidth-r.X)
	r.Height = min(r.Height, maxHeight-r.Y)
	return r
}

type painter struct {
	context *drawingContext
	width   uint32
	height  uint32

	mask   *image.Gray
	images []*image.RGBA

	vertices     []float32
	currentColor [4]float32

	currentFont     *gkit.Font
	currentFontSize uint32
}

var _ gkit.Painter = &painter{}
var _ glPainterInternal = &painter{}

func (p *painter) SubPainter(r gkit.Rect) gkit.Painter {
	r = normalizeCoords(r, p.width, p.height)
	return &painterProxy{
		impl:  p,
		frame: r,
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

func (p *painter) DrawRect(r gkit.Rect) {
	p.drawRect(r, 0)
}

func (p *painter) drawRect(r gkit.Rect, z uint32) {
	left, top, right, bottom, Z := float32(r.X), float32(r.Y), float32(r.X+r.Width), float32(r.Y+r.Height), float32(z)
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	U, V, W := float32(0), float32(0), float32(-1)
	p.vertices = append(p.vertices,
		left, top, Z, R, G, B, A, U, V, U, V, W,
		right, top, Z, R, G, B, A, U, V, U, V, W,
		left, bottom, Z, R, G, B, A, U, V, U, V, W,
		right, top, Z, R, G, B, A, U, V, U, V, W,
		right, bottom, Z, R, G, B, A, U, V, U, V, W,
		left, bottom, Z, R, G, B, A, U, V, U, V, W,
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

	left, top, right, bottom, Z := float32(x), float32(y), float32(x+size.Width), float32(y+size.Height), float32(z)
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	U, V, W := float32(0), float32(0), float32(-1)
	p.vertices = append(p.vertices,
		left, top, Z, R, G, B, A, left, top, U, V, W,
		right, top, Z, R, G, B, A, right, top, U, V, W,
		left, bottom, Z, R, G, B, A, left, bottom, U, V, W,
		right, top, Z, R, G, B, A, right, top, U, V, W,
		right, bottom, Z, R, G, B, A, right, bottom, U, V, W,
		left, bottom, Z, R, G, B, A, left, bottom, U, V, W,
	)
}

func (p *painter) DrawImage(r gkit.Rect, img image.Image) {
	p.drawImage(r, 0, img)
}

func (p *painter) drawImage(r gkit.Rect, z uint32, img image.Image) {
	size := textureSideSize(p.width, p.height)
	imageCopy := image.NewRGBA(image.Rectangle{
		Max: image.Point{int(size), int(size)},
	})
	bounds := img.Bounds()
	draw.Copy(imageCopy, image.Point{}, img, bounds, draw.Over, nil)
	W := float32(len(p.images))
	p.images = append(p.images, imageCopy)
	left, top, right, bottom, Z := float32(r.X), float32(r.Y), float32(r.X+r.Width), float32(r.Y+r.Height), float32(z)
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	U, V := float32(0), float32(0)
	imageWidth, imageHeight := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	uvLeft, uvTop, uvRight, uvBottom := float32(0), float32(0), float32(imageWidth)/float32(size), float32(imageHeight)/float32(size)
	p.vertices = append(p.vertices,
		left, top, Z, R, G, B, A, U, V, uvLeft, uvTop, W,
		right, top, Z, R, G, B, A, U, V, uvRight, uvTop, W,
		left, bottom, Z, R, G, B, A, U, V, uvLeft, uvBottom, W,
		right, top, Z, R, G, B, A, U, V, uvRight, uvTop, W,
		right, bottom, Z, R, G, B, A, U, V, uvRight, uvBottom, W,
		left, bottom, Z, R, G, B, A, U, V, uvLeft, uvBottom, W,
	)
}
