package fractal

import (
	"image"
	"image/color"
)

// Paletted_uint32 is an in-memory image of uint32 indices into a given palette.
type Paletted_uint32 struct {
	// Pix holds the image's pixels, as palette indices. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []uint32
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Palette is the image's palette.
	Palette color.Palette
}

func (p *Paletted_uint32) ColorModel() color.Model {
	// this should allow us to create a Truecolor 8/16 bit png when
	// using a large palatte

	if len(p.Palette) > 255 {
		// will cause the encoder to use Truecolor 8 bit
		return color.RGBAModel
	}

	if len(p.Palette) > 16777215 {
		// will cause the encoder to use Truecolor 16 bit
		return color.RGBA64Model
	}

	return p.Palette
}

func (p *Paletted_uint32) Bounds() image.Rectangle { return p.Rect }

func (p *Paletted_uint32) At(x, y int) color.Color {
	if len(p.Palette) == 0 {
		return nil
	}
	if !(image.Point{x, y}.In(p.Rect)) {
		return p.Palette[0]
	}
	i := p.PixOffset(x, y)
	return p.Palette[p.Pix[i]]
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Paletted_uint32) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

func (p *Paletted_uint32) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = uint32(p.Palette.Index(c))
}

func (p *Paletted_uint32) ColorIndexAt(x, y int) uint32 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return 0
	}
	i := p.PixOffset(x, y)
	return p.Pix[i]
}

func (p *Paletted_uint32) SetColorIndex(x, y int, index uint32) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = index
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Paletted_uint32) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are image.Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Paletted_uint32{
			Palette: p.Palette,
		}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Paletted_uint32{
		Pix:     p.Pix[i:],
		Stride:  p.Stride,
		Rect:    p.Rect.Intersect(r),
		Palette: p.Palette,
	}
}

// Opaque scans the entire image and returns whether or not it is fully opaque.
func (p *Paletted_uint32) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}

	for _, c := range p.Palette {
		_, _, _, a := c.RGBA()
		if a != 0xffff {
			return false
		}
	}
	return true
}

// NewPaletted_uint32 returns a new Paletted with the given width, height and palette.
func NewPaletted_uint32(r image.Rectangle, p color.Palette) *Paletted_uint32 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint32, 1*w*h)
	return &Paletted_uint32{pix, 1 * w, r, p}
}
