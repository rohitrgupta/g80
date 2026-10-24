package cart

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"g80/pkg/gfx"
	"g80/pkg/sprite"
	"g80/pkg/tilemap"
)

// Sprite sheets are stored vertically, with 256 tiles per sheet.
const (
	SpriteSheetCount = 16
	SpriteTilesX     = 16
	SpriteTilesY     = SpriteTilesX * SpriteSheetCount
)

// Cart holds all game data that is saved to disk.
type Cart struct {
	Palette gfx.Palette
	Sprites *sprite.Sheet
	Maps    [tilemap.Count]*tilemap.Map
}

func New() *Cart {
	c := &Cart{
		Palette: gfx.Sweetie16,
		Sprites: sprite.NewSheet(SpriteTilesX, SpriteTilesY),
	}
	for index := range c.Maps {
		c.Maps[index] = tilemap.New()
	}
	return c
}

// section is one "[name]" block of the plain-text cart file.
// Add entries to sections to persist new kinds of data (maps, code, sound...).
type section struct {
	name   string
	encode func(c *Cart, w io.Writer) error
	decode func(c *Cart, lines []string) error
}

var sections = []section{
	{
		name: "palette",
		encode: func(c *Cart, w io.Writer) error {
			_, err := fmt.Fprintln(w, c.Palette.String())
			return err
		},
		decode: func(c *Cart, lines []string) error {
			p, err := gfx.ParsePalette(strings.Join(lines, ""))
			c.Palette = p
			return err
		},
	},
	{
		name:   "sprites",
		encode: func(c *Cart, w io.Writer) error { return c.Sprites.EncodeHex(w) },
		decode: func(c *Cart, lines []string) error { return c.Sprites.DecodeHex(lines) },
	},
	{
		name: "maps",
		encode: func(c *Cart, w io.Writer) error {
			maps := make(map[int]*tilemap.Map)
			for index, m := range c.Maps {
				if m == nil {
					return fmt.Errorf("map %d is null", index)
				}
				if !m.IsEmpty() {
					maps[index] = m
				}
			}
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", " ")
			return encoder.Encode(maps)
		},
		decode: func(c *Cart, lines []string) error {
			var maps map[int]*tilemap.Map
			if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &maps); err != nil {
				return err
			}
			if maps == nil {
				return fmt.Errorf("expected an indexed maps object")
			}
			for index, m := range maps {
				if index < 0 || index >= tilemap.Count {
					return fmt.Errorf("invalid map index %d", index)
				}
				if m == nil {
					return fmt.Errorf("map %d is null", index)
				}
				if err := m.Validate(c.Sprites.NumTiles()); err != nil {
					return fmt.Errorf("map %d: %w", index, err)
				}
				c.Maps[index] = m
			}
			return nil
		},
	},
}

const header = "# g80 cart v1"

func (c *Cart) Encode(w io.Writer) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintln(bw, header)
	for _, s := range sections {
		fmt.Fprintf(bw, "\n[%s]\n", s.name)
		if err := s.encode(c, bw); err != nil {
			return fmt.Errorf("[%s]: %w", s.name, err)
		}
	}
	return bw.Flush()
}

// Decode parses a cart. Blank lines and lines starting with '#' are ignored,
// as are unknown sections.
func Decode(r io.Reader) (*Cart, error) {
	chunks := map[string][]string{}
	current := ""
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
		case len(line) > 2 && strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") && !strings.ContainsAny(line[1:len(line)-1], "[]{}\","):
			current = strings.TrimSpace(line[1 : len(line)-1])
		case current == "":
			return nil, fmt.Errorf("line %d: data outside of a section", n)
		default:
			chunks[current] = append(chunks[current], line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	c := New()
	for _, s := range sections {
		if lines, ok := chunks[s.name]; ok {
			if err := s.decode(c, lines); err != nil {
				return nil, fmt.Errorf("[%s]: %w", s.name, err)
			}
		}
	}
	return c, nil
}

func Load(path string) (*Cart, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

func (c *Cart) Save(path string) error {
	var buf bytes.Buffer
	if err := c.Encode(&buf); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
