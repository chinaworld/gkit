package gkit

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
	"os"

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

func (f *Font) DrawString(size uint32, text string) *image.Gray {
	options := truetype.Options{
		Size: float64(size),
		DPI:  0, // Fuck imperial system. 1pt = 1px now.
	}
	face := truetype.NewFace(f.font, &options)
	drawer := font.Drawer{
		Face: face,
	}
	metrics := face.Metrics()
	advance := drawer.MeasureString(text)

	height := metrics.Height.Ceil()
	width := advance.Ceil()

	dst := image.NewGray(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{width, height},
	})

	drawer.Dst = dst
	drawer.Src = dst
	drawer.Dot = fixed.Point26_6{0, metrics.Ascent}
	drawer.DrawString(text)

	return dst
}
