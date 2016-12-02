package gkit

type View interface {
	AddChild(View)
	DeleteChild(View)
	Children() []View

	MinSize() Size
	MaxSize() Size
	PrefSize() Size
	SetSize(Size)
	SetOrigin(Point)
	SetFrame(Rect)
	Frame() Rect
	Origin() Point
	Size() Size
	UpdateSizes()
	SetBorders(v SideValues)
	Borders() SideValues
	Bounds() Rect

	PropagateLayout()
	Layout()
	NeedsLayout() bool
	PrefSizeChanged() bool

	PropagateUpdate()
	Update()

	Layer
}
