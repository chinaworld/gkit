package gkit

type Application struct {
	windows      []Window
	glInited     bool
	windowSystem WindowSystem
}

func NewApplication(windowSystem WindowSystem) (*Application, error) {
	return &Application{
		windows:      make([]Window, 0),
		windowSystem: windowSystem,
	}, nil
}

func (a *Application) CreateWindow(w, h int, title string) (Window, error) {
	window, err := a.windowSystem.Create(uint32(w), uint32(h), title)
	if err != nil {
		return nil, err
	}
	a.windows = append(a.windows, window)
	return window, nil
}

func (a *Application) Shutdown() {
	for _, window := range a.windows {
		window.Destroy()
	}
	a.windowSystem.Terminate()
}

func (a *Application) Run() {
	for len(a.windows) > 0 {
		a.windowSystem.PollEvents()
		shouldCleanUp := false
		for _, window := range a.windows {
			shouldCleanUp = shouldCleanUp || window.ShouldClose()
		}
		if shouldCleanUp {
			keepedWindows := make([]Window, 0, len(a.windows)-1)
			for _, window := range a.windows {
				if !window.ShouldClose() {
					keepedWindows = append(keepedWindows, window)
				} else {
					window.Destroy()
				}
			}
			a.windows = keepedWindows
		}

		drawQueue := make(chan struct {
			drawingContext DrawingContext
			painter        Painter
		})
		for _, window := range a.windows {
			window.UpdateSize()
			go func(window Window) {
				window.Root().PropagateUpdate()
				window.Root().PropagateLayout()
				painter := window.BeginPaint()
				painter.DrawLayer(Rect{Size: window.Size()}, window.Root())
				drawQueue <- struct {
					drawingContext DrawingContext
					painter        Painter
				}{
					drawingContext: window,
					painter:        painter,
				}
			}(window)
		}

		for range a.windows {
			drawEntry := <-drawQueue
			context := drawEntry.drawingContext
			painter := drawEntry.painter
			context.EndPaint(painter)
		}
		close(drawQueue)
	}
}
