package ui

import (
	"image/color"

	"g80/pkg/gfx"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Label draws dynamic text with the built-in pixel font.
type Label struct {
	Base
	Text   func() string
	Color  color.Color
	Shadow color.Color
	Align  Align
}

func NewLabel(width int, align Align, c color.Color, text func() string, opts ...widget.WidgetOpt) *Label {
	return &Label{Base: NewBase(width, gfx.GlyphH, opts...), Text: text, Color: c, Align: align}
}

func (l *Label) Render(screen *ebiten.Image) {
	l.Base.Render(screen)
	r := l.Rect()
	s := l.Text()
	x := r.Min.X
	switch l.Align {
	case AlignCenter:
		x += (r.Dx() - gfx.TextWidth(s)) / 2
	case AlignRight:
		x += r.Dx() - gfx.TextWidth(s)
	}
	if l.Shadow != nil {
		gfx.DrawText(screen, s, x+1, r.Min.Y+1, l.Shadow)
	}
	gfx.DrawText(screen, s, x, r.Min.Y, l.Color)
}
