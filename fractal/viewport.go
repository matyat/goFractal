package fractal

import (
	"math"
)

// Struct to define the area rendered
type ViewPort struct {
	Location        complex128
	Scale, Rotation float64
	Width, Height   int
}

// translate a pixel co-ord to a complex number
func (view ViewPort) ComplexAt(X, Y float64) complex128 {
	//move origin to centre
	x := X - float64(view.Width)/2
	y := Y - float64(view.Height)/2

	// scale
	x /= view.Scale
	y /= view.Scale

	// rotate
	cos_rt := math.Cos(view.Rotation)
	sin_rt := math.Sin(view.Rotation)
	x, y = (x*cos_rt - y*sin_rt), (x*sin_rt + y*cos_rt)

	// move
	x += real(view.Location)
	y += imag(view.Location)

	return complex(x, y)
}
