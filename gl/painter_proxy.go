package gl

import (
	"image"

	"github.com/alex-ac/gkit"
)

type painterProxy struct {
	impl glPainterInternal

	frame gkit.Rect
}

var _ gkit.Painter = &painterProxy{}
var _ glPainterInternal = &painterProxy{}

func (p *painterProxy) DrawRect(r gkit.Rect) {
	p.drawRect(r, 0)
}

func (p *painterProxy) drawRect(r gkit.Rect, z uint32) {
	x, y, width, height := normalizeCoords(
		r.X, r.Y, r.Width, r.Height, p.frame.Width, p.frame.Height)
	p.impl.drawRect(gkit.Rect{
		gkit.Point{p.frame.X + x, p.frame.Y + y},
		gkit.Size{width, height},
	}, z+1)
}

func (p *painterProxy) SubPainter(r gkit.Rect) gkit.Painter {
	x, y, width, height := normalizeCoords(
		r.X, r.Y, r.Width, r.Height, p.frame.Width, p.frame.Height)
	return &painterProxy{
		impl: p,
		frame: gkit.Rect{
			gkit.Point{x, y},
			gkit.Size{width, height},
		},
	}
}

func (p *painterProxy) SetColor(c gkit.Color) {
	p.setColor(c)
}

func (p *painterProxy) setColor(c gkit.Color) {
	p.impl.setColor(c)
}

func (p *painterProxy) SetFont(f *gkit.Font) {
	p.setFont(f)
}
func (p *painterProxy) setFont(f *gkit.Font) {
	p.impl.setFont(f)
}

func (p *painterProxy) SetFontSize(size uint32) {
	p.setFontSize(size)
}

func (p *painterProxy) setFontSize(size uint32) {
	p.impl.setFontSize(size)
}

func (p *painterProxy) DrawText(x, y uint32, text string) {
	p.drawText(x, y, 0, text)
}
func (p *painterProxy) drawText(x, y, z uint32, text string) {
	x, y, _, _ = normalizeCoords(
		x, y, 0, 0, p.frame.Width, p.frame.Height)
	p.impl.drawText(x+p.frame.X, y+p.frame.Y, z+1, text)
}

func (p *painterProxy) DrawImage(r gkit.Rect, image image.Image) {
	p.drawImage(r, 0, image)
}

func (p *painterProxy) drawImage(r gkit.Rect, z uint32, image image.Image) {
	x, y, width, height := normalizeCoords(
		r.X, r.Y, r.Width, r.Height, p.frame.Width, p.frame.Height)
	p.impl.drawImage(gkit.Rect{
		gkit.Point{x + p.frame.X, y + p.frame.Y},
		gkit.Size{width, height},
	}, z+1, image)
}
