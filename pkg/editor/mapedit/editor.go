package mapedit

import (
	"fmt"

	"g80/pkg/cart"
	"g80/pkg/editor"
	"g80/pkg/editor/spriteedit"
	"g80/pkg/gfx"
	"g80/pkg/sprite"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui/widget"
)

var (
	iconMap      = gfx.NewBitmap("###.###", "#.#.#.#", "###.###", ".......", "###.###", "#.#.#.#", "###.###")
	iconSelect   = gfx.NewBitmap("##.##.#", "#.....#", ".......", "#.....#", "#.....#", ".......", "##.##.#")
	iconPrevious = gfx.NewBitmap("..#", ".#.", "#..", ".#.", "..#")
	iconNext     = gfx.NewBitmap("#..", ".#.", "..#", ".#.", "#..")
	iconAssign   = gfx.NewBitmap("...#...", "...##..", "#######", "...##..", "...#...")
)

type Editor struct {
	st   *state
	root *widget.Container
	img  *sprite.Renderer
}

var _ editor.Commands = (*Editor)(nil)

func New(c *cart.Cart) editor.Editor {
	e := &Editor{st: newState(c), img: sprite.NewRenderer(c.Sprites, &c.Palette)}
	e.root = ui.NewPanel(0, 0, ui.ContentWidth, ui.ContentHeight, ui.ColorBackground)
	e.build()
	return e
}

func (e *Editor) Name() string                               { return fmt.Sprintf("map editor #%02d", e.st.active) }
func (e *Editor) Icon() *gfx.Bitmap                          { return iconMap }
func (e *Editor) Root() widget.PreferredSizeLocateableWidget { return e.root }
func (e *Editor) Undo()                                      { e.st.undo() }
func (e *Editor) Redo()                                      { e.st.redo() }
func (e *Editor) Copy()                                      { e.st.copy() }
func (e *Editor) Cut()                                       { e.st.cut() }
func (e *Editor) Paste()                                     { e.st.paste() }

func (e *Editor) build() {
	e.root.RemoveChildren()
	st := e.st
	style := ui.ButtonStyle{Fg: ui.ColorDark, FgHover: ui.ColorMuted, FgActive: ui.ColorWhite}
	if st.overview {
		e.root.AddChild(
			ui.NewPanel(0, 0, ui.ContentWidth, 8, ui.ColorPanel,
				ui.NewIconButton(iconPrevious, 8, 8, style, func() { st.overview = false; e.build() }, ui.At(0, 0)),
				ui.NewLabel(120, ui.AlignLeft, ui.ColorWhite, func() string { return fmt.Sprintf("MAP #%02d", st.active) }, ui.At(12, 1))),
			newOverview(e, ui.At(0, 8)),
		)
		return
	}
	var tools []widget.PreferredSizeLocateableWidget
	icons := []*gfx.Bitmap{spriteedit.Tools[0].Icon(), spriteedit.Tools[1].Icon(), spriteedit.Tools[2].Icon(), iconSelect}
	for index, icon := range icons {
		button := ui.NewIconButton(icon, 8, 8, style, func() { st.tool = tool(index); st.pasting = false })
		button.Active = func() bool { return st.tool == tool(index) }
		tools = append(tools, button)
	}
	grid := ui.NewIconButton(iconMap, 8, 8, style, func() { st.grid = !st.grid }, ui.At(54, 120))
	grid.Active = func() bool { return st.grid }
	e.root.AddChild(ui.NewPanel(0, 0, sidebarWidth, ui.ContentHeight, ui.ColorPanel,
		ui.NewIconButton(iconMap, 8, 8, style, func() { st.commit(); st.overview = true; e.build() }, ui.At(3, 2)),
		ui.NewLabel(56, ui.AlignLeft, ui.ColorWhite, func() string { return fmt.Sprintf("MAP #%02d", st.active) }, ui.At(18, 3)),
		ui.NewRow(widget.DirectionVertical, 4, ui.At(3, 20), tools...),
		newPalette(e, ui.At(18, 15)),
		ui.NewIconButton(iconAssign, 8, 8, style, func() { st.choosing = !st.choosing; st.sheet = int(st.current().Sprites[st.slot]) / 256; e.build() }, ui.At(68, 16)),
		ui.NewLabel(60, ui.AlignLeft, ui.ColorWhite, func() string {
			if st.status != "" {
				return st.status
			}
			if st.pasting {
				return fmt.Sprintf("PASTE %dX%d", st.clipSize.X, st.clipSize.Y)
			}
			return fmt.Sprintf("%02d #%04d", st.slot, st.current().Sprites[st.slot])
		}, ui.At(18, 113)),
		ui.NewIconButton(iconPrevious, 8, 8, style, func() { st.offset = max(0, st.offset-4) }, ui.At(18, 120)),
		ui.NewIconButton(iconNext, 8, 8, style, func() { st.offset = min(maxOffset, st.offset+4) }, ui.At(30, 120)),
		grid,
	))
	if st.choosing {
		e.root.AddChild(
			ui.NewLabel(30, ui.AlignCenter, ui.ColorWhite, func() string { return fmt.Sprintf("SHEET %X", st.sheet) }, ui.At(81, 2)),
			ui.NewIconButton(iconPrevious, 8, 8, style, func() { st.sheet = (st.sheet + cart.SpriteSheetCount - 1) % cart.SpriteSheetCount }, ui.At(84, 13)),
			ui.NewIconButton(iconNext, 8, 8, style, func() { st.sheet = (st.sheet + 1) % cart.SpriteSheetCount }, ui.At(99, 13)),
			newSpritePicker(e, ui.At(112, 1)),
		)
	} else {
		e.root.AddChild(newCanvas(e, ui.At(sidebarWidth, 0)))
	}
}
