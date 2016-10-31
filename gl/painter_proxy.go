package gl

import (
	"github.com/alex-ac/gkit"
)

type painterProxy struct {
	impl glPainterInternal

	x, y, width, height uint32
}

var _ gkit.Painter = &painterProxy{}
var _ glPainterInternal = &painterProxy{}

func (p *painterProxy) DrawRect(x, y, width, height uint32) {
	p.drawRect(x, y, 0, width, height)
}

func (p *painterProxy) drawRect(x, y, z, width, height uint32) {
	x, y, width, height = normalizeCoords(
		x, y, width, height, p.width, p.height)
	p.impl.drawRect(p.x+x, p.y+y, z+1, width, height)
}

func (p *painterProxy) SubPainter(x, y, width, height uint32) gkit.Painter {
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
	p.impl.drawText(x+p.x, y+p.y, z+1, text)
}
