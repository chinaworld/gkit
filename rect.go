package gkit

type Point struct {
	X uint32
	Y uint32
}

func (p Point) Offset(p2 Point) Point {
	p.X += p2.X
	p.Y += p2.Y
	return p
}

func (p Point) AsFloats() (float32, float32) {
	return float32(p.X), float32(p.Y)
}

type Size struct {
	Width  uint32
	Height uint32
}

func max(values ...uint32) uint32 {
	var max uint32 = 0
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func (s Size) Outset(insets SideValues) Size {
	s.Width += insets.Left + insets.Right
	s.Height += insets.Top + insets.Bottom
	return s
}

type Rect struct {
	Point
	Size
}

type SideValues struct {
	Left   uint32
	Right  uint32
	Bottom uint32
	Top    uint32
}

func (r Rect) Inset(insets SideValues) Rect {
	if insets.Left+insets.Right > r.Width {
		insets.Left = r.Width / 2
		insets.Right = r.Width / 2
	}
	if insets.Top+insets.Bottom > r.Height {
		insets.Top = r.Height / 2
		insets.Bottom = r.Height / 2
	}
	r.X += insets.Left
	r.Y += insets.Top
	r.Width -= insets.Left + insets.Right
	r.Height -= insets.Top + insets.Bottom
	return r
}

func (r Rect) Offset(p Point) Rect {
	r.Point = r.Point.Offset(p)
	return r
}

func (r Rect) RightBottom() Point {
	return Point{r.X + r.Width, r.Y + r.Height}
}
