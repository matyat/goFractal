package fractal

import (
	"image/color"
	"math"
)

// internal type for ColorWheel to interpolate with, holds a colour and the it's angle
// in the wheel
type ColorNode struct {
	Color color.Color
	Angle float64
}

// Special Palette for fractals, which loops back on it self
// Set the InfColor to the color you want bounded regions (infinte
// iterations) to be
type ColorWheel struct {
	cPalette   color.Palette
	InfColor   color.Color
	ColorNodes []ColorNode
	Radius     float64
}

//adds a new colour to the list of nodes TODO: Implement this stub
func (col_wheel *ColorWheel) addColor(node ColorNode) {
	//col_wheel.ColorNodes.
	return
}

// Greates a new ColorWheel with the given radius, palette size and colour for
// infinty
func NewColorWheel(radius float64, palette_size int) *ColorWheel {
	return &ColorWheel{
		cPalette: make([]color.Color, palette_size),
		Radius:   radius,
	}
}

// internal method of getting the pos of the internal palette at a given angle
func (col_wheel *ColorWheel) getPalettePosAt(angle float64) float64 {
	p := 2 * math.Pi / float64(len(col_wheel.cPalette))
	return angle / p
}

func (col_wheel *ColorWheel) generate() {
	const M = float64(1<<16 - 1)
	palette_size := len(col_wheel.cPalette)
	//for each colour to interpolate between
	for i := range col_wheel.ColorNodes {
		col_A := col_wheel.ColorNodes[i]
		var col_B ColorNode
		if i != len(col_wheel.ColorNodes)-1 {
			col_B = col_wheel.ColorNodes[i+1]
		} else {
			// if this is the last colour, we want
			// to interpolate between this and the 
			// first colour
			col_B = col_wheel.ColorNodes[0]
			col_B.Angle += 2 * math.Pi
		}

		start_idx := int(math.Ceil(col_wheel.getPalettePosAt(col_A.Angle)))
		end_idx := int(math.Ceil(col_wheel.getPalettePosAt(col_B.Angle)))

		Ar, Ag, Ab, Aa := col_A.Color.RGBA()
		Br, Bg, Bb, Ba := col_B.Color.RGBA()

		// begin interpolation
		for n := start_idx; n != end_idx; n++ {
			cur_angle := 2 * math.Pi * float64(n) / float64(palette_size)
			m := (cur_angle - col_A.Angle) / (col_B.Angle - col_A.Angle)
			blend := func(a, b uint32) uint8 {
				fa := float64(a) * (1 - m)
				fb := float64(b) * m
				return uint8(255 * (fa + fb) / M)
			}
			col_wheel.cPalette[n%palette_size] = color.RGBA{blend(Ar, Br), blend(Ag, Bg), blend(Ab, Bb), blend(Aa, Ba)}
		}
	}
}

func (col_wheel *ColorWheel) ColorAt(itr float64) color.Color {
	// return the colour for infinity if infinity is given
	if math.IsNaN(itr) || math.IsInf(itr, 0) {
		return col_wheel.InfColor
	}

	angle := itr / col_wheel.Radius
	_, f := math.Modf(angle / (2 * math.Pi))
	angle = f * 2 * math.Pi

	idx := int(math.Floor(col_wheel.getPalettePosAt(angle) + 0.5))
	if idx == len(col_wheel.cPalette) {
		idx = 0
	}
	return col_wheel.cPalette[idx]
}
