package mapedit

import (
	"image"
	"testing"

	"g80/pkg/cart"
	"g80/pkg/tilemap"
)

func TestPaintHistory(t *testing.T) {
	s := newState(cart.New())
	s.slot = 31
	s.begin()
	s.paint(image.Pt(0, 0), image.Pt(10, 0), false)
	s.commit()
	for col := 0; col <= 10; col++ {
		if s.current().At(col, 0) != 31 {
			t.Fatal("brush stroke has gaps")
		}
	}
	s.undo()
	if s.current().At(5, 0) != tilemap.Empty {
		t.Fatal("undo failed")
	}
	s.redo()
	if s.current().At(5, 0) != 31 {
		t.Fatal("redo failed")
	}
	s.setMap(1)
	s.undo()
	if s.current().At(5, 0) != tilemap.Empty {
		t.Fatal("map histories leaked")
	}
}

func TestPaletteAndClipboard(t *testing.T) {
	s := newState(cart.New())
	s.assign(4095)
	s.begin()
	s.paint(image.Point{}, image.Point{}, false)
	s.commit()
	s.copy()
	s.setMap(63)
	s.paste()
	slot := s.current().At(0, 0)
	if slot == tilemap.Empty || s.current().Sprites[slot] != 4095 {
		t.Fatal("cross-map paste lost the sprite assignment")
	}
	s.cut()
	if s.current().At(0, 0) != tilemap.Empty {
		t.Fatal("cut did not clear selection")
	}
	s.undo()
	if s.current().At(0, 0) == tilemap.Empty {
		t.Fatal("cut was not undoable")
	}
}

func TestFlood(t *testing.T) {
	s := newState(cart.New())
	s.begin()
	s.flood(image.Point{}, false)
	s.commit()
	if s.current().At(31, 15) != 0 {
		t.Fatal("fill did not reach the map boundary")
	}
	s.undo()
	if s.current().At(31, 15) != tilemap.Empty {
		t.Fatal("fill undo failed")
	}
}

func TestPasteFullPaletteIsAtomic(t *testing.T) {
	s := newState(cart.New())
	for slot := 0; slot < tilemap.PaletteSize; slot++ {
		s.current().Set(slot, 0, uint8(slot))
	}
	before := *s.current()
	s.clipboard = []uint16{4095}
	s.clipSize = image.Pt(1, 1)
	s.paste()
	if *s.current() != before || s.status != "PALETTE FULL" {
		t.Fatal("paste with a full palette must not corrupt the map")
	}
}

func TestPaletteAssignmentHistory(t *testing.T) {
	s := newState(cart.New())
	s.slot = 31
	s.assign(4095)
	s.undo()
	if s.current().Sprites[31] != 31 {
		t.Fatal("palette assignment undo failed")
	}
	s.redo()
	if s.current().Sprites[31] != 4095 {
		t.Fatal("palette assignment redo failed")
	}
	s.assign(4096)
	if s.current().Sprites[31] != 4095 {
		t.Fatal("out-of-range sprite assignment accepted")
	}
}

func TestRepeatedPaste(t *testing.T) {
	s := newState(cart.New())
	s.current().Set(0, 0, 3)
	s.current().Set(1, 0, 7)
	s.selection = image.Rect(0, 0, 2, 1)
	s.copy()
	for _, origin := range []image.Point{image.Pt(4, 3), image.Pt(8, 5), image.Pt(31, 15)} {
		s.selection = image.Rectangle{Min: origin, Max: origin.Add(image.Pt(1, 1))}
		s.paste()
		if !s.pasting || len(s.clipboard) != 2 || s.current().At(origin.X, origin.Y) != 3 {
			t.Fatal("paste did not preserve the clipboard and repeat placement")
		}
		if origin.X < 31 && s.current().At(origin.X+1, origin.Y) != 7 {
			t.Fatal("paste lost part of the copied rectangle")
		}
	}
	s.undo()
	if s.current().At(31, 15) != tilemap.Empty || s.current().At(8, 5) != 3 {
		t.Fatal("repeated placements must undo independently")
	}
}

func TestRectangularBrush(t *testing.T) {
	s := newState(cart.New())
	s.selectBrush(image.Pt(2, 3), image.Pt(1, 2))
	s.begin()
	s.paint(image.Pt(5, 4), image.Pt(5, 4), false)
	s.commit()
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			if got := s.current().At(5+col, 4+row); got != uint8((row+2)*4+col+1) {
				t.Fatalf("brush cell %d,%d: got slot %d", col, row, got)
			}
		}
	}
	s.undo()
	if s.current().At(6, 5) != tilemap.Empty {
		t.Fatal("rectangle stamp must undo as one change")
	}
	s.begin()
	s.paint(image.Pt(31, 15), image.Pt(31, 15), false)
	s.commit()
	if s.current().At(31, 15) != 9 {
		t.Fatal("rectangular brush must clip at map edges")
	}
	s.begin()
	s.paint(image.Pt(31, 15), image.Pt(31, 15), true)
	s.commit()
	if s.current().At(31, 15) != tilemap.Empty {
		t.Fatal("rectangular eraser failed")
	}
}
