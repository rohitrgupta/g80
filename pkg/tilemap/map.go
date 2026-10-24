package tilemap

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	Count       = 64
	Width       = 32
	Height      = 16
	PaletteSize = 32
	Empty       = uint8(255)
)

type Map struct {
	Sprites [PaletteSize]uint16   `json:"sprites"`
	Tiles   [Width * Height]uint8 `json:"tiles"`
}

const tileSymbols = "0123456789abcdefghijklmnopqrstuv"

func (m Map) MarshalJSON() ([]byte, error) {
	var tiles [Width * Height]byte
	for cell, slot := range m.Tiles {
		switch {
		case slot == Empty:
			tiles[cell] = '.'
		case slot < PaletteSize:
			tiles[cell] = tileSymbols[slot]
		default:
			return nil, fmt.Errorf("invalid map palette slot %d", slot)
		}
	}
	return json.Marshal(struct {
		Sprites [PaletteSize]uint16 `json:"sprites"`
		Tiles   string              `json:"tiles"`
	}{Sprites: m.Sprites, Tiles: string(tiles[:])})
}

func (m *Map) UnmarshalJSON(data []byte) error {
	var encoded struct {
		Sprites [PaletteSize]uint16 `json:"sprites"`
		Tiles   string              `json:"tiles"`
	}
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}
	next := Map{Sprites: encoded.Sprites}
	if len(encoded.Tiles) != len(next.Tiles) {
		return fmt.Errorf("expected %d map cells, got %d", len(next.Tiles), len(encoded.Tiles))
	}
	for cell, symbol := range []byte(encoded.Tiles) {
		if symbol == '.' {
			next.Tiles[cell] = Empty
			continue
		}
		slot := strings.IndexByte(tileSymbols, symbol)
		if slot < 0 {
			return fmt.Errorf("cell %d: invalid map symbol %q", cell, symbol)
		}
		next.Tiles[cell] = uint8(slot)
	}
	*m = next
	return nil
}

func (m *Map) IsEmpty() bool {
	for _, slot := range m.Tiles {
		if slot != Empty {
			return false
		}
	}
	return true
}

func New() *Map {
	m := &Map{}
	for slot := range m.Sprites {
		m.Sprites[slot] = uint16(slot)
	}
	for cell := range m.Tiles {
		m.Tiles[cell] = Empty
	}
	return m
}

func (m *Map) At(x, y int) uint8 {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return Empty
	}
	return m.Tiles[y*Width+x]
}

func (m *Map) Set(x, y int, slot uint8) {
	if x >= 0 && x < Width && y >= 0 && y < Height && (slot < PaletteSize || slot == Empty) {
		m.Tiles[y*Width+x] = slot
	}
}

func (m *Map) Validate(spriteCount int) error {
	for _, index := range m.Sprites {
		if int(index) >= spriteCount {
			return fmt.Errorf("invalid sprite index %d", index)
		}
	}
	for _, slot := range m.Tiles {
		if slot >= PaletteSize && slot != Empty {
			return fmt.Errorf("invalid map palette slot %d", slot)
		}
	}
	return nil
}
