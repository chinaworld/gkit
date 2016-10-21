package gkit

type SubLayouter interface {
	Settings() LayoutSettings
	PreferedSize() Size
	HeightForWidth(width uint32) uint32
	Layout(Rect)
}

type LayoutStrategy interface {
	PreferedSize(s LayoutSettings, subLayouters ...SubLayouter) Size
	HeightForWidth(s LayoutSettings, width uint32, subLayouters ...SubLayouter) uint32
	Layout(s LayoutSettings, size Size, subLayouters ...SubLayouter)
}

type FlexJustifyType int

const (
	FlexStart FlexJustifyType = iota
	FlexEnd
	FlexCenter
	SpaceAround
	SpaceBetween
	Stretch
)

type LayoutSettings struct {
	FlexJustify FlexJustifyType
	FlexAlign   FlexJustifyType

	FlexBase   uint32
	FlexGrow   float32
	FlexShrink float32

	FlexGap uint32

	Padding SideValues
}

type noneStrategy struct{}

var NoneStrategy LayoutStrategy = noneStrategy{}

func (noneStrategy) PreferedSize(LayoutSettings, ...SubLayouter) Size {
	return Size{}
}

func (noneStrategy) HeightForWidth(LayoutSettings, uint32, ...SubLayouter) uint32 {
	return 0
}

func (noneStrategy) Layout(LayoutSettings, Size, ...SubLayouter) {}

type flexRow struct{}

var FlexRow LayoutStrategy = flexRow{}

func (flexRow) PreferedSize(s LayoutSettings, subLayouters ...SubLayouter) Size {
	size := Size{}

	for i, sublayout := range subLayouters {
		if i > 0 {
			size.Width += s.FlexGap
		}
		sublayoutSize := sublayout.PreferedSize()
		size.Width += sublayoutSize.Width
		size.Height = max(size.Height, sublayoutSize.Height)
	}

	size = size.Outset(s.Padding)

	if s.FlexBase > size.Width {
		size.Width = s.FlexBase
	}

	return size
}

func (flexRow) HeightForWidth(s LayoutSettings, width uint32, subLayouters ...SubLayouter) uint32 {
	size := Size{}

	settings := make([]LayoutSettings, 0, len(subLayouters))
	sizes := make([]Size, 0, len(subLayouters))

	for i, sublayout := range subLayouters {
		if i > 0 {
			size.Width += s.FlexGap
		}
		sublayoutSize := sublayout.PreferedSize()

		settings = append(settings, sublayout.Settings())
		sizes = append(sizes, sublayoutSize)

		size.Width += sublayoutSize.Width
		size.Height = max(size.Height, sublayoutSize.Height)
	}

	size = size.Outset(s.Padding)

	if s.FlexBase > size.Width {
		size.Width = s.FlexBase
	}

	if size.Width == width {
		return size.Height
	}

	if size.Width < width {
		growSum := float32(0.0)
		for _, setting := range settings {
			growSum += setting.FlexGrow
		}

		if growSum == 0 {
			return size.Height
		}

		delta := width - size.Width
		size = Size{}
		for i, sublayout := range subLayouters {
			if i > 0 {
				size.Width += s.FlexGap
			}

			sublayoutDelta := uint32(float32(delta) * settings[i].FlexGrow / growSum)
			if sublayoutDelta == 0 {
				sizes[i].Width += sublayoutDelta
				sizes[i].Height = sublayout.HeightForWidth(sizes[i].Width)
			}
			size.Width += sizes[i].Width
			size.Height = max(size.Height, sizes[i].Height)
		}

		size = size.Outset(s.Padding)
		return size.Height
	}

	return 0
}

func (flexRow) Layout(s LayoutSettings, size Size, subLayouters ...SubLayouter) {}
