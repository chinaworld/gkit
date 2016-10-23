package gkit

type Layouter interface {
	MinSize() Size
	PrefSize() Size
	MaxSize() Size

	HeightForWidth(uint32) uint32

	Layout(Rect)
}

type View struct {
	Background Color

	Frame Rect

	children []*View

	MinSize        Size
	MaxSize        Size
	PrefSize       Size
	HeightForWidth func(*View, uint32) uint32
	Layout         func(*View, Rect, []Layouter)
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
	p.DrawRect(0, 0, v.Frame.Width, v.Frame.Height)

	for _, child := range v.children {
		child.Draw(p.SubPainter(child.Frame.X, child.Frame.Y, child.Frame.Width, child.Frame.Height))
	}
}

func (v *View) Layouter() Layouter {
	return viewLayouter{v}
}

type viewLayouter struct{ *View }

var _ Layouter = viewLayouter{}

func (v viewLayouter) Layout(frame Rect) {
	v.Frame = frame
	if v.View.Layout != nil {
		layouters := make([]Layouter, 0, len(v.View.children))
		for _, child := range v.View.children {
			layouters = append(layouters, child.Layouter())
		}
		v.View.Layout(v.View, Rect{0, 0, frame.Width, frame.Height}, layouters)
	}
}

func (v viewLayouter) MaxSize() Size {
	return v.View.MaxSize
}

func (v viewLayouter) MinSize() Size {
	return v.View.MinSize
}

func (v viewLayouter) PrefSize() Size {
	return v.View.PrefSize
}

func (v viewLayouter) HeightForWidth(width uint32) uint32 {
	if v.View.HeightForWidth != nil {
		return v.View.HeightForWidth(v.View, width)
	}
	return 0
}
