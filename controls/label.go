package controls

import (
	"github.com/alex-ac/gkit"
)

type Label struct {
	gkit.ViewBase

	text            string
	font            *gkit.Font
	fontSize        uint32
	Color           gkit.Color
	BackgroundColor gkit.Color
}

var _ gkit.View = &Label{}

func NewLabel() *Label {
	label := Label{}
	label.ViewBase.View = &label
	return &label
}

func (l *Label) Layout() {}
func (l *Label) Update() {}

type subLayer struct {
	*Label
}

func (s subLayer) Draw(p gkit.Painter) {
	p.SetFont(s.font)
	p.SetFontSize(s.fontSize)
	p.SetColor(s.Color)
	p.DrawText(gkit.Point{}, s.text)
}

func (l *Label) Draw(p gkit.Painter) {
	if l.font == nil || l.fontSize == 0 || l.text == "" {
		return
	}
	p.SetColor(l.BackgroundColor)
	p.DrawRect(gkit.Rect{Size: l.Size()})
	p.DrawLayer(l.Bounds(), subLayer{l})
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

func (l *Label) Text() string {
	return l.text
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

func (l *Label) SetColor(color gkit.Color) {
	oldColor := l.Color
	l.Color = color
	if l.Color != oldColor {
		l.SetNeedsRedraw()
	}
}

func (l *Label) SetBackgroundColor(color gkit.Color) {
	oldColor := l.BackgroundColor
	l.BackgroundColor = color
	if l.BackgroundColor != oldColor {
		l.SetNeedsRedraw()
	}
}
