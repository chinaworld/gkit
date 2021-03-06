package gl

import (
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/alex-ac/gkit"
)

func NewWindowSystem() (gkit.WindowSystem, error) {
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
	return &WindowSystem{
		stopWait: make(chan struct{}),
	}, nil
}

type WindowSystem struct {
	glInited bool
	stopWait chan struct{}
}

var _ gkit.WindowSystem = &WindowSystem{}

func (s *WindowSystem) Create(w, h uint32, title string) (gkit.Window, error) {
	ok := false
	glfwWindow, err := glfw.CreateWindow(int(w), int(h), title, nil, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !ok {
			glfwWindow.Destroy()
		}
	}()
	if !s.glInited {
		glfwWindow.MakeContextCurrent()
		err = gl.Init()
		if err != nil {
			return nil, err
		}
		s.glInited = true
	}

	window := &Window{
		window: glfwWindow,
	}

	err = window.glSetup()
	if err != nil {
		return nil, err
	}

	ok = true
	return window, nil
}

func (s *WindowSystem) WaitEvents() {
	select {
	case <-s.stopWait:
		return
	case <-time.After(1 * time.Second / 60):
		glfw.PollEvents()
	}
}

func (s *WindowSystem) Interrupt() {
	s.stopWait <- struct{}{}
}

func (s *WindowSystem) Terminate() {
	glfw.Terminate()
}
