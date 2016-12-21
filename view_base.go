package gkit

type ViewBase struct {
	View View

	frame   Rect
	borders SideValues

	children []View

	needsLayout     bool
	needsRedraw     bool
	prefSizeChanged bool

	minSize  Size
	prefSize Size
	maxSize  Size
}

func (v *ViewBase) AddChild(view View) {
	for _, child := range v.children {
		if child == view {
			return
		}
	}
	if v.children == nil {
		v.children = make([]View, 0, 10)
	}
	v.children = append(v.children, view)
	v.needsLayout = true
	v.prefSizeChanged = true
}

func (v *ViewBase) DeleteChild(view View) {
	for i, child := range v.children {
		if child == view {
			v.children = append(v.children[:i], v.children[i+1:]...)
			v.needsLayout = true
			v.prefSizeChanged = true
			break
		}
	}
}

func (v *ViewBase) PropagateUpdate() {
	v.View.Update()
	needsUpdateSize := false
	for _, child := range v.children {
		child.PropagateUpdate()
		needsUpdateSize = needsUpdateSize || child.PrefSizeChanged()
	}
	if v.needsLayout || v.prefSizeChanged || needsUpdateSize {
		v.View.UpdateSizes()
	}
}

func (v *ViewBase) Children() []View {
	return v.children
}

func (v *ViewBase) Frame() Rect {
	return v.frame
}

func (v *ViewBase) PropagateDraw(p Painter) {
	for _, child := range v.children {
		frame := child.Frame()
		p.DrawLayer(frame, child)
	}
	v.needsRedraw = false
}

func (v *ViewBase) PropagateLayout() {
	if v.needsLayout || v.prefSizeChanged {
		v.View.Layout()
		v.needsLayout = false
		v.prefSizeChanged = true
	}
	for _, child := range v.children {
		child.PropagateLayout()
		v.needsRedraw = v.needsRedraw || child.NeedsRedraw()
	}
}

func (v *ViewBase) SetOrigin(p Point) {
	v.frame.X, v.frame.Y = p.X, p.Y
}

func (v *ViewBase) Origin() Point {
	return Point{v.frame.X, v.frame.Y}
}

func (v *ViewBase) SetSize(s Size) {
	if v.frame.Width != s.Width || v.frame.Height != s.Height {
		v.needsLayout = true
		v.needsRedraw = true
	}
	v.frame.Width, v.frame.Height = s.Width, s.Height
}

func (v *ViewBase) Size() Size {
	return Size{v.frame.Width, v.frame.Height}
}

func (v *ViewBase) SetFrame(frame Rect) {
	v.SetOrigin(frame.Point)
	v.SetSize(frame.Size)
}

func (v *ViewBase) SetMinSize(size Size) {
	v.minSize = size
}

func (v *ViewBase) MinSize() Size {
	return v.minSize
}

func (v *ViewBase) SetMaxSize(size Size) {
	v.maxSize = size
}

func (v *ViewBase) MaxSize() Size {
	return v.maxSize
}

func (v *ViewBase) SetPrefSize(size Size) {
	size = size.Outset(v.borders)
	if v.prefSize != size {
		v.prefSize = size
		v.prefSizeChanged = true
	}
}

func (v *ViewBase) SetPrefSizeChanged() {
	v.prefSizeChanged = true
}

func (v *ViewBase) PrefSize() Size {
	return v.prefSize
}

func (v *ViewBase) NeedsLayout() bool {
	return v.needsLayout
}

func (v *ViewBase) SetNeedsLayout() {
	v.needsLayout = true
}

func (v *ViewBase) NeedsRedraw() bool {
	return v.needsRedraw
}

func (v *ViewBase) SetNeedsRedraw() {
	v.needsRedraw = true
}

func (v *ViewBase) PrefSizeChanged() bool {
	return v.prefSizeChanged
}

func (v *ViewBase) SetBorders(borders SideValues) {
	v.borders = borders
}

func (v *ViewBase) Borders() SideValues {
	return v.borders
}

func (v *ViewBase) Bounds() Rect {
	return Rect{Size: v.Size()}.Inset(v.Borders())
}
