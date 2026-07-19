package main

import (
	"fmt"
	"strings"

	"github.com/lixenwraith/color"
)

const reset = "\x1b[0m"

// bg256 returns the SGR sequence for an 8-bit xterm-256 background.
// This is the sequence used by naked terminals (TTY, basic SSH).
func bg256(idx uint8) string { return fmt.Sprintf("\x1b[48;5;%dm", idx) }

// bgRGB returns the SGR sequence for a 24-bit truecolor background.
func bgRGB(c color.RGB) string { return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B) }

// block256 renders n background-colored spaces using the 256-color palette.
func block256(idx uint8, n int) string { return bg256(idx) + strings.Repeat(" ", n) + reset }

// blockRGB renders n background-colored spaces using truecolor.
func blockRGB(c color.RGB, n int) string { return bgRGB(c) + strings.Repeat(" ", n) + reset }

// =====================================================================
// PART 1: 256 Color Mechanics
// =====================================================================

func show256Cube() {
	fmt.Println("── Part 1: xterm-256 Color Cube (6x6x6) ──")
	fmt.Println("Displaying the 216-color cube layout (Indices 16-231).")

	// The standard xterm-256 cube consists of 6 "slices" of red,
	// containing 6x6 grids of green and blue.
	for r := uint8(0); r < 6; r++ {
		for g := uint8(0); g < 6; g++ {
			for b := uint8(0); b < 6; b++ {
				// Get exact index using the pure math function
				idx := color.Cube256(r, g, b)
				fmt.Print(block256(idx, 3))
			}
			fmt.Print("  ") // Space between green columns
		}
		fmt.Println()
	}
	fmt.Println()
}

func showGrayscaleRamp() {
	fmt.Println("── Part 1: Grayscale Ramp ──")
	fmt.Println("Displaying the 24-step grayscale ramp (Indices 232-255).")

	for step := uint8(0); step < 24; step++ {
		idx := color.Gray256(step)
		fmt.Print(block256(idx, 2))
	}
	fmt.Println("\n")
}

// =====================================================================
// PART 2: Naked Terminal (No Desktop Environment)
// =====================================================================

func showNakedTerminalDegradation() {
	fmt.Println("── Part 2: Naked Terminal (Quantization Fallback) ──")
	fmt.Println("Simulating how 24-bit Truecolor gracefully degrades in environments")
	fmt.Println("without a modern desktop compositor (e.g., bare TTY, tmux, legacy SSH).")
	fmt.Println()

	// 1. Force the lazy evaluation of the 262KB Redmean LUT.
	// This makes RGBTo256 O(1) during the render loop.
	color.WarmXterm256()

	// 2. Render a smooth gradient in both truecolor and degraded 256-color
	const steps = 40
	start, end := color.ElectricViolet, color.LimeGreen

	var truecolor strings.Builder
	var naked256 strings.Builder

	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		c := color.Lerp(start, end, t)

		// Desktop Environment (24-bit)
		truecolor.WriteString(blockRGB(c, 2))

		// Naked Terminal (8-bit Quantized via Redmean perceptual distance)
		idx := color.RGBTo256(c)
		naked256.WriteString(block256(idx, 2))
	}

	fmt.Println("Desktop Environment (24-bit smooth interpolation):")
	fmt.Println(truecolor.String())
	fmt.Println("Naked Terminal (8-bit perceptual quantization via RGBTo256):")
	fmt.Println(naked256.String())
	fmt.Println()

	// 3. Show exact index mapping for specific named colors
	fmt.Println("Nearest xterm-256 mappings for Named RGB colors:")
	named := []struct {
		name string
		c    color.RGB
	}{
		{"HotPink", color.HotPink},
		{"BurntOrange", color.BurntOrange},
		{"MintGreen", color.MintGreen},
		{"DeepNavy", color.DeepNavy},
		{"SlateGray", color.SlateGray},
	}

	for _, n := range named {
		// Calculate nearest palette index
		idx := color.RGBTo256(n.c)

		// Map index back to mathematical coordinates to see where it landed
		cubeStr := ""
		if idx >= 16 && idx <= 231 {
			r, g, b := color.CubeRGB256(idx)
			cubeStr = fmt.Sprintf("Cube(r:%d, g:%d, b:%d)", r, g, b)
		} else if idx >= 232 {
			cubeStr = fmt.Sprintf("GrayStep(%d)", idx-232)
		} else {
			cubeStr = "System(0-15)"
		}

		fmt.Printf("  %-12s %s (Truecolor) -> %s (Index %3d) %s\n",
			n.name,
			blockRGB(n.c, 4),
			block256(idx, 4),
			idx,
			cubeStr,
		)
	}
	fmt.Println()
}

func main() {
	// Execute Part 1
	show256Cube()
	showGrayscaleRamp()

	// Execute Part 2
	showNakedTerminalDegradation()
}
