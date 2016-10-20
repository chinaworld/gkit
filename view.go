package gkit

type View struct {
	Background Color

	Width  uint32
	Height uint32
}

func (v *View) Draw(p Painter) {
	p.SetColor(v.Background)
	p.DrawRect(0, 0, v.Width, v.Height)
}
