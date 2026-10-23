package ui

import (
	"image"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// Base implements the boilerplate of a custom fixed-size ebitenui widget.
// Embed it and override Render (and Update if needed).
type Base struct {
	widget        *widget.Widget
	width, height int
}

func NewBase(w, h int, opts ...widget.WidgetOpt) Base {
	return Base{widget: widget.NewWidget(opts...), width: w, height: h}
}

func (b *Base) GetWidget() *widget.Widget     { return b.widget }
func (b *Base) PreferredSize() (int, int)     { return b.width, b.height }
func (b *Base) SetLocation(r image.Rectangle) { b.widget.Rect = r }
func (b *Base) Validate()                     {}
func (b *Base) Update(u *widget.UpdateObject) { b.widget.Update(u) }
func (b *Base) Render(screen *ebiten.Image)   { b.widget.Render(screen) }
func (b *Base) Rect() image.Rectangle         { return b.widget.Rect }

// At places a widget at (x, y) relative to a parent using an AnchorLayout (e.g. a Panel).
func At(x, y int) widget.WidgetOpt {
	return widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionStart,
		VerticalPosition:   widget.AnchorLayoutPositionStart,
		Padding:            &widget.Insets{Left: x, Top: y},
	})
}

// NewPanel creates a fixed-size container at (x, y) whose children are positioned with At.
// bg may be nil for a transparent panel.
func NewPanel(x, y, w, h int, bg color.Color, children ...widget.PreferredSizeLocateableWidget) *widget.Container {
	opts := []widget.ContainerOpt{
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(w, h), At(x, y)),
	}
	if bg != nil {
		opts = append(opts, widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(bg)))
	}
	c := widget.NewContainer(opts...)
	if len(children) > 0 {
		c.AddChild(children...)
	}
	return c
}

// NewRow lays children out in a single row or column.
func NewRow(dir widget.Direction, spacing int, opt widget.WidgetOpt, children ...widget.PreferredSizeLocateableWidget) *widget.Container {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(dir),
			widget.RowLayoutOpts.Spacing(spacing),
		)),
		widget.ContainerOpts.WidgetOpts(opt),
	)
	c.AddChild(children...)
	return c
}
