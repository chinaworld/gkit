package controls

import (
	"github.com/alex-ac/gkit"
)

type Label struct {
	gkit.ViewBase

	text     string
	font     *gkit.Font
	fontSize uint32
	Color    gkit.Color
}

var _ gkit.View = &Label{}

func NewLabel() *Label {
	label := Label{}
	label.ViewBase.View = &label
	return &label
}

func (l *Label) Layout() {}
func (l *Label) Update() {}

func (l *Label) Draw(p gkit.Painter) {
	if l.font == nil || l.fontSize == 0 || l.text == "" {
		return
	}
	p.SetFont(l.font)
	p.SetFontSize(l.fontSize)
	p.SetColor(l.Color)
	p.DrawText(0, 0, l.text)
}

func (l *Label) SetFont(font *gkit.Font) {
	oldFont := l.font
	l.font = font
	if oldFont != font {
		l.SetPrefSizeChanged()
	}
}

func (l *Label) SetFontSize(size uint32) {
	oldSize := l.fontSize
	l.fontSize = size
	if oldSize != size {
		l.SetPrefSizeChanged()
	}
}

func (l *Label) SetText(text string) {
	oldText := l.text
	l.text = text
	if oldText != text {
		l.SetPrefSizeChanged()
	}
}

func (l *Label) UpdateSizes() {
	if l.font == nil || l.fontSize == 0 || l.text == "" {
		l.SetMinSize(gkit.Size{})
		l.SetPrefSize(gkit.Size{})
		l.SetMaxSize(gkit.Size{})
		return
	}
	size := l.font.StringSize(l.fontSize, l.text)
	l.SetPrefSize(size)
}
