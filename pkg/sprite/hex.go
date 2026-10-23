package sprite

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const hexDigits = "0123456789abcdef"

// EncodeHex writes every non-empty tile as one line: "<index>:<64 hex digits>",
// one digit per pixel in row-major order.
func (s *Sheet) EncodeHex(w io.Writer) error {
	buf := make([]byte, TileSize*TileSize)
	for i := 0; i < s.NumTiles(); i++ {
		empty := true
		for j, v := range s.Read(s.TileRect(i)) {
			buf[j] = hexDigits[v]
			empty = empty && v == 0
		}
		if empty {
			continue
		}
		if _, err := fmt.Fprintf(w, "%03d:%s\n", i, buf); err != nil {
			return err
		}
	}
	return nil
}

// DecodeHex reads tiles written by EncodeHex. Tiles not listed are left untouched.
func (s *Sheet) DecodeHex(lines []string) error {
	for n, line := range lines {
		idx, data, ok := strings.Cut(line, ":")
		if !ok {
			return fmt.Errorf("line %d: missing ':'", n+1)
		}
		i, err := strconv.Atoi(strings.TrimSpace(idx))
		if err != nil || i < 0 || i >= s.NumTiles() {
			return fmt.Errorf("line %d: invalid tile index %q", n+1, idx)
		}
		data = strings.TrimSpace(data)
		if len(data) != TileSize*TileSize {
			return fmt.Errorf("line %d: expected %d hex digits, got %d", n+1, TileSize*TileSize, len(data))
		}
		pix := make([]uint8, len(data))
		for j := 0; j < len(data); j++ {
			v := strings.IndexByte(hexDigits, toLower(data[j]))
			if v < 0 {
				return fmt.Errorf("line %d: invalid hex digit %q", n+1, data[j])
			}
			pix[j] = uint8(v)
		}
		s.Write(s.TileRect(i), pix)
	}
	return nil
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}
