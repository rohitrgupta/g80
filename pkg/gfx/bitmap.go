package gfx

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Bitmap is a 1-bit image (icons, glyphs) that can be drawn in any color.
type Bitmap struct {
	w, h int
	pix  []bool
	img  *ebiten.Image
}

// NewBitmap builds a bitmap from text rows where '#' marks a set pixel.
func NewBitmap(rows ...string) *Bitmap {
	b := &Bitmap{h: len(rows)}
	for _, r := range rows {
		b.w = max(b.w, len(r))
	}
	b.pix = make([]bool, b.w*b.h)
	for y, r := range rows {
		for x := 0; x < len(r); x++ {
			b.pix[y*b.w+x] = r[x] == '#'
		}
	}
	return b
}

func (b *Bitmap) Size() (int, int) { return b.w, b.h }

// Transpose returns the bitmap mirrored along its main diagonal.
func (b *Bitmap) Transpose() *Bitmap {
	t := &Bitmap{w: b.h, h: b.w, pix: make([]bool, len(b.pix))}
	for y := 0; y < b.h; y++ {
		for x := 0; x < b.w; x++ {
			t.pix[x*t.w+y] = b.pix[y*b.w+x]
		}
	}
	return t
}

func (b *Bitmap) image() *ebiten.Image {
	if b.img == nil {
		buf := make([]byte, b.w*b.h*4)
		for i, on := range b.pix {
			if on {
				copy(buf[i*4:], []byte{0xff, 0xff, 0xff, 0xff})
			}
		}
		b.img = ebiten.NewImage(b.w, b.h)
		b.img.WritePixels(buf)
	}
	return b.img
}

// Draw draws the set pixels of b at (x, y) using color c.
func (b *Bitmap) Draw(dst *ebiten.Image, x, y int, c color.Color) {
	if b.w == 0 || b.h == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(c)
	dst.DrawImage(b.image(), op)
}
