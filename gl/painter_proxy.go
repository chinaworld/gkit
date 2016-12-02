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
	r = normalizeCoords(r, p.frame.Size)
	p.impl.drawRect(gkit.Rect{
		gkit.Point{p.frame.X + r.X, p.frame.Y + r.Y},
		r.Size,
	}, z+1)
}

func (p *painterProxy) SubPainter(r gkit.Rect) gkit.Painter {
	r = normalizeCoords(r, p.frame.Size)
	return &painterProxy{
		impl:  p,
		frame: r,
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

func (p *painterProxy) DrawText(o gkit.Point, text string) {
	p.drawText(o, 0, text)
}
func (p *painterProxy) drawText(o gkit.Point, z uint32, text string) {
	r := normalizeCoords(
		gkit.Rect{Point: o}, p.frame.Size)
	p.impl.drawText(gkit.Point{r.X + p.frame.X, r.Y + p.frame.Y}, z+1, text)
}

func (p *painterProxy) DrawImage(r gkit.Rect, image image.Image) {
	p.drawImage(r, 0, image)
}

func (p *painterProxy) drawImage(r gkit.Rect, z uint32, image image.Image) {
	r = normalizeCoords(r, p.frame.Size)
	p.impl.drawImage(gkit.Rect{
		gkit.Point{r.X + p.frame.X, r.Y + p.frame.Y},
		r.Size,
	}, z+1, image)
}
