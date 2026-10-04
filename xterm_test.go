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
