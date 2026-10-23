package sprite

import "image"

// TileSize is the width and height of a single sprite tile in pixels.
const TileSize = 8

// Sheet stores sprite pixels as 4-bit palette indices laid out in a grid of tiles.
type Sheet struct {
	tilesX, tilesY int
	pix            []uint8
	version        uint64
}

func NewSheet(tilesX, tilesY int) *Sheet {
	return &Sheet{
		tilesX: tilesX,
		tilesY: tilesY,
		pix:    make([]uint8, tilesX*tilesY*TileSize*TileSize),
	}
}

func (s *Sheet) Width() int              { return s.tilesX * TileSize }
func (s *Sheet) Height() int             { return s.tilesY * TileSize }
func (s *Sheet) TilesX() int             { return s.tilesX }
func (s *Sheet) TilesY() int             { return s.tilesY }
func (s *Sheet) NumTiles() int           { return s.tilesX * s.tilesY }
func (s *Sheet) Bounds() image.Rectangle { return image.Rect(0, 0, s.Width(), s.Height()) }

// Version changes whenever a pixel changes; used to invalidate cached renders.
func (s *Sheet) Version() uint64 { return s.version }

func (s *Sheet) At(x, y int) uint8 {
	if !image.Pt(x, y).In(s.Bounds()) {
		return 0
	}
	return s.pix[y*s.Width()+x]
}

func (s *Sheet) Set(x, y int, c uint8) {
	if !image.Pt(x, y).In(s.Bounds()) {
		return
	}
	c &= 0x0f
	i := y*s.Width() + x
	if s.pix[i] != c {
		s.pix[i] = c
		s.version++
	}
}

// TileRect returns the pixel rectangle of the tile with the given index.
func (s *Sheet) TileRect(index int) image.Rectangle {
	x, y := (index%s.tilesX)*TileSize, (index/s.tilesX)*TileSize
	return image.Rect(x, y, x+TileSize, y+TileSize)
}

// Read copies the pixels of r in row-major order.
func (s *Sheet) Read(r image.Rectangle) []uint8 {
	out := make([]uint8, 0, r.Dx()*r.Dy())
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			out = append(out, s.At(x, y))
		}
	}
	return out
}

// Write stores data (row-major, as returned by Read) into r.
func (s *Sheet) Write(r image.Rectangle, data []uint8) {
	i := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s.Set(x, y, data[i])
			i++
		}
	}
}
