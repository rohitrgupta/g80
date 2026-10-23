package spriteedit

import (
	"image"

	"g80/pkg/cart"
	"g80/pkg/sprite"
)

// PageTiles is the width and height, in tiles, of one sprite sheet page shown at a time.
const PageTiles = 16

const pageSize = PageTiles * PageTiles

// Sizes lists the selectable edit region sizes in pixels.
var Sizes = []int{8, 16}

// State is shared by all sprite editor widgets.
type State struct {
	Cart    *cart.Cart
	Sheet   *sprite.Sheet
	History *sprite.History

	Tool               Tool
	Primary, Secondary uint8
	Size               int
	Page               int
	selected           int
}

func NewState(c *cart.Cart) *State {
	return &State{
		Cart:      c,
		Sheet:     c.Sprites,
		History:   sprite.NewHistory(256),
		Tool:      Tools[0],
		Primary:   12,
		Secondary: 0,
		Size:      Sizes[0],
	}
}

func (s *State) Pages() int { return s.Sheet.TilesY() / PageTiles }

// Selected is the index of the top-left tile of the edit region.
func (s *State) Selected() int { return s.selected }

// Select makes tile index the top-left of the edit region, shifting it so the region stays on its page.
func (s *State) Select(index int) {
	if index < 0 || index >= s.Sheet.NumTiles() {
		return
	}
	page, local := index/pageSize, index%pageSize
	limit := PageTiles - s.Size/sprite.TileSize
	tx, ty := min(local%PageTiles, limit), min(local/PageTiles, limit)
	s.Page = page
	s.selected = page*pageSize + ty*PageTiles + tx
}

func (s *State) SetSize(n int) {
	s.Size = n
	s.Select(s.selected)
}

func (s *State) SetPage(p int) {
	if p >= 0 && p < s.Pages() {
		s.Select(p*pageSize + s.selected%pageSize)
	}
}

// Region is the sheet rectangle currently being edited.
func (s *State) Region() image.Rectangle {
	min := s.Sheet.TileRect(s.selected).Min
	return image.Rectangle{Min: min, Max: min.Add(image.Pt(s.Size, s.Size))}
}

// Edit runs fn on the current region as a single undoable step.
func (s *State) Edit(fn func(r image.Rectangle)) {
	r := s.Region()
	s.History.Begin(s.Sheet, r)
	fn(r)
	s.History.Commit(s.Sheet)
}
