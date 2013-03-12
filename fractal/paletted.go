package fractal

import (
	"image/color"
	"math"
)

// internal type for ColorWheel to interpolate with, holds a colour and the it's angle
// in the wheel
type interpolColor struct {
	color color.Color
	angle float64
}

// Special Palette for fractals, which loops back on it self
// Set the InfColor to the color you want bounded regions (infinte
// iterations) to be
type ColorWheel struct {
	cPalette       color.Palette
	InfColor       color.Color
	interpolColors []interpolColor
	Radius         float64
}

// Greates a new ColorWheel with the given radius, palette size and colour for
// infinty
func NewColorWheel(radius float64, palette_size int, inf_color color.Color) *ColorWheel {
	return &ColorWheel{
		cPalette: make([]color.Color, palette_size),
		InfColor: inf_color,
		Radius:   radius,
	}
}

// add a colour to the ColorWheel
func (col_wheel *ColorWheel) AddColor(col color.Color, angle float64) {
	col_wheel.interpolColors = append(col_wheel.interpolColors, interpolColor{col, angle})
}

// internal method of getting the pos of the internal palette at a given angle
func (col_wheel *ColorWheel) getPalettePosAt(angle float64) float64 {
	// wrap around 2pi
	_, f := math.Modf(angle / (2 * math.Pi))
	s := f * 2 * math.Pi
	p := 2 * math.Pi / float64(len(col_wheel.cPalette))
	return s / p
}

func (col_wheel *ColorWheel) generate() {
	const M = float64(1<<16 - 1)
	//for each colour to interpolate between
	for i := range col_wheel.interpolColors {
		col_A := col_wheel.interpolColors[i]
		var col_B interpolColor
		if i != len(col_wheel.interpolColors)-1 {
			col_B = col_wheel.interpolColors[i+1]
		} else {
			// if this is the last colour, we want
			// to interpolate between this and the 
			// first colour
			col_B = col_wheel.interpolColors[0]
			col_B.angle += 2 * math.Pi
		}

		start_idx := int(math.Ceil(col_wheel.getPalettePosAt(col_A.angle)))
		end_idx := int(math.Ceil(col_wheel.getPalettePosAt(col_B.angle)))

		Ar, Ag, Ab, Aa := col_A.color.RGBA()
		Br, Bg, Bb, Ba := col_B.color.RGBA()
		
		// begin interpolation
		for n := start_idx; n != end_idx; n++ {
			cur_angle := 2 * math.Pi * float64(n) / float64(len(col_wheel.cPalette))
			m := (cur_angle - col_A.angle) / (col_B.angle - col_A.angle)
			blend := func(a, b uint32) uint8 {
				fa := float64(a) * (1 - m)
				fb := float64(b) * m
				return uint8(255 * (fa + fb) / M)
			}

			col_wheel.cPalette[n] = color.RGBA{blend(Ar, Br), blend(Ag, Bg), blend(Ab, Bb), blend(Aa, Ba)}

			if n == len(col_wheel.cPalette)-1 {
				n = -1
			}
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
