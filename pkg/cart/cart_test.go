package cart

import (
	"bytes"
	"testing"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	c := New()
	c.Sprites.Set(0, 0, 3)
	c.Sprites.Set(15, 9, 15)
	c.Sprites.Set(127, 255, 7)
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
