package mapedit

import (
	"image"

	"g80/pkg/gfx"
	"g80/pkg/sprite"
	"g80/pkg/tilemap"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	sidebarWidth   = 80
	canvasWidth    = ui.ContentWidth - sidebarWidth
	visibleColumns = canvasWidth / sprite.TileSize
	maxOffset      = tilemap.Width - visibleColumns
	paletteCell    = 12
)

func drawSprite(screen *ebiten.Image, e *Editor, index int, x, y int, scale float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(e.img.Image().SubImage(e.st.cart.Sprites.TileRect(index)).(*ebiten.Image), op)
}

type palette struct {
	ui.Base
	e        *Editor
	dragging bool
	anchor   image.Point
}

func newPalette(e *Editor, pos widget.WidgetOpt) *palette {
	return &palette{Base: ui.NewBase(4*paletteCell, 8*paletteCell, pos), e: e}
}

func (p *palette) Update(update *widget.UpdateObject) {
	p.Base.Update(update)
	mx, my := ebiten.CursorPosition()
	point := image.Pt(mx, my).Sub(p.Rect().Min)
	cell := image.Pt(min(max(point.X, 0)/paletteCell, 3), min(max(point.Y, 0)/paletteCell, 7))
	if p.dragging {
		p.e.st.selectBrush(p.anchor, cell)
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.dragging = false
		}
		return
	}
	if !image.Pt(mx, my).In(p.Rect()) {
		return
	}
	left := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	right := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	if !left && !right {
		return
	}
	p.e.st.selectBrush(cell, cell)
	p.e.st.status = ""
	if left {
		p.anchor, p.dragging = cell, true
	}
	if right {
		p.e.st.sheet = int(p.e.st.current().Sprites[p.e.st.slot]) / 256
		p.e.st.choosing = true
		p.e.build()
	}
}

func (p *palette) Render(screen *ebiten.Image) {
	p.Base.Render(screen)
	for slot, index := range p.e.st.current().Sprites {
		origin := p.Rect().Min.Add(image.Pt(slot%4*paletteCell, slot/4*paletteCell))
		cell := image.Rect(0, 0, paletteCell, paletteCell).Add(origin)
		gfx.FillRect(screen, cell.Inset(1), ui.ColorBackground)
		drawSprite(screen, p.e, int(index), origin.X+2, origin.Y+2, 1)
		color := ui.ColorDark
		if image.Pt(slot%4, slot/4).In(p.e.st.brushRegion()) {
			color = ui.ColorWhite
		}
		gfx.StrokeRect(screen, cell, color)
	}
	brush := p.e.st.brushRegion()
	outline := image.Rectangle{Min: brush.Min.Mul(paletteCell).Add(p.Rect().Min), Max: brush.Max.Mul(paletteCell).Add(p.Rect().Min)}
	gfx.StrokeRect(screen, outline, ui.ColorWhite)
}

type spritePicker struct {
	ui.Base
	e *Editor
}

func newSpritePicker(e *Editor, pos widget.WidgetOpt) *spritePicker {
	return &spritePicker{Base: ui.NewBase(128, 128, pos), e: e}
}

func (p *spritePicker) Update(update *widget.UpdateObject) {
	p.Base.Update(update)
	mx, my := ebiten.CursorPosition()
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		p.e.st.choosing = false
		p.e.build()
		return
	}
	if !image.Pt(mx, my).In(p.Rect()) {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		point := image.Pt(mx, my).Sub(p.Rect().Min)
		p.e.st.assign(p.e.st.sheet*256 + point.Y/8*16 + point.X/8)
		p.e.build()
	}
}

func (p *spritePicker) Render(screen *ebiten.Image) {
	p.Base.Render(screen)
	origin := p.Rect().Min
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(origin.X), float64(origin.Y))
	page := image.Rect(0, p.e.st.sheet*128, 128, (p.e.st.sheet+1)*128)
	screen.DrawImage(p.e.img.Image().SubImage(page).(*ebiten.Image), op)
	index := int(p.e.st.current().Sprites[p.e.st.slot])
	if index/256 == p.e.st.sheet {
		cell := image.Rect(0, 0, 8, 8).Add(origin).Add(image.Pt(index%16*8, index%256/16*8))
		gfx.StrokeRect(screen, cell, ui.ColorWhite)
	}
}

type canvas struct {
	ui.Base
	e                   *Editor
	dragging            bool
	button              ebiten.MouseButton
	last, anchor, hover image.Point
	hovering            bool
}

func newCanvas(e *Editor, pos widget.WidgetOpt) *canvas {
	return &canvas{Base: ui.NewBase(canvasWidth, tilemap.Height*8, pos), e: e}
}

func (c *canvas) cell(mx, my int) image.Point {
	point := image.Pt(mx, my).Sub(c.Rect().Min)
	return image.Pt(min(max(point.X, 0)/8, visibleColumns-1)+c.e.st.offset, min(max(point.Y, 0)/8, tilemap.Height-1))
}

func (c *canvas) Update(update *widget.UpdateObject) {
	c.Base.Update(update)
	s := c.e.st
	mx, my := ebiten.CursorPosition()
	c.hovering = image.Pt(mx, my).In(c.Rect())
	c.hover = c.cell(mx, my)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if s.pasting {
			s.pasting = false
			return
		}
		s.commit()
		s.overview = true
		c.e.build()
		return
	}
	if !c.dragging {
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			s.offset = max(0, s.offset-1)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			s.offset = min(maxOffset, s.offset+1)
		}
		if !c.hovering {
			return
		}
		for _, button := range []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight} {
			if !inpututil.IsMouseButtonJustPressed(button) {
				continue
			}
			c.button, c.last, c.anchor = button, c.hover, c.hover
			if s.pasting {
				if button == ebiten.MouseButtonRight {
					s.pasting = false
					return
				}
				s.selection = image.Rectangle{Min: c.hover, Max: c.hover.Add(s.clipSize)}.Intersect(image.Rect(0, 0, tilemap.Width, tilemap.Height))
				s.paste()
				return
			}
			s.selection = image.Rectangle{Min: c.hover, Max: c.hover.Add(image.Pt(1, 1))}
			s.status = ""
			if s.tool == picker {
				if slot := s.current().At(c.hover.X, c.hover.Y); slot != tilemap.Empty {
					point := image.Pt(int(slot)%4, int(slot)/4)
					s.selectBrush(point, point)
				}
				return
			}
			if s.tool == selectTool {
				c.dragging = true
				return
			}
			s.begin()
			if s.tool == fill {
				s.flood(c.hover, button == ebiten.MouseButtonRight)
				s.commit()
				return
			}
			c.dragging = true
			s.paint(c.hover, c.hover, button == ebiten.MouseButtonRight)
			return
		}
		return
	}
	if !ebiten.IsMouseButtonPressed(c.button) {
		s.commit()
		c.dragging = false
		return
	}
	if s.tool == selectTool {
		s.selection = image.Rect(min(c.anchor.X, c.hover.X), min(c.anchor.Y, c.hover.Y), max(c.anchor.X, c.hover.X)+1, max(c.anchor.Y, c.hover.Y)+1)
	} else if c.hover != c.last {
		s.paint(c.last, c.hover, c.button == ebiten.MouseButtonRight)
		c.last = c.hover
	}
}

func (c *canvas) Render(screen *ebiten.Image) {
	c.Base.Render(screen)
	s := c.e.st
	origin := c.Rect().Min
	for row := 0; row < tilemap.Height; row++ {
		for col := 0; col < visibleColumns; col++ {
			cell := image.Rect(0, 0, 8, 8).Add(origin).Add(image.Pt(col*8, row*8))
			if slot := s.current().At(col+s.offset, row); slot != tilemap.Empty {
				drawSprite(screen, c.e, int(s.current().Sprites[slot]), cell.Min.X, cell.Min.Y, 1)
			}
			if s.grid {
				gfx.StrokeRect(screen, cell, ui.ColorDark)
			}
		}
	}
	selection := image.Rectangle{Min: s.selection.Min.Sub(image.Pt(s.offset, 0)).Mul(8).Add(origin), Max: s.selection.Max.Sub(image.Pt(s.offset, 0)).Mul(8).Add(origin)}
	selection = selection.Intersect(c.Rect())
	if !selection.Empty() {
		gfx.StrokeRect(screen, selection, ui.ColorPanel)
	}
	if c.hovering {
		size := image.Pt(1, 1)
		if s.pasting {
			size = s.clipSize
		} else if s.tool == pen {
			size = s.brushRegion().Size()
		}
		cell := image.Rectangle{Max: size.Mul(8)}.Add(origin).Add(c.hover.Sub(image.Pt(s.offset, 0)).Mul(8)).Intersect(c.Rect())
		gfx.StrokeRect(screen, cell, ui.ColorWhite)
	}
}

type overview struct {
	ui.Base
	e *Editor
}

func newOverview(e *Editor, pos widget.WidgetOpt) *overview {
	return &overview{Base: ui.NewBase(ui.ContentWidth, 120, pos), e: e}
}

func (v *overview) Update(update *widget.UpdateObject) {
	v.Base.Update(update)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		v.e.st.overview = false
		v.e.build()
		return
	}
	mx, my := ebiten.CursorPosition()
	if image.Pt(mx, my).In(v.Rect()) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		point := image.Pt(mx, my).Sub(v.Rect().Min)
		v.e.st.setMap(point.Y/15*8 + point.X/30)
		v.e.build()
	}
}

func (v *overview) Render(screen *ebiten.Image) {
	v.Base.Render(screen)
	for index, m := range v.e.st.cart.Maps {
		origin := v.Rect().Min.Add(image.Pt(index%8*30, index/8*15))
		cell := image.Rect(0, 0, 30, 15).Add(origin)
		for row := 0; row < 13; row++ {
			for col := 0; col < 28; col++ {
				slot := m.At(col*tilemap.Width/28, row*tilemap.Height/13)
				if slot == tilemap.Empty {
					continue
				}
				region := v.e.st.cart.Sprites.TileRect(int(m.Sprites[slot]))
				colorIndex := uint8(0)
				counts := [16]int{}
				for py := region.Min.Y; py < region.Max.Y; py++ {
					for px := region.Min.X; px < region.Max.X; px++ {
						pixel := v.e.st.cart.Sprites.At(px, py)
						if pixel != 0 {
							counts[pixel]++
							if counts[pixel] > counts[colorIndex] {
								colorIndex = pixel
							}
						}
					}
				}
				gfx.FillRect(screen, image.Rect(origin.X+col+1, origin.Y+row+1, origin.X+col+2, origin.Y+row+2), v.e.st.cart.Palette[colorIndex])
			}
		}
		color := ui.ColorDark
		if index == v.e.st.active {
			color = ui.ColorPanel
		}
		gfx.StrokeRect(screen, cell, color)
	}
}
