package spriteedit

import (
	"testing"

	"g80/pkg/cart"
	"g80/pkg/ui"
)

func TestSpriteSheets(t *testing.T) {
	st := NewState(cart.New())
	if st.Pages() != 16 || pageSize != 256 {
		t.Fatalf("got %d sheets of %d sprites, want 16 sheets of 256", st.Pages(), pageSize)
	}
	if len(pageIcons) != st.Pages() {
		t.Fatalf("got %d tab icons for %d sheets", len(pageIcons), st.Pages())
	}
	tabs := (&Editor{st: st}).buildPageTabs(ui.At(102, 1))
	_, height := tabs.PreferredSize()
	if height+1 > ui.ContentHeight {
		t.Fatalf("tabs extend to %d, beyond content height %d", height+1, ui.ContentHeight)
	}
	st.Select(17)
	for sheet := 0; sheet < st.Pages(); sheet++ {
		st.SetPage(sheet)
		if st.Page != sheet || st.Selected() != sheet*pageSize+17 {
			t.Fatalf("sheet %d: got page %d, sprite %d", sheet, st.Page, st.Selected())
		}
		st.Sheet.Set(st.Region().Min.X, st.Region().Min.Y, uint8(sheet))
	}
	for sheet := 0; sheet < st.Pages(); sheet++ {
		st.SetPage(sheet)
		if got := st.Sheet.At(st.Region().Min.X, st.Region().Min.Y); got != uint8(sheet) {
			t.Errorf("sheet %d: got pixel %d", sheet, got)
		}
	}
}

func TestSpriteSheetSelectionBounds(t *testing.T) {
	st := NewState(cart.New())
	st.Select(4095)
	if st.Page != 15 || st.Selected() != 4095 {
		t.Fatalf("last sprite: got page %d, sprite %d", st.Page, st.Selected())
	}
	st.SetSize(16)
	if st.Selected() != 4078 || !st.Region().In(st.Sheet.Bounds()) {
		t.Fatalf("large selection extends beyond final sheet: sprite %d, region %v", st.Selected(), st.Region())
	}
	for _, page := range []int{-1, 16} {
		st.SetPage(page)
		if st.Page != 15 || st.Selected() != 4078 {
			t.Fatalf("invalid sheet %d changed selection", page)
		}
	}
}
