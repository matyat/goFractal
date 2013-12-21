package fractal

import (
	"math"
)

// Struct to define the area rendered
type ViewPort struct {
	Location       complex128
	LocationStr    string  `xml:"Location,attr"`
	Scale          float64 `xml:",attr"`
	Rotation       float64 `xml:",attr"`
	Width,Height          int     `xml:",attr"`
	Multisampling  int     `xml:",attr"`
	cosRot, sinRot float64
	initialised    bool
}

// translate a pixel co-ord to a complex number
func (vp ViewPort) ComplexAt(x, y float64) complex128 {
	// only need to run once
	if !vp.initialised {
		vp.cosRot = math.Cos(vp.Rotation)
		vp.sinRot = math.Sin(vp.Rotation)
		vp.Location, _ = parseCmplxString(vp.LocationStr)
		vp.initialised = true
	}

	//move origin to centre
	x -= float64(vp.Width) / 2
	y -= float64(vp.Height) / 2

	// scale
	x /= vp.Scale
	y /= vp.Scale

	// rotate
	x, y = (x*vp.cosRot - y*vp.sinRot), (x*vp.sinRot + y*vp.cosRot)

	// move
	x += real(vp.Location)
	y += imag(vp.Location)

	return complex(x, y)
}
