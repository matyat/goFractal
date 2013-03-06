package main

import (
	"fmt"
	"goFractal/fractal"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"strconv"
	"time"
)

// convert from HSV to RGB color space
func HSVToRGB(h, s, v float64) color.RGBA {
	hh := h / 60
	i := int(hh)
	ff := hh - float64(i)

	p := v * (1 - s)
	q := v * (1 - s*ff)
	t := v * (1 - s*(1-ff))

	var r, g, b float64
	switch i {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	case 5:
		r = v
		g = p
		b = q
	}
	return color.RGBA{uint8(255 * r), uint8(255 * g), uint8(255 * b), 255}
}

// Generates a colour palette 
func GenerateColorPalette(levels uint32) color.Palette {
	palette := make([]color.Color, levels)
	for i := uint32(0); i < levels; i++ {
		n := float64(i) / float64(levels)
		palette[i] = HSVToRGB(300*n, 0.8, 1.0)
	}
	return palette
}

func main() {
	// get the number of CPUs, and set the go runtime to utilise them
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var threads int
	if cpus == 1 {
		// no point in multithreading on a single-core
		threads = 1
	} else {
		// the number of threads should be higher than the number of CPUs
		// so the load across all the CPUs remains fairly constant if a few
		// threads finish early

		// TODO: work out the optimal number of threads per CPU
		threads = 2 * cpus
	}

	palette := GenerateColorPalette(255)
	size := 1024

	viewport := fractal.Viewport{
		Location: complex(-0.5, 0),
		Scale:    2 / float64(size),
		Rotation: 0,
		Width:    size,
		Height:   size,
	}

	fractal_generator := fractal.Generator{
		Bailout:       2.0,
		MaxIterations: uint32(len(palette) - 1),
		Function:      fractal.Mandelbrot,
	}

	monitor := fractal.NewMonitor()
	img := fractal.Render(viewport, fractal_generator, palette, monitor, 1, threads)
	// loop until the monitor says the render is complete
	t := time.Now()
	var last_prog_str string
	for {
		if prog, done := monitor.Progress(); done {
			break
		} else {
			prog_str := strconv.FormatFloat(100*prog, 'f', 1, 64)
			if prog_str != last_prog_str {
				fmt.Print("\r", prog_str, "%")
			}
			last_prog_str = prog_str
		}
	}

	fmt.Print("\rFinished rendering in ", time.Since(t), " on ", cpus, " CPUs with ", threads, " threads\n")

	imgFile, _ := os.Create("image.png")
	png.Encode(imgFile, img)
}
