package fractal

import (
	"encoding/xml"
	"errors"
	"image/color"
	"io/ioutil"
)

// intermediant stucts for parsing xml files
type rendererIntr struct {
	Type          string   `xml:",attr"`
	Bailout       float64  `xml:",attr"`
	MaxIterations int      `xml:",attr"`
	C             string   `xml:",attr"`
	ViewPort      ViewPort `xml:"ViewPort"`
	ColorWheel    colorWheelIntr
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

// Parse an .xml file and return a Renderer
func ParseXml(filename string) (Renderer, error) {
	renderer := Renderer{}

	content, fErr := ioutil.ReadFile(filename)
	if fErr != nil {
		return renderer, fErr
	}

	intr := rendererIntr{}
	xErr := xml.Unmarshal(content, &intr)
	if xErr != nil {
		return renderer, xErr
	}

	switch intr.Type {
	case "Mandelbrot":
		renderer.Generator = Mandelbrot(intr.Bailout, intr.MaxIterations)

	case "Julia":
		C, jErr := parseCmplxString(intr.C)
		if jErr != nil {
			return renderer, jErr
		}
		renderer.Generator = Julia(C, intr.Bailout, intr.MaxIterations)

	default:
		return renderer, errors.New("Error no known fractal type: " + intr.Type)
	}

	renderer.ViewPort = intr.ViewPort
	renderer.ViewPort.Rotation = degToRad(renderer.ViewPort.Rotation)

	//get the colour for infinity
	colorwheelintr := intr.ColorWheel
	infcolor := color.RGBA{
		colorwheelintr.InfColor.Red,
		colorwheelintr.InfColor.Green,
		colorwheelintr.InfColor.Blue,
		colorwheelintr.InfColor.Alpha,
	}

	//populate the colour nodes
	colourNodes := make([]ColorNode, len(colorwheelintr.Nodes))
	for i := range colourNodes {
		nodeintr := colorwheelintr.Nodes[i]
		color := color.RGBA{
			nodeintr.Red,
			nodeintr.Green,
			nodeintr.Blue,
			nodeintr.Alpha,
		}

		colourNodes[i] = ColorNode{
			Color: color,
			Angle: degToRad(nodeintr.Angle),
		}
	}

	renderer.ColorWheel = ColorWheel{
		Radius:      colorwheelintr.Radius,
		PaletteSize: colorwheelintr.Res,
		InfColor:    infcolor,
		ColorNodes:  colourNodes,
	}

	return renderer, nil
}
