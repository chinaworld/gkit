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
	drawText(o gkit.Point, z uint32, text string)
	drawImage(r gkit.Rect, z uint32, image image.Image)
	enableRedraw()
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func normalizeCoords(r gkit.Rect, max gkit.Size) gkit.Rect {
	r.X = min(r.X, max.Width)
	r.Y = min(r.Y, max.Height)
	r.Width = min(r.Width, max.Width-r.X)
	r.Height = min(r.Height, max.Height-r.Y)
	return r
}

func getFloats(r gkit.Rect, z uint32) (float32, float32, float32, float32, float32) {
	left, top := r.Point.AsFloats()
	right, bottom := r.RightBottom().AsFloats()
	Z := float32(z)
	return left, top, right, bottom, Z
}

type instruction func(*painter)

type painter struct {
	context *drawingContext
	size    gkit.Size

	mask   *image.Gray
	images []*image.RGBA

	vertices     []float32
	currentColor [4]float32

	currentFont     *gkit.Font
	currentFontSize uint32

	scaleFactor float32
	doRedraw    bool

	instructions []instruction
}

var _ gkit.Painter = &painter{}
var _ glPainterInternal = &painter{}

func (p *painter) DrawLayer(r gkit.Rect, l gkit.Layer) {
	r = normalizeCoords(r, p.size)
	painter := &painterProxy{
		impl:  p,
		frame: r,
	}

	if l.NeedsRedraw() {
		p.enableRedraw()
	}
	l.Draw(painter)
	l.PropagateDraw(painter)
	if p.doRedraw {
		for _, i := range p.instructions {
			i(p)
		}
		p.instructions = p.instructions[0:0]
	}
}

func (p *painter) enableRedraw() {
	p.doRedraw = true
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
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	p.addInstruction(func(p *painter) {
		left, top, right, bottom, Z := getFloats(r, z)
		U, V, W := float32(0), float32(0), float32(-1)
		p.vertices = append(p.vertices,
			left, top, Z, R, G, B, A, U, V, U, V, W,
			right, top, Z, R, G, B, A, U, V, U, V, W,
			left, bottom, Z, R, G, B, A, U, V, U, V, W,
			right, top, Z, R, G, B, A, U, V, U, V, W,
			right, bottom, Z, R, G, B, A, U, V, U, V, W,
			left, bottom, Z, R, G, B, A, U, V, U, V, W,
		)
	})
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

func (p *painter) DrawText(o gkit.Point, text string) {
	p.drawText(o, 0, text)
}

func (p *painter) drawText(o gkit.Point, z uint32, text string) {
	font := p.currentFont
	fontSize := p.currentFontSize
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	p.addInstruction(func(p *painter) {
		size := font.StringSize(fontSize, text)
		font.DrawString(uint32(float32(fontSize)*p.scaleFactor), text, o.Scale(p.scaleFactor), p.mask)
		r := gkit.Rect{o, size}

		left, top, right, bottom, Z := getFloats(r, z)
		U, V, W := float32(0), float32(0), float32(-1)
		p.vertices = append(p.vertices,
			left, top, Z, R, G, B, A, left, top, U, V, W,
			right, top, Z, R, G, B, A, right, top, U, V, W,
			left, bottom, Z, R, G, B, A, left, bottom, U, V, W,
			right, top, Z, R, G, B, A, right, top, U, V, W,
			right, bottom, Z, R, G, B, A, right, bottom, U, V, W,
			left, bottom, Z, R, G, B, A, left, bottom, U, V, W,
		)
	})
}

func (p *painter) DrawImage(r gkit.Rect, img image.Image) {
	p.drawImage(r, 0, img)
}

func (p *painter) drawImage(r gkit.Rect, z uint32, img image.Image) {
	R, G, B, A := p.currentColor[0], p.currentColor[1], p.currentColor[2], p.currentColor[3]
	p.addInstruction(func(p *painter) {
		size := textureSideSize(p.size)
		size = uint32(float32(size) * p.scaleFactor)
		imageCopy := image.NewRGBA(image.Rectangle{
			Max: image.Point{int(size), int(size)},
		})
		bounds := img.Bounds()
		draw.Copy(imageCopy, image.Point{}, img, bounds, draw.Over, nil)
		left, top, right, bottom, Z := getFloats(r, z)
		W := float32(len(p.images))
		p.images = append(p.images, imageCopy)
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
	})
}

func (p *painter) addInstruction(i instruction) {
	p.instructions = append(p.instructions, i)
}
