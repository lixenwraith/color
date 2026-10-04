# color

24-bit RGB values, perceptual metrics, and blend operations. No output device, no dependencies.

Extracted from [lixenwraith/terminal](https://github.com/lixenwraith/terminal) so renderers — terminal, GUI, image, framebuffer — share one color type and one set of operations without linking terminal I/O, `x/sys`, or termios.

## Install

```
go get github.com/lixenwraith/color
```

Go 1.26+. Standard library only.

## Type

```go
type RGB struct {
    R uint8 `toml:"r"`
    G uint8 `toml:"g"`
    B uint8 `toml:"b"`
}
```

Three bytes, comparable, pointer-free: safe to embed in dense cell or pixel buffers and to bind directly from config.

`RGB` implements `image/color.Color`, so values pass into `image`, `draw`, and GUI toolkit pipelines unchanged. `From` converts back, un-premultiplying alpha and discarding it.

```go
img.Set(x, y, color.Amber)          // RGB satisfies image/color.Color
c := color.From(img.At(x, y))       // back to RGB
```

Where both packages are needed at one site, alias the standard library: `import stdcolor "image/color"`.

## Operations

| Call | Behavior |
| --- | --- |
| `Blend(dst, src, alpha)` | Linear alpha compositing |
| `SoftLight(dst, src, intensity)` | Perez soft light; gentler than linear alpha |
| `Overlay(dst, src, alpha)` | Multiply on darks, screen on lights |
| `Screen(dst, src, alpha)` | Always lightens; glow accumulation without `Add` clipping |
| `Add(dst, src, alpha)` | Saturating additive |
| `Max(dst, src, alpha)` | Per-channel maximum |
| `Scale(c, factor)` | Channel multiply, saturating |
| `Grayscale(c)`, `Luma(c)` | Rec. 601 luma |
| `c.Lerp(other, t)` | Linear interpolation |
| `RedmeanDistance(a, b)` | Squared perceptual distance, for nearest-color search |

All operations are pure. `alpha` and `t` clamp to `[0,1]`; channels saturate. Integer paths avoid division; soft light is table-driven, no `sqrt` per channel.

```go
bg   := color.Obsidian
glow := color.Screen(bg, color.Amber, 0.4)
edge := bg.Lerp(color.Amber, 0.75)
warm := color.SoftLight(edge, color.Terracotta, 0.3)
```

## Palette

~120 named colors, grouped by hue and ordered dark-to-light: `Obsidian`, `Amber`, `EmeraldGreen`, `LightSkyBlue`, `Vermilion`, … Standard names (CSS, X11) where the RGB matches; descriptive compounds otherwise.

Package-level `var`s. **Read-only by contract** — the language permits assignment, the package does not. Alias them into domain parameter files rather than mutating them.

## Hex

```go
c, err := color.ParseHex("#4a90d9")     // also "4a90d9", "#abc", "abc"
s := c.Hex()                            // "#4a90d9"
var Accent = color.MustParseHex("#ff8800")
```

`RGB` deliberately does **not** implement `encoding.TextUnmarshaler`. TOML and JSON decoders prefer it over struct-field unification, which would silently break table-form config:

```toml
accent = { r = 255, g = 136, b = 0 }
```

## Terminal palettes

Quantization for terminals that cannot take 24-bit color:

```go
idx := color.RGBTo256(color.EmeraldGreen)   // xterm-256 index, Redmean LUT
ansi := color.RGBTo16(color.EmeraldGreen)   // ANSI 16 index, for text consoles
```

`ANSIBlack` .. `ANSIBrightWhite` name the 16 indices (SGR `30-37`/`90-97`).
The terminal's theme picks their shade, so they suit semantic colors such as
log levels. `RGBTo16` measures against xterm's default shades.

`terminal.Cell` carries `color.RGB` directly. When `Cell.Attrs` sets
`AttrFg256` / `AttrBg256`, `Cell.Fg.R` / `Cell.Bg.R` hold a palette index, not
a channel. Such values are not colors and must not be passed to this package.

## Concurrency

Values are immutable, operations are pure, lookup tables are built at package init. Safe for concurrent use.

## License

See `LICENSE`.
