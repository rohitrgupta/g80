package ui

import (
	"image/color"

	"g80/pkg/gfx"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// ButtonStyle configures IconButton colors. Nil colors are not drawn / fall back to Fg/Bg.
type ButtonStyle struct {
	Fg, FgHover, FgActive color.Color
	Bg, BgActive          color.Color
	// Marker, if set, draws an arrow above the button while it is active.
	Marker color.Color
}

// IconButton is a clickable 1-bit icon, optionally showing an active (selected) state.
type IconButton struct {
	Base
	Icon    *gfx.Bitmap
	Style   ButtonStyle
	Active  func() bool
	OnClick func()
	hover   bool
}

var markerArrow = gfx.NewBitmap("#####", ".###.", "..#..")

func NewIconButton(icon *gfx.Bitmap, w, h int, style ButtonStyle, onClick func(), opts ...widget.WidgetOpt) *IconButton {
	b := &IconButton{Icon: icon, Style: style, OnClick: onClick}
	b.Base = NewBase(w, h, append(opts,
		widget.WidgetOpts.CursorEnterHandler(func(*widget.WidgetCursorEnterEventArgs) { b.hover = true }),
		widget.WidgetOpts.CursorExitHandler(func(*widget.WidgetCursorExitEventArgs) { b.hover = false }),
		widget.WidgetOpts.MouseButtonClickedHandler(func(a *widget.WidgetMouseButtonClickedEventArgs) {
			if a.Button == ebiten.MouseButtonLeft && b.OnClick != nil {
				b.OnClick()
			}
		}),
	)...)
	return b
}

func (b *IconButton) Render(screen *ebiten.Image) {
	b.Base.Render(screen)
	r := b.Rect()
	active := b.Active != nil && b.Active()

	bg, fg := b.Style.Bg, b.Style.Fg
	switch {
	case active:
		if b.Style.BgActive != nil {
			bg = b.Style.BgActive
		}
		if b.Style.FgActive != nil {
			fg = b.Style.FgActive
		}
	case b.hover && b.Style.FgHover != nil:
		fg = b.Style.FgHover
	}

	if bg != nil {
		gfx.FillRect(screen, r, bg)
	}
	iw, ih := b.Icon.Size()
	b.Icon.Draw(screen, r.Min.X+(r.Dx()-iw)/2, r.Min.Y+(r.Dy()-ih)/2, fg)

	if active && b.Style.Marker != nil {
		mw, _ := markerArrow.Size()
		markerArrow.Draw(screen, r.Min.X+(r.Dx()-mw)/2, r.Min.Y-4, b.Style.Marker)
	}
}
