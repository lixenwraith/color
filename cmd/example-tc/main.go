package main

import (
	"context"
	"fmt"
	stdcolor "image/color"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/lixenwraith/color"
)

const reset = "\x1b[0m"

// bg returns the SGR sequence for a 24-bit background.
func bg(c color.RGB) string { return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B) }

// fg returns the SGR sequence for a 24-bit foreground.
func fg(c color.RGB) string { return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B) }

// block renders n background-colored spaces as a swatch.
func block(c color.RGB, n int) string { return bg(c) + strings.Repeat(" ", n) + reset }

func showPalette() {
	fmt.Println("── Named palette ──")
	entries := []struct {
		name string
		c    color.RGB
	}{
		{"Black", color.Black}, {"White", color.White}, {"Red", color.Red},
		{"Orange", color.Orange}, {"Gold", color.Gold}, {"Yellow", color.Yellow},
		{"Lime", color.Lime}, {"ForestGreen", color.ForestGreen}, {"Teal", color.Teal},
		{"Cyan", color.Cyan}, {"RoyalBlue", color.RoyalBlue}, {"Blue", color.Blue},
		{"Magenta", color.Magenta}, {"HotPink", color.HotPink}, {"Coral", color.Coral},
		{"Silver", color.Silver},
	}
	for _, e := range entries {
		fmt.Printf("%s %-12s %s\n", block(e.c, 4), e.name, e.c.Hex())
	}
	fmt.Println()
}

func showBlendModes() {
	fmt.Println("── Blend modes (dst=NavyBlue, src=Orange, alpha=0.6) ──")
	dst, src := color.NavyBlue, color.Orange
	const a = 0.6
	fmt.Printf("%s dst   %s src\n", block(dst, 6), block(src, 6))
	modes := []struct {
		name string
		out  color.RGB
	}{
		{"Blend", color.Blend(dst, src, a)},
		{"Add", color.Add(dst, src, a)},
		{"Screen", color.Screen(dst, src, a)},
		{"Overlay", color.Overlay(dst, src, a)},
		{"SoftLight", color.SoftLight(dst, src, a)},
		{"Max", color.Max(dst, src, a)},
	}
	for _, m := range modes {
		fmt.Printf("%s %-10s %s\n", block(m.out, 6), m.name, m.out.Hex())
	}
	fmt.Println()
}

func showGradient() {
	fmt.Println("── Lerp gradient: DarkCrimson → Gold → Teal ──")
	stops := []color.RGB{color.DarkCrimson, color.Gold, color.Teal}
	const steps = 36
	var line strings.Builder
	for i := 0; i < steps; i++ {
		seg := (float64(i) / float64(steps-1)) * float64(len(stops)-1)
		idx := int(seg)
		if idx >= len(stops)-1 {
			idx = len(stops) - 2
		}
		line.WriteString(block(color.Lerp(stops[idx], stops[idx+1], seg-float64(idx)), 1))
	}
	fmt.Println(line.String())
	fmt.Println()
}

func showGrayscale() {
	fmt.Println("── Grayscale, Desaturate, Luma ──")
	base := color.Coral
	fmt.Printf("%s base       %s  luma=%d\n", block(base, 6), base.Hex(), color.Luma(base))
	fmt.Printf("%s grayscale  %s\n", block(color.Grayscale(base), 6), color.Grayscale(base).Hex())
	var line strings.Builder
	for i := 0; i <= 8; i++ {
		line.WriteString(block(color.Desaturate(base, float64(i)/8.0), 3))
	}
	fmt.Printf("%s desaturate t=0→1\n\n", line.String())
}

func showHex() {
	fmt.Println("── Hex parse / format ──")
	for _, s := range []string{"#ff8800", "0f0", "#1a1b26", "bad!"} {
		c, err := color.ParseHex(s)
		if err != nil {
			fmt.Printf("  %-8q → error: %v\n", s, err)
			continue
		}
		fmt.Printf("  %-8q → %s %s%s%s\n", s, block(c, 3), fg(c), c.Hex(), reset)
	}
	c := color.MustParseHex("#41c7c7")
	fmt.Printf("  MustParseHex(#41c7c7) → %s %s\n\n", block(c, 3), c.Hex())
}

func showNearest() {
	fmt.Println("── Nearest palette color (RedmeanDistance) ──")
	named := []struct {
		name string
		c    color.RGB
	}{
		{"Red", color.Red}, {"Orange", color.Orange}, {"Gold", color.Gold},
		{"Lime", color.Lime}, {"Teal", color.Teal}, {"Blue", color.Blue},
		{"Magenta", color.Magenta}, {"White", color.White}, {"Black", color.Black},
	}
	for _, hex := range []string{"#ff5522", "#118ab2", "#2b2b2b"} {
		target := color.MustParseHex(hex)
		best, bestD := named[0], color.RedmeanDistance(target, named[0].c)
		for _, n := range named[1:] {
			if d := color.RedmeanDistance(target, n.c); d < bestD {
				best, bestD = n, d
			}
		}
		fmt.Printf("  %s %s  ≈  %s %-8s %s\n", block(target, 3), hex, block(best.c, 3), best.name, best.c.Hex())
	}
	fmt.Println()
}

func showLerpFixed() {
	fmt.Println("── LerpFixed (Q16.16 fixed-point) ──")
	const shift = 16
	a, b := color.RoyalBlue, color.Gold
	var line strings.Builder
	for i := 0; i <= 16; i++ {
		t := int64(i) << (shift - 4) // i/16 expressed in Q16.16; spans [0, 1<<shift]
		line.WriteString(block(color.LerpFixed(a, b, t, shift), 3))
	}
	fmt.Printf("%s RoyalBlue → Gold\n\n", line.String())
}

func showImageBridge() {
	fmt.Println("── image/color bridge (From / RGBA) ──")
	// image/color values are alpha-premultiplied; From un-premultiplies and drops alpha.
	opaque := stdcolor.RGBA{R: 200, G: 100, B: 50, A: 255}
	c := color.From(opaque)
	fmt.Printf("  From(RGBA{200,100,50,255}) = %s %s\n", block(c, 3), c.Hex())

	half := stdcolor.RGBA{R: 100, G: 50, B: 25, A: 128} // 50%% alpha, premultiplied
	c2 := color.From(half)
	fmt.Printf("  From(50%% alpha)            = %s %s  (recovers base hue)\n", block(c2, 3), c2.Hex())

	var _ stdcolor.Color = color.Orange // RGB satisfies image/color.Color
	r, g, b, alpha := color.Orange.RGBA()
	fmt.Printf("  color.Orange.RGBA()        = (%d,%d,%d,%d)\n\n", r, g, b, alpha)
}

// showPulse animates a Blend between two colors on one line until ctx is cancelled.
func showPulse(ctx context.Context) {
	const width = 32
	fmt.Println("── Live pulse (Ctrl+C to exit) ──")
	base, glow := color.DodgerBlue, color.Gold
	ticker := time.NewTicker(60 * time.Millisecond)
	defer ticker.Stop()
	phase := 0.0
	for {
		select {
		case <-ctx.Done():
			fmt.Print("\r" + strings.Repeat(" ", width+10) + "\r")
			fmt.Println("interrupted — bye")
			return
		case <-ticker.C:
			phase += 0.15
			t := (math.Sin(phase) + 1) / 2 // 0..1
			fmt.Printf("\r%s t=%.2f", block(color.Blend(base, glow, t), width), t)
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	showPalette()
	showBlendModes()
	showGradient()
	showGrayscale()
	showHex()
	showNearest()
	showLerpFixed()
	showImageBridge()
	showPulse(ctx)
}
