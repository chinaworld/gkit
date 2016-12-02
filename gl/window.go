package gl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/alex-ac/gkit"
)

type Window struct {
	*drawingContext

	window *glfw.Window

	size gkit.Size

	root gkit.View
}

var _ gkit.Window = &Window{}

func (w *Window) Size() gkit.Size {
	return w.size
}

func (w *Window) SetRoot(view gkit.View) {
	w.root = view
	if w.root != nil {
		w.root.SetSize(w.size)
	}
}

func (w *Window) Root() gkit.View {
	return w.root
}

func (w *Window) glSetup() error {
	w.window.MakeContextCurrent()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	context, err := newDrawingContext()
	if err != nil {
		return err
	}
	w.drawingContext = context

	return nil
}

func (w *Window) Destroy() {
	w.window.MakeContextCurrent()
	w.drawingContext.Destroy()
}

func (w *Window) BeginPaint() gkit.Painter {
	return w.drawingContext.BeginPaint(w.size)
}

func (w *Window) UpdateSize() {
	viewportSize := make([]int32, 4)
	gl.GetIntegerv(gl.VIEWPORT, &viewportSize[0])
	w.size.Width, w.size.Height = uint32(viewportSize[2]), uint32(viewportSize[3])
	w.root.SetSize(w.size)
}

func (w *Window) EndPaint(painter gkit.Painter) {
	w.window.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	defer w.window.SwapBuffers()
	w.drawingContext.EndPaint(painter)
}

func (w *Window) Maximize() error {
	return w.window.Maximize()
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}
