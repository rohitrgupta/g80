package gfx

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// FillRect fills r on dst, clipped to dst's bounds.
func FillRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	r = r.Intersect(dst.Bounds())
	if r.Empty() {
		return
	}
	dst.SubImage(r).(*ebiten.Image).Fill(c)
}

// StrokeRect draws a 1px outline just inside r.
func StrokeRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	FillRect(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+1), c)
	FillRect(dst, image.Rect(r.Min.X, r.Max.Y-1, r.Max.X, r.Max.Y), c)
	FillRect(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+1, r.Max.Y), c)
	FillRect(dst, image.Rect(r.Max.X-1, r.Min.Y, r.Max.X, r.Max.Y), c)
}

// DottedRect draws a 1px outline just inside r alternating between colors a and b.
func DottedRect(dst *ebiten.Image, r image.Rectangle, a, b color.Color) {
	dot := func(x, y int) {
		c := a
		if (x+y)%2 != 0 {
			c = b
		}
		FillRect(dst, image.Rect(x, y, x+1, y+1), c)
	}
	for x := r.Min.X; x < r.Max.X; x++ {
		dot(x, r.Min.Y)
		dot(x, r.Max.Y-1)
	}
	for y := r.Min.Y + 1; y < r.Max.Y-1; y++ {
		dot(r.Min.X, y)
		dot(r.Max.X-1, y)
	}
}
