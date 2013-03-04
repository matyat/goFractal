package fractal

import(
	"image"
	"math/cmplx"
)

type Mandelbrot struct {
	Start, End    Pt_f64
	ImageSize     image.Rectangle
	Bailout       float64
	MaxIterations int
}

// gets the colour index at pixel X, Y for the Mandelbrot set
// using the escape time algorithm
func (m Mandelbrot) At(X, Y int) uint8 {
	// convert X & Y out of pixel space
	x0 := float64(X-m.ImageSize.Min.X)/
		float64(m.ImageSize.Max.X-m.ImageSize.Min.X)*
		(m.End.X-m.Start.X) + m.Start.X
	y0 := float64(Y-m.ImageSize.Min.Y)/
		float64(m.ImageSize.Max.Y-m.ImageSize.Min.Y)*
		(m.End.Y-m.Start.Y) + m.Start.Y

	z := complex(0, 0)
	c := complex(x0, y0)

	itr := 0
	for ; cmplx.Abs(z) < m.Bailout && itr < m.MaxIterations; itr++ {
		z = z*z + c
	}

	return uint8(itr)
}
