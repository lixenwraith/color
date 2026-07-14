package color

import (
	"errors"
	stdcolor "image/color"
)

// RGB represents a 24-bit color
type RGB struct {
	R uint8 `toml:"r"`
	G uint8 `toml:"g"`
	B uint8 `toml:"b"`
}

var _ stdcolor.Color = RGB{}

// RGBA implements image/color.Color. Alpha is always opaque, so the
// premultiplied result equals the non-premultiplied one.
func (c RGB) RGBA() (r, g, b, a uint32) {
	return uint32(c.R) * 0x101, uint32(c.G) * 0x101, uint32(c.B) * 0x101, 0xffff
}

// From converts any image/color.Color to RGB, un-premultiplying by alpha and
// discarding it. Fully transparent input yields Black.
func From(c stdcolor.Color) RGB {
	r, g, b, a := c.RGBA()
	switch a {
	case 0:
		return RGB{}
	case 0xffff:
		return RGB{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8)}
	}
	// r,g,b <= a <= 0xffff, so r*0xffff <= 0xfffe0001 and stays in uint32
	return RGB{
		R: uint8(r * 0xffff / a >> 8),
		G: uint8(g * 0xffff / a >> 8),
		B: uint8(b * 0xffff / a >> 8),
	}
}

// Lerp linearly interpolates from c to other by t, clamped to [0,1]
func (c RGB) Lerp(other RGB, t float64) RGB {
	if t <= 0 {
		return c
	}
	if t >= 1 {
		return other
	}
	return RGB{
		R: clampU8(float64(c.R) + (float64(other.R)-float64(c.R))*t),
		G: clampU8(float64(c.G) + (float64(other.G)-float64(c.G))*t),
		B: clampU8(float64(c.B) + (float64(other.B)-float64(c.B))*t),
	}
}

// Luma returns Rec. 601 luminance: R*0.299 + G*0.587 + B*0.114
func Luma(c RGB) uint8 {
	return uint8((int(c.R)*299 + int(c.G)*587 + int(c.B)*114) / 1000)
}

// RedmeanDistance returns squared perceptually-weighted distance between a and b.
// Monotonic in perceived difference; use for nearest-color search, not as an
// absolute metric.
// Formula: https://en.wikipedia.org/wiki/Color_difference#sRGB
func RedmeanDistance(a, b RGB) int {
	rmean := (int(a.R) + int(b.R)) / 2
	dr := int(a.R) - int(b.R)
	dg := int(a.G) - int(b.G)
	db := int(a.B) - int(b.B)
	return (((512 + rmean) * dr * dr) >> 8) + 4*dg*dg + (((767 - rmean) * db * db) >> 8)
}

// Hex returns the color as "#rrggbb"
func (c RGB) Hex() string {
	const d = "0123456789abcdef"
	b := [7]byte{
		'#',
		d[c.R>>4], d[c.R&0xf],
		d[c.G>>4], d[c.G&0xf],
		d[c.B>>4], d[c.B&0xf],
	}
	return string(b[:])
}

var errHex = errors.New("color: invalid hex string")

// ParseHex accepts "#rgb", "#rrggbb", and the same forms without the leading '#'
func ParseHex(s string) (RGB, error) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	switch len(s) {
	case 3:
		r, ok0 := nibble(s[0])
		g, ok1 := nibble(s[1])
		b, ok2 := nibble(s[2])
		if !ok0 || !ok1 || !ok2 {
			return RGB{}, errHex
		}
		return RGB{R: r * 0x11, G: g * 0x11, B: b * 0x11}, nil
	case 6:
		var v [3]uint8
		for i := range v {
			hi, ok0 := nibble(s[i*2])
			lo, ok1 := nibble(s[i*2+1])
			if !ok0 || !ok1 {
				return RGB{}, errHex
			}
			v[i] = hi<<4 | lo
		}
		return RGB{R: v[0], G: v[1], B: v[2]}, nil
	}
	return RGB{}, errHex
}

// MustParseHex panics on invalid input; for package-level initializers
func MustParseHex(s string) RGB {
	c, err := ParseHex(s)
	if err != nil {
		panic(err)
	}
	return c
}

func nibble(b byte) (uint8, bool) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', true
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, true
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, true
	}
	return 0, false
}
