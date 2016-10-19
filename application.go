package gkit

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Application struct {
	windows  []*Window
	glInited bool
}

func NewApplication() (*Application, error) {
	ok := false
	err := glfw.Init()
	if err != nil {
		return nil, err
	}
	defer func() {
		if !ok {
			glfw.Terminate()
		}
	}()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	ok = true
	return &Application{
		windows: make([]*Window, 0),
	}, nil
}

func (a *Application) CreateWindow(w, h int, title string) (*Window, error) {
	ok := false
	glfwWindow, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !ok {
			glfwWindow.Destroy()
		}
	}()
	if !a.glInited {
		glfwWindow.MakeContextCurrent()
		if err := gl.Init(); err != nil {
			return nil, err
		}
		a.glInited = true
	}

	window := &Window{
		window: glfwWindow,
	}

	if err := window.glSetup(); err != nil {
		return nil, err
	}

	ok = true
	a.windows = append(a.windows, window)
	return window, err
}

func (a *Application) Shutdown() {
	for _, window := range a.windows {
		window.Destroy()
	}
	glfw.Terminate()
}

func (a *Application) Run() {
	for len(a.windows) > 0 {
		glfw.PollEvents()
		for _, window := range a.windows {
			window.Draw()
		}
		shouldCleanUp := false
		for _, window := range a.windows {
			shouldCleanUp = shouldCleanUp || window.ShouldClose()
		}
		if shouldCleanUp {
			keepedWindows := make([]*Window, 0, len(a.windows)-1)
			for _, window := range a.windows {
				if !window.ShouldClose() {
					keepedWindows = append(keepedWindows, window)
				} else {
					window.Destroy()
				}
			}
			a.windows = keepedWindows
		}
	}
}
