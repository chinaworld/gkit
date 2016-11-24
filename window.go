package gkit

type WindowSystem interface {
	Create(width, height uint32, title string) (Window, error)
	PollEvents()
	Terminate()
}

type Window interface {
	DrawingContext
	SetRoot(View)
	Root() View
	Destroy()
	UpdateSize()
	Maximize() error
	ShouldClose() bool
	Size() Size
}
