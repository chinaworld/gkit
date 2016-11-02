package gkit

import (
	"image"
	"image/png"
	"os"
)

func LoadPng(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}
