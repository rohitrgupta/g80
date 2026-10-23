package sprite

import (
	"g80/pkg/gfx"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderer caches a Sheet as an ebiten image, refreshing it when the sheet or palette changes.
type Renderer struct {
	sheet   *Sheet
	palette *gfx.Palette
	img     *ebiten.Image
	buf     []byte
	version uint64
	lastPal gfx.Palette
}

func NewRenderer(s *Sheet, p *gfx.Palette) *Renderer {
	return &Renderer{sheet: s, palette: p}
}

func (r *Renderer) Image() *ebiten.Image {
	if r.img != nil && r.version == r.sheet.Version() && r.lastPal == *r.palette {
		return r.img
	}
	if r.img == nil {
		r.img = ebiten.NewImage(r.sheet.Width(), r.sheet.Height())
		r.buf = make([]byte, len(r.sheet.pix)*4)
	}
	for i, v := range r.sheet.pix {
		c := r.palette[v]
		r.buf[i*4], r.buf[i*4+1], r.buf[i*4+2], r.buf[i*4+3] = c.R, c.G, c.B, c.A
	}
	r.img.WritePixels(r.buf)
	r.version = r.sheet.Version()
	r.lastPal = *r.palette
	return r.img
}
