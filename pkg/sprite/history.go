package sprite

import (
	"bytes"
	"image"
)

// History records region snapshots of a Sheet for undo/redo.
type History struct {
	limit      int
	undo, redo []change
	pending    *change
}

type change struct {
	rect          image.Rectangle
	before, after []uint8
}

func NewHistory(limit int) *History {
	return &History{limit: limit}
}

// Begin snapshots r before an edit. Call Commit when the edit is finished.
func (h *History) Begin(s *Sheet, r image.Rectangle) {
	h.pending = &change{rect: r, before: s.Read(r)}
}

// Commit records the pending edit if it changed anything.
func (h *History) Commit(s *Sheet) {
	c := h.pending
	h.pending = nil
	if c == nil {
		return
	}
	c.after = s.Read(c.rect)
	if bytes.Equal(c.before, c.after) {
		return
	}
	h.undo = append(h.undo, *c)
	if len(h.undo) > h.limit {
		h.undo = h.undo[1:]
	}
	h.redo = h.redo[:0]
}

func (h *History) Undo(s *Sheet) {
	if len(h.undo) == 0 {
		return
	}
	c := h.undo[len(h.undo)-1]
	h.undo = h.undo[:len(h.undo)-1]
	s.Write(c.rect, c.before)
	h.redo = append(h.redo, c)
}

func (h *History) Redo(s *Sheet) {
	if len(h.redo) == 0 {
		return
	}
	c := h.redo[len(h.redo)-1]
	h.redo = h.redo[:len(h.redo)-1]
	s.Write(c.rect, c.after)
	h.undo = append(h.undo, c)
}
