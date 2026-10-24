package cart

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"g80/pkg/tilemap"
)

func TestMapRoundTrip(t *testing.T) {
	c := New()
	c.Maps[63].Sprites[31] = 4095
	c.Maps[63].Set(31, 15, 31)
	var data bytes.Buffer
	if err := c.Encode(&data); err != nil {
		t.Fatal(err)
	}
	got, err := Decode(&data)
	if err != nil {
		t.Fatal(err)
	}
	if *got.Maps[63] != *c.Maps[63] {
		t.Fatal("map did not round trip")
	}
}

func TestCartWithoutMapsHasEmptyMaps(t *testing.T) {
	c, err := Decode(strings.NewReader("# g80 cart v1\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range c.Maps {
		if m == nil || m.At(0, 0) != tilemap.Empty {
			t.Fatal("cart without map data must have empty maps")
		}
	}
}

func TestSaveOmitsEmptyMaps(t *testing.T) {
	for _, populated := range []bool{false, true} {
		c := New()
		c.Maps[0].Sprites[0] = 4095
		if populated {
			c.Maps[63].Set(31, 15, 0)
		}
		var data bytes.Buffer
		if err := c.Encode(&data); err != nil {
			t.Fatal(err)
		}
		_, payload, ok := strings.Cut(data.String(), "[maps]\n")
		if !ok {
			t.Fatal("missing map section")
		}
		var maps map[int]json.RawMessage
		if err := json.Unmarshal([]byte(payload), &maps); err != nil {
			t.Fatal(err)
		}
		if populated {
			if len(maps) != 1 || maps[63] == nil {
				t.Fatal("only populated map 63 should be saved")
			}
		} else if len(maps) != 0 {
			t.Fatal("empty maps must not be saved")
		}
		got, err := Decode(&data)
		if err != nil {
			t.Fatal(err)
		}
		if *got.Maps[63] != *c.Maps[63] || !got.Maps[0].IsEmpty() {
			t.Fatal("sparse map indices were not preserved")
		}
	}
}

func TestRejectInvalidMapData(t *testing.T) {
	for _, payload := range []string{
		`[]`,
		`[{}]`,
		`null`,
		`{"-1":{}}`,
		`{"64":{}}`,
		`{"0":null}`,
		`{"0":{"tiles":[` + strings.Repeat("0,", tilemap.Width*tilemap.Height-1) + `0]}}`,
	} {
		if _, err := Decode(strings.NewReader("# g80 cart v1\n[maps]\n" + payload)); err == nil {
			t.Errorf("invalid map data accepted: %s", payload)
		}
	}
}
