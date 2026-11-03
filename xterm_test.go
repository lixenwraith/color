package color

import "testing"

// Each reference shade maps to its own index, so the constants, the table and
// the search agree; the named primaries land on the bright indices like xterm.
func TestRGBTo16MapsReferenceShadesToTheirIndex(t *testing.T) {
	for i, c := range ansi16 {
		if got := RGBTo16(c); got != uint8(i) {
			t.Errorf("shade %d %s: index %d", i, c.Hex(), got)
		}
	}
	for c, want := range map[RGB]uint8{Red: ANSIBrightRed, Black: ANSIBlack, White: ANSIBrightWhite, {170, 0, 0}: ANSIRed} {
		if got := RGBTo16(c); got != want {
			t.Errorf("%s: index %d, want %d", c.Hex(), got, want)
		}
	}
}

// A gray maps to the nearest of the palette's grays, the ramp's and the cube's:
// each of them maps to itself, and any gray lands no farther than the nearest
// one plus the 4 levels a LUT cell spans
func TestGraysMapToTheNearestPaletteGray(t *testing.T) {
	var grays []uint8
	for i := 16; i < 256; i++ {
		if p := Xterm256(uint8(i)); p.R == p.G && p.G == p.B {
			grays = append(grays, uint8(i))
		}
	}
	for _, i := range grays {
		if got := RGBTo256(Xterm256(i)); got != i {
			t.Errorf("palette gray %d (%s) maps to %d", i, Xterm256(i).Hex(), got)
		}
	}
	for v := range 256 {
		got := Xterm256(RGBTo256(RGB{uint8(v), uint8(v), uint8(v)}))
		near := 255
		for _, i := range grays {
			near = min(near, max(int(Xterm256(i).R)-v, v-int(Xterm256(i).R)))
		}
		if off := max(int(got.R)-v, v-int(got.R)); got.R != got.G || got.G != got.B || off > near+4 {
			t.Errorf("gray %d maps to %s, where one %d levels off exists", v, got.Hex(), near)
		}
	}
}

// Palette indices show xterm's default colors, and each converts to the ANSI
// 16 index of its color
func TestXterm256ColorsAndTheirANSI16(t *testing.T) {
	for i, want := range map[uint8]RGB{1: {205, 0, 0}, 16: {0, 0, 0}, 21: {0, 0, 255}, 67: {95, 135, 175},
		196: {255, 0, 0}, 231: {255, 255, 255}, 232: {8, 8, 8}, 255: {238, 238, 238}} {
		if got := Xterm256(i); got != want {
			t.Errorf("index %d: %s, want %s", i, got.Hex(), want.Hex())
		}
	}
	for i := range 256 {
		if got, want := Xterm256To16(uint8(i)), RGBTo16(Xterm256(uint8(i))); got != want {
			t.Errorf("index %d (%s): ANSI %d, want %d", i, Xterm256(uint8(i)).Hex(), got, want)
		}
	}
}
