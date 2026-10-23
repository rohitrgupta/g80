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

const sheetViewSize = PageTiles * sprite.TileSize

// SheetView shows one page of the sprite sheet; click or drag to select the edit region.
type SheetView struct {
	ui.Base
	st       *State
	img      *sprite.Renderer
	dragging bool
}

func newSheetView(st *State, img *sprite.Renderer, opts ...widget.WidgetOpt) *SheetView {
	return &SheetView{Base: ui.NewBase(sheetViewSize, sheetViewSize, opts...), st: st, img: img}
}

func (v *SheetView) Update(u *widget.UpdateObject) {
	v.Base.Update(u)
	mx, my := ebiten.CursorPosition()
	r := v.Rect()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && image.Pt(mx, my).In(r) {
		v.dragging = true
	}
	if !v.dragging {
		return
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		v.dragging = false
		return
	}
	tx := min(max((mx-r.Min.X)/sprite.TileSize, 0), PageTiles-1)
	ty := min(max((my-r.Min.Y)/sprite.TileSize, 0), PageTiles-1)
	v.st.Select(v.st.Page*pageSize + ty*PageTiles + tx)
}

func (v *SheetView) Render(screen *ebiten.Image) {
	v.Base.Render(screen)
	r := v.Rect()
	pageOrigin := image.Pt(0, v.st.Page*sheetViewSize)
	page := image.Rectangle{Min: pageOrigin, Max: pageOrigin.Add(image.Pt(sheetViewSize, sheetViewSize))}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	screen.DrawImage(v.img.Image().SubImage(page).(*ebiten.Image), op)

	sel := v.st.Region().Sub(pageOrigin).Add(r.Min)
	gfx.StrokeRect(screen, sel.Inset(-1), ui.ColorWhite)
}
