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

	PropagateLayout()
	Layout()
	NeedsLayout() bool
	PrefSizeChanged() bool

	PropagateUpdate()
	Update()

	PropagateDraw(Painter)
	Draw(Painter)
	NeedsRedraw() bool
}
