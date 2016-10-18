package gkit

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	window *glfw.Window
}

func (w *Window) Draw() {
	w.window.MakeContextCurrent()
}

func (w *Window) Maximize() error {
	return w.window.Maximize()
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *Window) loadShaders() error {
	_ := gl.CreateProgram()

}
