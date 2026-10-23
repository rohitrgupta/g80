package app

import (
	"log"
	"strings"

	"g80/pkg/cart"
	"g80/pkg/editor"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const statusTicks = 90

// App hosts the editors under a shared top bar and implements ebiten.Game.
type App struct {
	ui      *ebitenui.UI
	cart    *cart.Cart
	path    string
	editors []editor.Editor
	active  int
	content *widget.Container
	status  string
	ticks   int
}

// New creates the app with one editor per factory; the first is shown initially.
func New(c *cart.Cart, path string, factories ...editor.Factory) *App {
	a := &App{cart: c, path: path}
	for _, f := range factories {
		a.editors = append(a.editors, f(c))
	}

	a.content = ui.NewPanel(0, ui.BarHeight, ui.ContentWidth, ui.ContentHeight, nil)
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(ui.ColorBackground)),
	)
	root.AddChild(a.buildBar(), a.content)
	a.ui = &ebitenui.UI{Container: root}
	a.switchTo(0)
	return a
}

func (a *App) buildBar() *widget.Container {
	tabStyle := ui.ButtonStyle{Fg: ui.ColorMuted, FgHover: ui.ColorPanel, FgActive: ui.ColorBackground, BgActive: ui.ColorMuted}
	var tabs []widget.PreferredSizeLocateableWidget
	for i, e := range a.editors {
		b := ui.NewIconButton(e.Icon(), 8, ui.BarHeight, tabStyle, func() { a.switchTo(i) })
		b.Active = func() bool { return a.active == i }
		tabs = append(tabs, b)
	}

	cmdStyle := ui.ButtonStyle{Fg: ui.ColorMuted, FgHover: ui.ColorBackground}
	cmd := func(f func(editor.Commands)) func() { return func() { a.command(f) } }
	buttons := []widget.PreferredSizeLocateableWidget{
		ui.NewIconButton(iconCut, 8, ui.BarHeight, cmdStyle, cmd(editor.Commands.Cut)),
		ui.NewIconButton(iconCopy, 8, ui.BarHeight, cmdStyle, cmd(editor.Commands.Copy)),
		ui.NewIconButton(iconPaste, 8, ui.BarHeight, cmdStyle, cmd(editor.Commands.Paste)),
		ui.NewIconButton(iconUndo, 8, ui.BarHeight, cmdStyle, cmd(editor.Commands.Undo)),
		ui.NewIconButton(iconRedo, 8, ui.BarHeight, cmdStyle, cmd(editor.Commands.Redo)),
		ui.NewIconButton(iconSave, 8, ui.BarHeight, cmdStyle, a.save),
	}

	title := ui.NewLabel(100, ui.AlignLeft, ui.ColorMuted, a.title, ui.At(44, 1))
	return ui.NewPanel(0, 0, ui.ScreenWidth, ui.BarHeight, ui.ColorWhite,
		ui.NewRow(widget.DirectionHorizontal, 0, ui.At(0, 0), tabs...),
		title,
		ui.NewRow(widget.DirectionHorizontal, 0, ui.At(ui.ScreenWidth-len(buttons)*8-2, 0), buttons...),
	)
}

func (a *App) title() string {
	if a.ticks > 0 {
		return a.status
	}
	return strings.ToUpper(a.editors[a.active].Name())
}

func (a *App) setStatus(s string) {
	a.status, a.ticks = s, statusTicks
}

func (a *App) switchTo(i int) {
	a.active = i
	a.content.RemoveChildren()
	a.content.AddChild(a.editors[i].Root())
}

func (a *App) command(f func(editor.Commands)) {
	if c, ok := a.editors[a.active].(editor.Commands); ok {
		f(c)
	}
}

func (a *App) save() {
	if err := a.cart.Save(a.path); err != nil {
		log.Printf("save %s: %v", a.path, err)
		a.setStatus("SAVE FAILED")
		return
	}
	a.setStatus("SAVED")
}

func (a *App) handleShortcuts() {
	if !ebiten.IsKeyPressed(ebiten.KeyControl) && !ebiten.IsKeyPressed(ebiten.KeyMeta) {
		return
	}
	pressed := inpututil.IsKeyJustPressed
	switch {
	case pressed(ebiten.KeyS):
		a.save()
	case pressed(ebiten.KeyZ):
		a.command(editor.Commands.Undo)
	case pressed(ebiten.KeyY):
		a.command(editor.Commands.Redo)
	case pressed(ebiten.KeyX):
		a.command(editor.Commands.Cut)
	case pressed(ebiten.KeyC):
		a.command(editor.Commands.Copy)
	case pressed(ebiten.KeyV):
		a.command(editor.Commands.Paste)
	}
}

func (a *App) Update() error {
	a.ui.Update()
	a.handleShortcuts()
	if a.ticks > 0 {
		a.ticks--
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.ui.Draw(screen)
}

func (a *App) Layout(int, int) (int, int) {
	return ui.ScreenWidth, ui.ScreenHeight
}
