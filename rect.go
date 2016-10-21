package gkit

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
	X      uint32
	Y      uint32
	Width  uint32
	Height uint32
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
