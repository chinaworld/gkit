package gkit

type View struct {
	Background Color

	Frame Rect

	children []*View

	LayoutStrategy LayoutStrategy
	LayoutSettings LayoutSettings
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

func (v *View) DoLayout() {
	newViewLayout(v).Layout(v.Frame)
}

func newViewLayout(v *View) SubLayouter {
	strategy := v.LayoutStrategy
	if strategy == nil {
		strategy = NoneStrategy
	}
	children := make([]SubLayouter, 0, len(v.children))
	for _, child := range v.children {
		children = append(children, newViewLayout(child))
	}
	return &viewLayout{
		Strategy:       strategy,
		settings:       v.LayoutSettings,
		children:       children,
		heightForWidth: make(map[uint32]uint32),
		frame:          &v.Frame,
	}
}

type viewLayout struct {
	Strategy LayoutStrategy
	settings LayoutSettings
	children []SubLayouter

	preferredSize  *Size
	heightForWidth map[uint32]uint32

	frame *Rect
}

var _ SubLayouter = &viewLayout{}

func (l *viewLayout) Settings() LayoutSettings {
	return l.settings
}

func (l *viewLayout) PreferedSize() Size {
	if l.preferredSize == nil {
		size := l.Strategy.PreferedSize(l.settings, l.children...)
		l.preferredSize = &size
	}
	return *l.preferredSize
}

func (l *viewLayout) HeightForWidth(width uint32) uint32 {
	if _, ok := l.heightForWidth[width]; !ok {
		l.heightForWidth[width] = l.Strategy.HeightForWidth(l.settings, width, l.children...)
	}
	return l.heightForWidth[width]
}

func (l *viewLayout) Layout(frame Rect) {
	*l.frame = frame
	l.Strategy.Layout(l.settings, Size{
		Width:  frame.Width,
		Height: frame.Height,
	}, l.children...)
}
