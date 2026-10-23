package gfx

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strings"
)

// PaletteSize is the number of colors available to sprites.
const PaletteSize = 16

// Palette maps a 4-bit color index to an RGB color.
type Palette [PaletteSize]color.RGBA

// Sweetie16 is the default palette (same as TIC-80).
var Sweetie16 = MustParsePalette("1a1c2c5d275db13e53ef7d57ffcd75a7f07038b76425717929366f3b5dc941a6f673eff7f4f4f494b0c2566c86333c57")

// ParsePalette parses 16 colors written as a single string of RRGGBB hex values.
func ParsePalette(s string) (Palette, error) {
	var p Palette
	b, err := hex.DecodeString(strings.TrimSpace(s))
	if err != nil {
		return p, fmt.Errorf("palette: %w", err)
	}
	if len(b) != PaletteSize*3 {
		return p, fmt.Errorf("palette: expected %d colors, got %d bytes", PaletteSize, len(b))
	}
	for i := range p {
		p[i] = color.RGBA{b[i*3], b[i*3+1], b[i*3+2], 0xff}
	}
	return p, nil
}

func MustParsePalette(s string) Palette {
	p, err := ParsePalette(s)
	if err != nil {
		panic(err)
	}
	return p
}

// String returns the palette in the format accepted by ParsePalette.
func (p Palette) String() string {
	b := make([]byte, 0, PaletteSize*3)
	for _, c := range p {
		b = append(b, c.R, c.G, c.B)
	}
	return hex.EncodeToString(b)
}
