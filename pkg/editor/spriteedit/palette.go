package spriteedit

import (
	"image"

	"g80/pkg/gfx"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	swatchSize    = 8
	swatchColumns = 8
)

// PaletteView shows the 16 cart colors. Left click picks the primary color, right click the secondary.
type PaletteView struct {
	ui.Base
	st *State
}

func newPaletteView(st *State, opts ...widget.WidgetOpt) *PaletteView {
	p := &PaletteView{st: st}
	rows := gfx.PaletteSize / swatchColumns
	p.Base = ui.NewBase(swatchColumns*swatchSize, rows*swatchSize, append(opts,
		widget.WidgetOpts.MouseButtonPressedHandler(func(a *widget.WidgetMouseButtonPressedEventArgs) {
			i := uint8((a.OffsetY/swatchSize)*swatchColumns + a.OffsetX/swatchSize)
			switch a.Button {
			case ebiten.MouseButtonLeft:
				st.Primary = i
			case ebiten.MouseButtonRight:
				st.Secondary = i
			}
		}),
	)...)
	return p
}

func (p *PaletteView) swatch(i uint8) image.Rectangle {
	x, y := int(i)%swatchColumns*swatchSize, int(i)/swatchColumns*swatchSize
	return image.Rect(x, y, x+swatchSize, y+swatchSize).Add(p.Rect().Min)
}

func (p *PaletteView) Render(screen *ebiten.Image) {
	p.Base.Render(screen)
	gfx.StrokeRect(screen, p.Rect().Inset(-1), ui.ColorBackground)
	for i, c := range p.st.Cart.Palette {
		gfx.FillRect(screen, p.swatch(uint8(i)), c)
	}
	gfx.DottedRect(screen, p.swatch(p.st.Secondary), ui.ColorWhite, ui.ColorBackground)
	gfx.StrokeRect(screen, p.swatch(p.st.Primary), ui.ColorWhite)
}
