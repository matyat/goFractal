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
	cPalette    color.Palette
	PaletteSize int
	InfColor    color.Color
	ColorNodes  []ColorNode
	Radius      float64
}

// Creates a new ColorWheel with the given radius, palette size and colour for
// infinty
func NewColorWheel(radius float64, paletteSize int) *ColorWheel {
	return &ColorWheel{
		PaletteSize: paletteSize,
		Radius:      radius,
	}
}

// internal method of getting the pos of the internal palette at a given angle
func (colWheel *ColorWheel) getPalettePosAt(angle float64) float64 {
	p := 2 * math.Pi / float64(len(colWheel.cPalette))
	return angle / p
}

func (colWheel *ColorWheel) generate() {
	colWheel.cPalette = make([]color.Color, colWheel.PaletteSize)

	const M = float64(1<<16 - 1)
	paletteSize := len(colWheel.cPalette)
	//for each colour to interpolate between
	for i := range colWheel.ColorNodes {
		colA := colWheel.ColorNodes[i]
		var colB ColorNode

		if i != len(colWheel.ColorNodes)-1 {
			colB = colWheel.ColorNodes[i+1]
		} else {
			// if this is the last colour, we want
			// to interpolate between this and the
			// first colour
			colB = colWheel.ColorNodes[0]
			colB.Angle += 2 * math.Pi
		}

		startIdx := int(math.Ceil(colWheel.getPalettePosAt(colA.Angle)))
		endIdx := int(math.Ceil(colWheel.getPalettePosAt(colB.Angle)))

		Ar, Ag, Ab, Aa := colA.Color.RGBA()
		Br, Bg, Bb, Ba := colB.Color.RGBA()

		// begin interpolation
		for n := startIdx; n != endIdx; n++ {
			curAngle := 2 * math.Pi * float64(n) / float64(paletteSize)
			m := (curAngle - colA.Angle) / (colB.Angle - colA.Angle)
			blend := func(a, b uint32) uint8 {
				fa := float64(a) * (1 - m)
				fb := float64(b) * m
				return uint8(255 * (fa + fb) / M)
			}
			colWheel.cPalette[n%paletteSize] = color.RGBA{blend(Ar, Br), blend(Ag, Bg), blend(Ab, Bb), blend(Aa, Ba)}
		}
	}
}

// Get the colour of a given number of normalised iterations
func (colWheel *ColorWheel) ColorAt(itr float64) color.Color {
	// return the colour for infinity if infinity is given
	if math.IsNaN(itr) || math.IsInf(itr, 0) {
		return colWheel.InfColor
	}

	if itr < 0 {
		itr = 0
	}

	angle := itr / colWheel.Radius
	_, f := math.Modf(angle / (2 * math.Pi))
	angle = f * 2 * math.Pi

	idx := int(math.Floor(colWheel.getPalettePosAt(angle) + 0.5))
	if idx == len(colWheel.cPalette) {
		idx = 0
	}
	return colWheel.cPalette[idx]
}
