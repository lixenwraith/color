// Package color provides 24-bit RGB values, perceptual metrics, and blend
// operations, independent of any output device.
//
// RGB is a plain 3-byte value with toml field tags, safe for config binding
// and for embedding in dense cell or pixel buffers. It satisfies
// image/color.Color, so values pass directly into image, draw, and GUI
// toolkit pipelines without conversion.
//
// Device concerns — terminal capability detection, xterm-256 quantization,
// SGR emission, framebuffer formats — belong in the consuming package.
//
// # Operations
//
//	Blend      linear alpha compositing
//	SoftLight  Perez soft light, gentler than linear alpha
//	Overlay    multiply on darks, screen on lights
//	Screen     always lightens, avoids the clipping harshness of Add
//	Add        saturating additive
//	Max        per-channel maximum
//	Scale      channel multiply, saturating
//	Grayscale  Rec. 601 luma
//	Lerp       linear interpolation
//
// # Concurrency
//
// All values are immutable and all operations are pure. Package-level palette
// variables are writable by the language but are read-only by contract.
package color
