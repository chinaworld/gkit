package gkit

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	window  *glfw.Window
	painter *glPainter

	Size Size

	Root View
}

func (w *Window) glSetup() error {
	w.window.MakeContextCurrent()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

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

func (w *Window) UpdateSize() {
	viewportSize := make([]int32, 4)
	gl.GetIntegerv(gl.VIEWPORT, &viewportSize[0])
	w.Size.Width, w.Size.Height = uint32(viewportSize[2]), uint32(viewportSize[3])
}

func (w *Window) Draw() {
	w.window.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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

func (w *Window) Layout() {
	w.Root.Layouter().Layout(Rect{0, 0, w.Size.Width, w.Size.Height})
}
