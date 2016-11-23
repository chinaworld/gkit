package controls

import (
	"github.com/alex-ac/gkit"
)

type Label struct {
	gkit.View

	text     string
	font     *gkit.Font
	fontSize uint32
	Color    gkit.Color
}

func NewLabel() *Label {
	var label Label
	label.View.Draw = func(v *gkit.View, p gkit.Painter) {
		label.draw(p)
	}
	return &label
}

func (l *Label) draw(p gkit.Painter) {
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
		l.updateSizes()
	}
}

func (l *Label) SetFontSize(size uint32) {
	oldSize := l.fontSize
	l.fontSize = size
	if oldSize != size {
		l.updateSizes()
	}
}

func (l *Label) SetText(text string) {
	oldText := l.text
	l.text = text
	if oldText != text {
		l.updateSizes()
	}
}

func (l *Label) updateSizes() {
	if l.font == nil || l.fontSize == 0 || l.text == "" {
		l.View.MinSize = gkit.Size{}
		l.View.PrefSize = gkit.Size{}
		l.View.MaxSize = gkit.Size{}
		return
	}
	size := l.font.StringSize(l.fontSize, l.text)
	l.View.MinSize = size
	l.View.PrefSize = size
}
