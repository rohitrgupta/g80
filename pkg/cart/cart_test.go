package cart

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	c := New()
	c.Sprites.Set(0, 0, 3)
	c.Sprites.Set(15, 9, 15)
	c.Sprites.Set(127, 255, 7)
	c.Sprites.Set(127, SpriteTilesY*8-1, 9)
	c.Palette[1].R = 0x12

	var buf bytes.Buffer
	if err := c.Encode(&buf); err != nil {
		t.Fatal(err)
	}
	got, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.Palette != c.Palette {
		t.Errorf("palette mismatch")
	}
	if !bytes.Equal(got.Sprites.Read(got.Sprites.Bounds()), c.Sprites.Read(c.Sprites.Bounds())) {
		t.Errorf("sprite mismatch")
	}
}

func TestDecodeLegacySprites(t *testing.T) {
	data := "# g80 cart v1\n[sprites]\n511:" + strings.Repeat("7", 64) + "\n"
	got, err := Decode(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if got.Sprites.NumTiles() != 16*256 {
		t.Fatalf("got %d sprites, want 4096", got.Sprites.NumTiles())
	}
	if got.Sprites.At(127, 255) != 7 {
		t.Error("legacy sprite pixels were not preserved")
	}
	if got.Sprites.At(127, SpriteTilesY*8-1) != 0 {
		t.Error("new sheets should start empty")
	}
}
