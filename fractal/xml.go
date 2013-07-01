package fractal

import (
	"encoding/xml"
	"errors"
	"image/color"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

// intermediant stucts for parsing xml files
type rendererIntr struct {
	Type          string  `xml:",attr"`
	Bailout       float64 `xml:",attr"`
	MaxIterations int     `xml:",attr"`
	C             string  `xml:",attr"`
	ViewPort      viewPortIntr
	ColorWheel    colorWheelIntr
}

type viewPortIntr struct {
	Location      string  `xml:",attr"`
	Scale         float64 `xml:",attr"`
	Rotation      float64 `xml:",attr"`
	Width         int     `xml:",attr"`
	Height        int     `xml:",attr"`
	Multisampling int     `xml:",attr"`
}

type colorWheelIntr struct {
	Res      int             `xml:",attr"`
	Radius   float64         `xml:",attr"`
	Nodes    []colorNodeIntr `xml:"Color"`
	InfColor colorNodeIntr
}

type colorNodeIntr struct {
	Red   uint8   `xml:",attr"`
	Green uint8   `xml:",attr"`
	Blue  uint8   `xml:",attr"`
	Alpha uint8   `xml:",attr"`
	Angle float64 `xml:",attr"`
}

// Parse a string into a complex number type.
// The complex number can have any number of real and imaginary
// components in any order i.e. "0.5 - 2 + 12i - 3 + 0.1i" is 
// perfectly valid and will return -4.5 + 12.1i.
func parseCmplxString(str string) (complex128, error) {
	var R float64
	var I float64
	var val_ptr *float64

	value_strs := strings.Split(str, " ")

	sign := 1.0
	for i := range value_strs {
		// if the string is a '+' or '-', we change the sign
		// and go the the next string
		if value_strs[i] == "+" {
			sign = 1.0
			continue
		} else if value_strs[i] == "-" {
			sign = -1.0
			continue
		}

		// if the last char is 'i', this is an imaginary number, so trim off
		// the i and change the val_ptr to I. Else, we must be dealing with
		// a real number.
		if last := len(value_strs[i]) - 1; last >= 0 && value_strs[i][last] == 'i' {
			value_strs[i] = value_strs[i][:last]
			val_ptr = &I
		} else {
			val_ptr = &R
		}

		val, err := strconv.ParseFloat(value_strs[i], 64)
		if err != nil {
			return complex(0, 0), err
		}
		*val_ptr += sign * val
	}

	return complex(R, I), nil
}

// convert degree to radians
func degToRad(deg float64) float64 {
	return math.Pi * deg / 180
}

// Parse an .xml file and return a Renderer
func ParseXml(filename string) (Renderer, error) {
	renderer := Renderer{}

	content, f_err := ioutil.ReadFile(filename)
	if f_err != nil {
		return renderer, f_err
	}

	intr := rendererIntr{}
	x_err := xml.Unmarshal(content, &intr)
	if x_err != nil {
		return renderer, x_err
	}

	viewport_loc, c_err := parseCmplxString(intr.ViewPort.Location)
	if c_err != nil {
		return renderer, c_err
	}

	renderer.Multisampling = intr.ViewPort.Multisampling

	switch intr.Type {
	case "Mandelbrot":
		renderer.Generator = Mandelbrot(intr.Bailout, intr.MaxIterations)

	case "Julia":
		C, j_err := parseCmplxString(intr.C)
		if j_err != nil {
			return renderer, j_err
		}
		renderer.Generator = Julia(C, intr.Bailout, intr.MaxIterations)

	default:
		return renderer, errors.New("Error no known fractal type: " + intr.Type)
	}

	renderer.ViewPort = ViewPort{
		Location: viewport_loc,
		Scale:    intr.ViewPort.Scale,
		Rotation: intr.ViewPort.Rotation,
		Width:    intr.ViewPort.Width,
		Height:   intr.ViewPort.Height,
	}

	//get the colour for infinity
	colorwheelintr := intr.ColorWheel
	infcolor := color.RGBA{
		colorwheelintr.InfColor.Red,
		colorwheelintr.InfColor.Green,
		colorwheelintr.InfColor.Blue,
		colorwheelintr.InfColor.Alpha,
	}

	//populate the colour nodes
	colour_nodes := make([]ColorNode, len(colorwheelintr.Nodes))
	for i := range colour_nodes {
		nodeintr := colorwheelintr.Nodes[i]
		color := color.RGBA{
			nodeintr.Red,
			nodeintr.Green,
			nodeintr.Blue,
			nodeintr.Alpha,
		}

		colour_nodes[i] = ColorNode{
			Color: color,
			Angle: degToRad(nodeintr.Angle),
		}
	}

	renderer.ColorWheel = ColorWheel{
		Radius:      colorwheelintr.Radius,
		PaletteSize: colorwheelintr.Res,
		InfColor:    infcolor,
		ColorNodes:  colour_nodes,
	}

	return renderer, nil
}
