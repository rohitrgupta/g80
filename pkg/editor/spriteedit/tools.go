package spriteedit

import (
	"image"

	"g80/pkg/gfx"
	"g80/pkg/sprite"
)

// Stroke is one press-drag-release gesture on the canvas. Points are sheet coordinates
// inside Region.
type Stroke struct {
	State     *State
	Region    image.Rectangle
	Color     uint8
	Secondary bool // started with the right mouse button
	Last      image.Point
}

// Tool is a canvas drawing tool. Register new tools by appending to Tools.
type Tool interface {
	Icon() *gfx.Bitmap
	Press(s *Stroke, p image.Point)
	Drag(s *Stroke, p image.Point)
}

var Tools = []Tool{Pen{}, Picker{}, Fill{}}

type Pen struct{}

func (Pen) Icon() *gfx.Bitmap { return iconPen }

func (Pen) Press(s *Stroke, p image.Point) { s.State.Sheet.Set(p.X, p.Y, s.Color) }

func (Pen) Drag(s *Stroke, p image.Point) {
	sprite.Line(s.State.Sheet, s.Region, s.Last.X, s.Last.Y, p.X, p.Y, s.Color)
}

// Picker sets the primary (left button) or secondary (right button) color from the canvas.
type Picker struct{}

func (Picker) Icon() *gfx.Bitmap { return iconPicker }

func (Picker) Press(s *Stroke, p image.Point) {
	c := s.State.Sheet.At(p.X, p.Y)
	if s.Secondary {
		s.State.Secondary = c
	} else {
		s.State.Primary = c
	}
}

func (t Picker) Drag(s *Stroke, p image.Point) { t.Press(s, p) }

type Fill struct{}

func (Fill) Icon() *gfx.Bitmap { return iconFill }

func (Fill) Press(s *Stroke, p image.Point) {
	sprite.FloodFill(s.State.Sheet, s.Region, p.X, p.Y, s.Color)
}

func (Fill) Drag(*Stroke, image.Point) {}

// Action is a one-shot operation on the edit region. Register new actions by appending to Actions.
type Action struct {
	Icon *gfx.Bitmap
	Run  func(s *sprite.Sheet, r image.Rectangle)
}

var Actions = []Action{
	{iconFlipH, sprite.FlipH},
	{iconFlipV, sprite.FlipV},
	{iconRotate, sprite.Rotate},
	{iconClear, func(s *sprite.Sheet, r image.Rectangle) { sprite.Fill(s, r, 0) }},
}
