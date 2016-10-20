package gkit

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	window  *glfw.Window
	painter *glPainter

	Root View
}

func (w *Window) glSetup() error {
	w.window.MakeContextCurrent()

	painter, err := newGlPainter()
	if err != nil {
		return err
	}
	w.painter = painter

	return nil
}

func (w *Window) Destroy() {
	w.window.MakeContextCurrent()
	w.painter.Destroy()
}

func (w *Window) Draw() {
	w.window.MakeContextCurrent()
	viewportSize := make([]int32, 4)
	gl.GetIntegerv(gl.VIEWPORT, &viewportSize[0])
	width, height := uint32(viewportSize[2]), uint32(viewportSize[3])
	w.Root.Width, w.Root.Height = width, height

	gl.Clear(gl.COLOR_BUFFER_BIT)
	defer w.window.SwapBuffers()

	painter, endPaint := w.painter.BeginPaint()
	defer endPaint()

	w.Root.Draw(painter)
}

func (w *Window) Maximize() error {
	return w.window.Maximize()
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}
