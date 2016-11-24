package controls

import (
	"image"

	"github.com/alex-ac/gkit"
)

type Image struct {
	gkit.ViewBase

	image image.Image
}

func NewImage() *Image {
	image := &Image{}
	image.ViewBase.View = image
	return image
}

func (i *Image) SetImage(img image.Image) {
	width, height := 0, 0
	if i.image != nil {
		bounds := i.image.Bounds()
		width = bounds.Max.X - bounds.Min.X
		height = bounds.Max.Y - bounds.Min.Y
	}
	i.image = img
	newWidth, newHeight := 0, 0
	if i.image != nil {
		bounds := i.image.Bounds()
		newWidth = bounds.Max.X - bounds.Min.X
		newHeight = bounds.Max.Y - bounds.Min.Y
	}
	if width != newWidth || height != newHeight {
		i.SetPrefSize(gkit.Size{uint32(newWidth), uint32(newHeight)})
	}
}

func (i *Image) UpdateSizes() {
	width, height := 0, 0
	if i.image != nil {
		bounds := i.image.Bounds()
		width = bounds.Max.X - bounds.Min.X
		height = bounds.Max.Y - bounds.Min.Y
	}
	i.SetPrefSize(gkit.Size{uint32(width), uint32(height)})
}

func (i *Image) Layout() {}

func (i *Image) Update() {}

func (i *Image) Draw(p gkit.Painter) {
	if i.image != nil {
		size := i.Size()
		p.DrawImage(0, 0, size.Width, size.Height, i.image)
	}
}
