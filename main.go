package main

import (
	"fmt"
	"goFractal/fractal"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"
)

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
		threads = 4 * cpus
	}

	color_wheel := fractal.NewColorWheel(3, 255)
	color_wheel.InfColor = color.RGBA{0, 0, 0, 255}
	color_wheel.AddColor(color.RGBA{92, 4, 4, 255}, 0)
	color_wheel.AddColor(color.RGBA{95, 53, 35, 255}, math.Pi*1/3)
	color_wheel.AddColor(color.RGBA{138, 125, 192, 255}, math.Pi*2/3)
	color_wheel.AddColor(color.RGBA{185, 160, 222, 255}, math.Pi)
	color_wheel.AddColor(color.RGBA{248, 251, 218, 255}, math.Pi*4/3)
	color_wheel.AddColor(color.RGBA{123, 96, 71, 255}, math.Pi*5/3)

	size := 8000

	viewport := fractal.Viewport{
		Location: complex(0, 0),
		Scale:    1.5 / float64(size),
		Rotation: -25 * math.Pi / 180,
		Width:    2 * size,
		Height:   size,
	}

	monitor := fractal.NewMonitor()
	img := fractal.Render(viewport, fractal.Julia(complex(-0.8, 0.156), 1600),
		color_wheel, monitor, 4, threads)

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
