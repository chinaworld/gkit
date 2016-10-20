package gkit

type View struct {
	Background Color

	X      uint32
	Y      uint32
	Width  uint32
	Height uint32

	children []*View
}

func (v *View) AddChild(view *View) {
	for _, child := range v.children {
		if child == view {
			return
		}
	}
	if v.children == nil {
		v.children = make([]*View, 0, 10)
	}
	v.children = append(v.children, view)
}

func (v *View) DeleteChild(view *View) {
	for i, child := range v.children {
		if child == view {
			v.children = append(v.children[:i], v.children[i+1:]...)
			break
		}
	}
}

func (v *View) Draw(p Painter) {
	p.SetColor(v.Background)
	p.DrawRect(0, 0, v.Width, v.Height)

	for _, child := range v.children {
		child.Draw(p.SubPainter(child.X, child.Y, child.Width, child.Height))
	}
}
