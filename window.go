package gkit

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	window  *glfw.Window
	program uint32

	Root View
}

func (w *Window) glSetup() error {
	w.window.MakeContextCurrent()
	program, err := loadShaders()
	if err != nil {
		return err
	}
	w.program = program

	return nil
}

func (w *Window) Destroy() {
	w.window.MakeContextCurrent()
	gl.DeleteProgram(w.program)
}

func (w *Window) Draw() {
	w.window.MakeContextCurrent()
	gl.UseProgram(w.program)
	defer gl.UseProgram(0)

	gl.EnableVertexAttribArray(0)

	gl.Clear(gl.COLOR_BUFFER_BIT)
	w.Root.Draw()
}

func (w *Window) Maximize() error {
	return w.window.Maximize()
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}
