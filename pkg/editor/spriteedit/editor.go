package spriteedit

import (
	"fmt"
	"image"

	"g80/pkg/cart"
	"g80/pkg/editor"
	"g80/pkg/gfx"
	"g80/pkg/sprite"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui/widget"
)

// Editor is the sprite editor: canvas, tools, palette and sprite sheet.
type Editor struct {
	st       *State
	root     *widget.Container
	clip     []uint8
	clipSize int
}

var _ editor.Commands = (*Editor)(nil)

func New(c *cart.Cart) editor.Editor {
	e := &Editor{st: NewState(c)}
	e.root = e.build(sprite.NewRenderer(c.Sprites, &c.Palette))
	return e
}

func (e *Editor) Name() string                               { return "sprite editor" }
func (e *Editor) Icon() *gfx.Bitmap                          { return iconEditor }
func (e *Editor) Root() widget.PreferredSizeLocateableWidget { return e.root }

const leftPanelWidth = 110

func (e *Editor) build(img *sprite.Renderer) *widget.Container {
	st := e.st

	index := ui.NewLabel(canvasSize, ui.AlignCenter, ui.ColorWhite,
		func() string { return fmt.Sprintf("#%d", st.Selected()) }, ui.At(24, 5))
	index.Shadow = ui.ColorBackground

	left := ui.NewPanel(0, 0, leftPanelWidth, ui.ContentHeight, ui.ColorPanel,
		index,
		newSizeSlider(st, ui.At(9, 40)),
		newCanvas(st, img, ui.At(24, 13)),
		e.buildToolbar(ui.At(19, 89)),
		newPaletteView(st, ui.At(24, 104)),
		e.buildPageTabs(ui.At(102, 1)),
	)
	sheet := ui.NewPanel(leftPanelWidth, 0, ui.ContentWidth-leftPanelWidth, ui.ContentHeight, ui.ColorWhite,
		newSheetView(st, img, ui.At(1, 0)),
	)
	return ui.NewPanel(0, 0, ui.ContentWidth, ui.ContentHeight, nil, left, sheet)
}

func (e *Editor) buildToolbar(pos widget.WidgetOpt) *widget.Container {
	st := e.st
	toolStyle := ui.ButtonStyle{Fg: ui.ColorDark, FgHover: ui.ColorMuted, FgActive: ui.ColorWhite, Marker: ui.ColorBackground}
	actionStyle := ui.ButtonStyle{Fg: ui.ColorDark, FgHover: ui.ColorWhite}

	var buttons []widget.PreferredSizeLocateableWidget
	for _, t := range Tools {
		b := ui.NewIconButton(t.Icon(), 8, 8, toolStyle, func() { st.Tool = t })
		b.Active = func() bool { return st.Tool == t }
		buttons = append(buttons, b)
	}
	for _, a := range Actions {
		buttons = append(buttons, ui.NewIconButton(a.Icon, 8, 8, actionStyle, func() {
			st.Edit(func(r image.Rectangle) { a.Run(st.Sheet, r) })
		}))
	}
	return ui.NewRow(widget.DirectionHorizontal, 3, pos, buttons...)
}

func (e *Editor) buildPageTabs(pos widget.WidgetOpt) *widget.Container {
	st := e.st
	style := ui.ButtonStyle{Fg: ui.ColorDark, FgHover: ui.ColorMuted, FgActive: ui.ColorBackground, BgActive: ui.ColorWhite}

	var tabs []widget.PreferredSizeLocateableWidget
	for p := 0; p < st.Pages(); p++ {
		b := ui.NewIconButton(pageIcons[p%len(pageIcons)], 8, 8, style, func() { st.SetPage(p) })
		b.Active = func() bool { return st.Page == p }
		tabs = append(tabs, b)
	}
	return ui.NewRow(widget.DirectionVertical, 1, pos, tabs...)
}

func (e *Editor) Undo() { e.st.History.Undo(e.st.Sheet) }
func (e *Editor) Redo() { e.st.History.Redo(e.st.Sheet) }

func (e *Editor) Copy() {
	e.clip = e.st.Sheet.Read(e.st.Region())
	e.clipSize = e.st.Size
}

func (e *Editor) Cut() {
	e.Copy()
	e.st.Edit(func(r image.Rectangle) { sprite.Fill(e.st.Sheet, r, 0) })
}

// Paste writes the clipboard to the top-left of the region, cropping if sizes differ.
func (e *Editor) Paste() {
	if e.clip == nil {
		return
	}
	e.st.Edit(func(r image.Rectangle) {
		n := min(e.clipSize, r.Dx())
		for y := 0; y < n; y++ {
			for x := 0; x < n; x++ {
				e.st.Sheet.Set(r.Min.X+x, r.Min.Y+y, e.clip[y*e.clipSize+x])
			}
		}
	})
}
