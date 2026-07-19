package color

import (
	"sync"
	"sync/atomic"
)

// Generic xterm 256-color palette indices
// Color cube: index = 16 + 36*r + 6*g + b where r,g,b ∈ [0,5]
// Grayscale ramp: indices 232-255, level = 8 + 10*(index-232)
// Ordered dark-to-light within each hue group

const (
	P256DeepNavy  uint8 = 17 // (0,0,1)
	P256DarkBlue  uint8 = 18 // (0,0,2)
	P256SteelBlue uint8 = 75 // (1,3,5)
	P256LightBlue uint8 = 81 // (1,4,5)

	P256DeepTeal  uint8 = 23 // (0,1,1)
	P256Teal      uint8 = 44 // (0,4,4)
	P256Green     uint8 = 46 // (0,5,0)
	P256Cyan      uint8 = 51 // (0,5,5)
	P256LightCyan uint8 = 87 // (1,5,5)

	P256CobaltBlue     uint8 = 33  // (0,2,5)
	P256DarkPurpleBlue uint8 = 54  // (1,0,2)
	P256Indigo         uint8 = 63  // (1,1,5)
	P256Purple         uint8 = 129 // (3,0,5)
	P256Violet         uint8 = 134 // (3,1,4)
	P256MediumPurple   uint8 = 135 // (3,1,5)
	P256Orchid         uint8 = 176 // (4,2,4)

	P256YellowGreen uint8 = 154 // (3,5,0)

	P256Maroon      uint8 = 52  // (1,0,0)
	P256DarkCrimson uint8 = 88  // (2,0,0)
	P256Crimson     uint8 = 160 // (4,0,0)

	P256Red       uint8 = 196 // (5,0,0)
	P256Rose      uint8 = 198 // (5,0,2)
	P256RedOrange uint8 = 202 // (5,1,0)
	P256Orange    uint8 = 208 // (5,2,0)
	P256Amber     uint8 = 214 // (5,3,0)
	P256Gold      uint8 = 220 // (5,4,0)
	P256Yellow    uint8 = 226 // (5,5,0)

	P256DarkAmber uint8 = 94 // (2,1,0)

	P256Gray uint8 = 240 // Grayscale step 8, level ~88
)

// 6-bit quantized LUT for Redmean-based 256-color mapping
// 64×64×64 = 262,144 bytes, fits in L2 cache
const lut256Size = 64 * 64 * 64

var (
	lut256Ptr  atomic.Pointer[[lut256Size]uint8]
	lut256Once sync.Once
)

// lut256 returns the palette LUT, building it on first use.
func lut256() *[lut256Size]uint8 {
	if p := lut256Ptr.Load(); p != nil {
		return p
	}
	return lut256Build()
}

//go:noinline
func lut256Build() *[lut256Size]uint8 {
	lut256Once.Do(func() {
		t := new([lut256Size]uint8)
		for r := range 64 {
			for g := range 64 {
				for b := range 64 {
					c := RGB{
						R: uint8(r<<2 | 2),
						G: uint8(g<<2 | 2),
						B: uint8(b<<2 | 2),
					}
					t[r<<12|g<<6|b] = computeRedmean256(c)
				}
			}
		}
		lut256Ptr.Store(t)
	})
	return lut256Ptr.Load()
}

// WarmXterm256 forces LUT construction. Idempotent, safe for concurrent use.
// Optional, prevents latency spike on first RGBTo256 render in the terminal loop.
func WarmXterm256() { _ = lut256() }

func computeRedmean256(c RGB) uint8 {
	if c.R == c.G && c.G == c.B {
		if c.R < 8 {
			return 16
		}
		if c.R > 238 {
			return 231
		}
		return uint8(232 + (int(c.R)-8)/10)
	}

	bestIdx := uint8(16)
	minDist := 1 << 30

	for i := range 216 {
		cand := RGB{
			R: cubeValues[i/36],
			G: cubeValues[(i/6)%6],
			B: cubeValues[i%6],
		}
		if d := RedmeanDistance(c, cand); d < minDist {
			minDist = d
			bestIdx = uint8(16 + i)
		}
	}

	for i := range 24 {
		g := uint8(8 + i*10)
		if d := RedmeanDistance(c, RGB{R: g, G: g, B: g}); d < minDist {
			minDist = d
			bestIdx = uint8(232 + i)
		}
	}

	return bestIdx
}

var cubeValues = [6]uint8{0, 95, 135, 175, 215, 255}

// RGBTo256 converts 24-bit RGB to the nearest xterm-256 palette index.
// O(1) via the Redmean LUT; the first call builds it.
func RGBTo256(c RGB) uint8 {
	return lut256()[int(c.R>>2)<<12|int(c.G>>2)<<6|int(c.B>>2)]
}

// Cube256 returns the xterm 256-palette index for an RGB cube coordinate.
func Cube256(r, g, b uint8) uint8 {
	if r > 5 {
		r = 5
	}
	if g > 5 {
		g = 5
	}
	if b > 5 {
		b = 5
	}
	return 16 + 36*r + 6*g + b
}

// CubeRGB256 returns the (r, g, b) cube coordinates for a 256-palette color cube index.
func CubeRGB256(index uint8) (r, g, b uint8) {
	if index < 16 || index > 231 {
		return 0, 0, 0
	}
	n := index - 16
	r = n / 36
	g = (n % 36) / 6
	b = n % 6
	return r, g, b
}

// Gray256 returns the xterm 256-palette index for a grayscale step.
func Gray256(step uint8) uint8 {
	if step > 23 {
		step = 23
	}
	return 232 + step
}
