package gkit

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type Font struct {
	font *truetype.Font
}

func LoadFontFile(filename string) (*Font, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return LoadFont(data)
}

func LoadFont(data []byte) (*Font, error) {
	font, err := freetype.ParseFont(data)
	if err != nil {
		return nil, err
	}
	return &Font{
		font: font,
	}, nil
}

func (f *Font) drawer(size uint32) font.Drawer {
	options := truetype.Options{
		Size: float64(size),
		DPI:  0, // Fuck imperial system. 1pt = 1px now.
	}
	face := truetype.NewFace(f.font, &options)
	return font.Drawer{
		Face: face,
	}
}

func (f *Font) StringSize(size uint32, text string) Size {
	drawer := f.drawer(size)
	metrics := drawer.Face.Metrics()
	advance := drawer.MeasureString(text)
	height := metrics.Height.Ceil()
	width := advance.Ceil()
	return Size{
		Width:  uint32(width),
		Height: uint32(height),
	}
}

func (f *Font) DrawString(size uint32, text string, x, y uint32, dst draw.Image) {
	drawer := f.drawer(size)
	metrics := drawer.Face.Metrics()

	drawer.Dst = dst
	drawer.Src = image.NewUniform(color.Gray{0xff})
	drawer.Dot = fixed.Point26_6{
		X: fixed.I(int(x)),
		Y: fixed.I(int(y)) + metrics.Ascent,
	}
	drawer.DrawString(text)
}
