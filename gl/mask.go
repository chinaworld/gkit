package gl

import (
	"image"
	"image/color"
	"image/draw"
)

type bitmap8bpp struct {
	data []uint8
	Rect image.Rectangle
}

func newBitmap8bpp(width, height int) *bitmap8bpp {
	return &bitmap8bpp{
		data: make([]uint8, width*height),
		Rect: image.Rectangle{Max: image.Point{width, height}},
	}
}

func (b *bitmap8bpp) ColorModel() color.Model {
	return color.GrayModel
}

func (b *bitmap8bpp) Bounds() image.Rectangle {
	return b.Rect
}

func (b *bitmap8bpp) At(x, y int) color.Color {
	width := b.Rect.Max.X - b.Rect.Min.X
	return color.Gray{b.data[y*width+x]}
}

func (b *bitmap8bpp) Set(x, y int, c color.Color) {
	width := b.Rect.Max.X - b.Rect.Min.X
	b.data[y*width+x] = b.ColorModel().Convert(c).(color.Gray).Y
}

var _ draw.Image = &bitmap8bpp{}
