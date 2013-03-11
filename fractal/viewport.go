package fractal

import (
	"math"
)

type Viewport struct {
	Location        complex128
	Scale, Rotation float64
	Width, Height   int
}

// translate a pixel co-ord to a complex number
func (view Viewport) ComplexAt(X, Y float64) complex128 {
	//move origin to centre
	x := X - float64(view.Width)/2
	y := Y - float64(view.Height)/2

	// scale
	x *= view.Scale
	y *= view.Scale

	// rotate
	x0, y0 := x, y
	x = x0*math.Cos(view.Rotation) - y0*math.Sin(view.Rotation)
	y = x0*math.Sin(view.Rotation) + y0*math.Cos(view.Rotation)

	// move
	x += real(view.Location)
	y += imag(view.Location)

	return complex(x, y)
}
