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
		if p := palette256(uint8(i)); p.R == p.G && p.G == p.B {
			grays = append(grays, uint8(i))
		}
	}
	for _, i := range grays {
		if got := RGBTo256(palette256(i)); got != i {
			t.Errorf("palette gray %d (%s) maps to %d", i, palette256(i).Hex(), got)
		}
	}
	for v := range 256 {
		got := palette256(RGBTo256(RGB{uint8(v), uint8(v), uint8(v)}))
		near := 255
		for _, i := range grays {
			near = min(near, max(int(palette256(i).R)-v, v-int(palette256(i).R)))
		}
		if off := max(int(got.R)-v, v-int(got.R)); got.R != got.G || got.G != got.B || off > near+4 {
			t.Errorf("gray %d maps to %s, where one %d levels off exists", v, got.Hex(), near)
		}
	}
}

// palette256 is the color of cube or ramp index i
func palette256(i uint8) RGB {
	if i >= 232 {
		v := 8 + 10*(i-232)
		return RGB{v, v, v}
	}
	r, g, b := CubeRGB256(i)
	return RGB{cubeValues[r], cubeValues[g], cubeValues[b]}
}
