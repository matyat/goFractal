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
func (view Viewport) ComplexAt(X, Y int) complex128 {
	//move origin to centre
	x := float64(X) - float64(view.Width)/2
	y := float64(Y) - float64(view.Height)/2
	// rotate
	x = x*math.Cos(view.Rotation) + y*math.Sin(view.Rotation)
	y = -x*math.Sin(view.Rotation) + y*math.Cos(view.Rotation)
	// scale
	x *= view.Scale
	y *= view.Scale
	// move
	x += real(view.Location)
	y += imag(view.Location)

	return complex(x, y)
}
