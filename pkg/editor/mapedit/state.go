package mapedit

import (
	"image"

	"g80/pkg/cart"
	"g80/pkg/tilemap"
)

type tool int

const (
	pen tool = iota
	picker
	fill
	selectTool
)

type change struct {
	before, after tilemap.Map
}

type history struct {
	undo, redo []change
}

type state struct {
	cart               *cart.Cart
	active             int
	slot               uint8
	tool               tool
	selection          image.Rectangle
	offset             int
	overview, choosing bool
	sheet              int
	grid               bool
	status             string
	history            [tilemap.Count]history
	pending            *tilemap.Map
	pendingMap         int
	clipboard          []uint16
	clipSize           image.Point
	pasting            bool
	brush              image.Rectangle
}

func newState(c *cart.Cart) *state {
	return &state{cart: c, selection: image.Rect(0, 0, 1, 1), grid: true}
}

func (s *state) current() *tilemap.Map { return s.cart.Maps[s.active] }

func (s *state) brushRegion() image.Rectangle {
	if s.brush.Empty() {
		origin := image.Pt(int(s.slot)%4, int(s.slot)/4)
		return image.Rectangle{Min: origin, Max: origin.Add(image.Pt(1, 1))}
	}
	return s.brush
}

func (s *state) selectBrush(from, to image.Point) {
	s.brush = image.Rect(min(from.X, to.X), min(from.Y, to.Y), max(from.X, to.X)+1, max(from.Y, to.Y)+1).Intersect(image.Rect(0, 0, 4, 8))
	if !s.brush.Empty() {
		s.slot = uint8(s.brush.Min.Y*4 + s.brush.Min.X)
	}
	s.pasting = false
}

func (s *state) begin() {
	s.commit()
	before := *s.current()
	s.pending, s.pendingMap = &before, s.active
}

func (s *state) commit() {
	if s.pending == nil {
		return
	}
	after := *s.cart.Maps[s.pendingMap]
	if after != *s.pending {
		h := &s.history[s.pendingMap]
		h.undo = append(h.undo, change{*s.pending, after})
		if len(h.undo) > 128 {
			h.undo = h.undo[1:]
		}
		h.redo = nil
	}
	s.pending = nil
}

func (s *state) setMap(index int) {
	if index < 0 || index >= tilemap.Count {
		return
	}
	s.commit()
	s.active, s.overview, s.choosing = index, false, false
	s.status = ""
}

func (s *state) assign(index int) {
	if index < 0 || index >= s.cart.Sprites.NumTiles() {
		return
	}
	s.begin()
	s.current().Sprites[s.slot] = uint16(index)
	s.commit()
	s.choosing = false
}

func (s *state) paint(from, to image.Point, erase bool) {
	brush := s.brushRegion()
	delta := to.Sub(from)
	steps := max(abs(delta.X), abs(delta.Y), 1)
	for step := 0; step <= steps; step++ {
		point := from.Add(image.Pt(delta.X*step/steps, delta.Y*step/steps))
		for row := 0; row < brush.Dy(); row++ {
			for col := 0; col < brush.Dx(); col++ {
				slot := uint8((brush.Min.Y+row)*4 + brush.Min.X + col)
				if erase {
					slot = tilemap.Empty
				}
				s.current().Set(point.X+col, point.Y+row, slot)
			}
		}
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func (s *state) flood(point image.Point, erase bool) {
	target := s.current().At(point.X, point.Y)
	replacement := s.slot
	if erase {
		replacement = tilemap.Empty
	}
	if target == replacement {
		return
	}
	bounds := image.Rect(0, 0, tilemap.Width, tilemap.Height)
	queue := []image.Point{point}
	s.current().Set(point.X, point.Y, replacement)
	for len(queue) > 0 {
		point, queue = queue[0], queue[1:]
		for _, neighbor := range []image.Point{point.Add(image.Pt(1, 0)), point.Add(image.Pt(-1, 0)), point.Add(image.Pt(0, 1)), point.Add(image.Pt(0, -1))} {
			if neighbor.In(bounds) && s.current().At(neighbor.X, neighbor.Y) == target {
				s.current().Set(neighbor.X, neighbor.Y, replacement)
				queue = append(queue, neighbor)
			}
		}
	}
}

func (s *state) undo() {
	s.commit()
	h := &s.history[s.active]
	if len(h.undo) == 0 {
		return
	}
	entry := h.undo[len(h.undo)-1]
	h.undo = h.undo[:len(h.undo)-1]
	*s.current() = entry.before
	h.redo = append(h.redo, entry)
}

func (s *state) redo() {
	s.commit()
	h := &s.history[s.active]
	if len(h.redo) == 0 {
		return
	}
	entry := h.redo[len(h.redo)-1]
	h.redo = h.redo[:len(h.redo)-1]
	*s.current() = entry.after
	h.undo = append(h.undo, entry)
}

const emptySprite = uint16(65535)

func (s *state) copy() {
	s.pasting = false
	s.clipSize = s.selection.Size()
	s.clipboard = nil
	for row := s.selection.Min.Y; row < s.selection.Max.Y; row++ {
		for col := s.selection.Min.X; col < s.selection.Max.X; col++ {
			index := emptySprite
			if slot := s.current().At(col, row); slot != tilemap.Empty {
				index = s.current().Sprites[slot]
			}
			s.clipboard = append(s.clipboard, index)
		}
	}
}

func (s *state) cut() {
	s.copy()
	s.begin()
	for row := s.selection.Min.Y; row < s.selection.Max.Y; row++ {
		for col := s.selection.Min.X; col < s.selection.Max.X; col++ {
			s.current().Set(col, row, tilemap.Empty)
		}
	}
	s.commit()
}

func (s *state) paste() {
	if len(s.clipboard) == 0 {
		return
	}
	next := *s.current()
	used := [tilemap.PaletteSize]bool{}
	for _, slot := range next.Tiles {
		if slot != tilemap.Empty {
			used[slot] = true
		}
	}
	for row := 0; row < s.clipSize.Y; row++ {
		for col := 0; col < s.clipSize.X; col++ {
			point := s.selection.Min.Add(image.Pt(col, row))
			if point.X >= tilemap.Width || point.Y >= tilemap.Height {
				continue
			}
			index := s.clipboard[row*s.clipSize.X+col]
			slot := tilemap.Empty
			if index != emptySprite {
				for candidate, spriteIndex := range next.Sprites {
					if spriteIndex == index {
						slot = uint8(candidate)
						break
					}
				}
				if slot == tilemap.Empty {
					for candidate, occupied := range used {
						if !occupied {
							slot = uint8(candidate)
							next.Sprites[slot] = index
							break
						}
					}
				}
				if slot == tilemap.Empty {
					s.status = "PALETTE FULL"
					return
				}
				used[slot] = true
			}
			next.Set(point.X, point.Y, slot)
		}
	}
	s.begin()
	*s.current() = next
	s.commit()
	s.status = ""
	s.pasting = true
}
