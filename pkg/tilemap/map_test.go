package tilemap

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	m := New()
	if m.At(0, 0) != Empty || m.At(-1, 0) != Empty {
		t.Fatal("new map must be empty")
	}
	m.Set(31, 15, 31)
	m.Set(32, 15, 2)
	m.Set(31, 15, 32)
	if m.At(31, 15) != 31 {
		t.Fatal("valid cell was not preserved")
	}
	m.Sprites[31] = 4095
	if err := m.Validate(4096); err != nil {
		t.Fatal(err)
	}
	m.Sprites[31] = 4096
	if m.Validate(4096) == nil {
		t.Fatal("invalid sprite accepted")
	}
}

func TestMapASCIIEncoding(t *testing.T) {
	m := New()
	for slot := 0; slot < PaletteSize; slot++ {
		m.Set(slot, 0, uint8(slot))
	}
	m.Sprites[31] = 4095
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var encoded struct {
		Tiles string `json:"tiles"`
	}
	if err := json.Unmarshal(data, &encoded); err != nil {
		t.Fatal(err)
	}
	want := "0123456789abcdefghijklmnopqrstuv" + strings.Repeat(".", Width*(Height-1))
	if encoded.Tiles != want {
		t.Fatalf("incorrect ASCII map: %q", encoded.Tiles)
	}
	var got Map
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got != *m {
		t.Fatal("ASCII map did not round trip")
	}
}

func TestMapRejectsLegacyJSON(t *testing.T) {
	m := New()
	m.Set(31, 15, 31)
	m.Sprites[31] = 4095
	type legacyMap Map
	data, err := json.Marshal((*legacyMap)(m))
	if err != nil {
		t.Fatal(err)
	}
	var got Map
	if err := json.Unmarshal(data, &got); err == nil {
		t.Fatal("legacy numeric map was accepted")
	}
}

func TestMapRejectsInvalidJSONTiles(t *testing.T) {
	for _, tiles := range []string{
		`""`,
		`"0123"`,
		`"` + strings.Repeat(".", Width*Height-1) + `w"`,
		`"` + strings.Repeat(".", Width*Height+1) + `"`,
		`[0,1]`,
		`[` + strings.Repeat("0,", Width*Height-1) + `32]`,
		`null`,
	} {
		m := New()
		before := *m
		if err := json.Unmarshal([]byte(`{"tiles":`+tiles+`}`), m); err == nil {
			t.Errorf("invalid tiles accepted: %s", tiles)
		}
		if *m != before {
			t.Error("failed decode changed the map")
		}
	}
}
