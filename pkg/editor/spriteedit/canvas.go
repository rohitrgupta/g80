package spriteedit

import (
	"image"

	"g80/pkg/gfx"
	"g80/pkg/sprite"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const canvasSize = 64

// Canvas shows the edit region zoomed in and forwards mouse input to the current tool.
type Canvas struct {
	ui.Base
	st       *State
	img      *sprite.Renderer
	stroke   *Stroke
	button   ebiten.MouseButton
	hover    image.Point
	hovering bool
}

func newCanvas(st *State, img *sprite.Renderer, opts ...widget.WidgetOpt) *Canvas {
	return &Canvas{Base: ui.NewBase(canvasSize, canvasSize, opts...), st: st, img: img}
}

func (c *Canvas) scale() int { return canvasSize / c.st.Size }

// cell converts screen coordinates to a sheet pixel, clamped to the edit region.
func (c *Canvas) cell(mx, my int) image.Point {
	r, reg, s := c.Rect(), c.st.Region(), c.scale()
	x := min(max((mx-r.Min.X)/s, 0), reg.Dx()-1)
	y := min(max((my-r.Min.Y)/s, 0), reg.Dy()-1)
	return reg.Min.Add(image.Pt(x, y))
}

func (c *Canvas) Update(u *widget.UpdateObject) {
	c.Base.Update(u)
	mx, my := ebiten.CursorPosition()
	c.hovering = image.Pt(mx, my).In(c.Rect())
	c.hover = c.cell(mx, my)

	if c.stroke == nil {
		if !c.hovering {
			return
		}
		for _, b := range []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight} {
			if inpututil.IsMouseButtonJustPressed(b) {
				c.begin(b)
				return
			}
		}
		return
	}

	if !ebiten.IsMouseButtonPressed(c.button) {
		c.st.History.Commit(c.st.Sheet)
		c.stroke = nil
		return
	}
	if c.hover != c.stroke.Last {
		c.st.Tool.Drag(c.stroke, c.hover)
		c.stroke.Last = c.hover
	}
}

func (c *Canvas) begin(b ebiten.MouseButton) {
	st := c.st
	s := &Stroke{State: st, Region: st.Region(), Color: st.Primary, Last: c.hover}
	if b == ebiten.MouseButtonRight {
		s.Color, s.Secondary = st.Secondary, true
	}
	st.History.Begin(st.Sheet, s.Region)
	st.Tool.Press(s, c.hover)
	c.stroke, c.button = s, b
}

func (c *Canvas) Render(screen *ebiten.Image) {
	c.Base.Render(screen)
	r, reg, s := c.Rect(), c.st.Region(), c.scale()

	gfx.StrokeRect(screen, r.Inset(-1), ui.ColorBackground)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s), float64(s))
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	screen.DrawImage(c.img.Image().SubImage(reg).(*ebiten.Image), op)

	if c.hovering || c.stroke != nil {
		cell := image.Rect(0, 0, s, s).Add(r.Min).Add(c.hover.Sub(reg.Min).Mul(s))
		gfx.StrokeRect(screen, cell, ui.ColorWhite)
	}
}
