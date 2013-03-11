package fractal

import (
	"image"
	"image/color"
	"math"
)
// convert from HSV to RGB color space
func HSVToRGB(h, s, v float64) color.RGBA {
	hh := h / 60
	i := int(hh)
	ff := hh - float64(i)

	p := v * (1 - s)
	q := v * (1 - s*ff)
	t := v * (1 - s*(1-ff))

	var r, g, b float64
	switch i {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	case 5:
		r = v
		g = p
		b = q
	}
	return color.RGBA{uint8(255 * r), uint8(255 * g), uint8(255 * b), 255}
}

// Paletted_float64 is an in-memory image of float64 indices into a given palette.
type Paletted_float64 struct {
	// Pix holds the image's pixels, as palette indices. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*1].
	Pix []float64
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
	// Palette is the image's palette.
	Palette *Palette
}

func (p *Paletted_float64) ColorModel() color.Model {
	// this should allow us to create a Truecolor 8/16 bit png when
	// using a large palatte

	if len(p.Palette.cPalette) > 255 {
		// will cause the encoder to use Truecolor 8 bit
		return color.RGBAModel
	}

	if len(p.Palette.cPalette) > 16777215 {
		// will cause the encoder to use Truecolor 16 bit
		return color.RGBA64Model
	}
	// FIXME: should return palette
	return color.RGBAModel
}

func (p *Paletted_float64) Bounds() image.Rectangle { return p.Rect }

func (p *Paletted_float64) At(x, y int) color.Color {
	if len(p.Palette.cPalette) == 0 {
		return nil
	}
	if !(image.Point{x, y}.In(p.Rect)) {
		return p.Palette.cPalette[0]
	}
	i := p.PixOffset(x, y)
	return p.Palette.ColorAt(p.Pix[i])
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Paletted_float64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

//func (p *Paletted_float64) Set(x, y int, c color.Color) {
//	if !(image.Point{x, y}.In(p.Rect)) {
//		return
//	}
//	i := p.PixOffset(x, y)
//	p.Pix[i] = float64(p.Palette.Index(c))
//}

func (p *Paletted_float64) ColorIndexAt(x, y int) float64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return 0
	}
	i := p.PixOffset(x, y)
	return p.Pix[i]
}

func (p *Paletted_float64) SetColorIndex(x, y int, index float64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = index
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Paletted_float64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are image.Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Paletted_float64{
			Palette: p.Palette,
		}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Paletted_float64{
		Pix:     p.Pix[i:],
		Stride:  p.Stride,
		Rect:    p.Rect.Intersect(r),
		Palette: p.Palette,
	}
}

// Opaque scans the entire image and returns whether or not it is fully opaque.
func (p *Paletted_float64) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}

	for _, c := range p.Palette.cPalette {
		_, _, _, a := c.RGBA()
		if a != 0xffff {
			return false
		}
	}
	return true
}

// NewPaletted_float64 returns a new Paletted with the given width, height and palette.
func NewPaletted_float64(r image.Rectangle, p *Palette) *Paletted_float64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 1*w*h)
	return &Paletted_float64{pix, 1 * w, r, p}
}

// Special Palette for fractals which can support normailised iteration
type Palette struct {
	cPalette            color.Palette
	InfColor            color.Color
	Length, SmoothScale int
}

func NewPalette(size, smooth_scale int) *Palette {
	return &Palette{
		Length:      size,
		SmoothScale: smooth_scale,
	}
}

func (palette *Palette) generate() {
	palette_len := palette.Length*palette.SmoothScale
	palette.cPalette = make([]color.Color, palette_len)
	for i := range palette.cPalette {
		n := float64(i) / float64(palette_len)
		palette.cPalette[i] = HSVToRGB(60 + 120*n, n, 1-n)
	}
}

func (palette *Palette) ColorAt(itr float64) color.Color {
	if math.IsNaN(itr) || math.IsInf(itr, 0) {
		return palette.InfColor
	}

	_, f:= math.Modf(itr/float64(2*palette.Length))
	idx := int(math.Abs(1-2*f)*float64(palette.Length*palette.SmoothScale))
	return palette.cPalette[idx]
}
