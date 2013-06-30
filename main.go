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
	threads := cpus

	color_wheel := fractal.NewColorWheel(3, 255)
	color_wheel.InfColor = color.RGBA{0, 0, 0, 255}

	red_node := fractal.ColorNode{color.RGBA{255, 0, 0, 255}, 0}
	green_node := fractal.ColorNode{color.RGBA{0, 255, 0, 255}, math.Pi * 2 / 3}
	blue_node := fractal.ColorNode{color.RGBA{0, 0, 255, 255}, math.Pi * 4 / 3}

	color_wheel.ColorNodes = []fractal.ColorNode{red_node, green_node, blue_node}

	size := 256
	viewport := fractal.Viewport{
		Location: complex(0, 0),
		Scale:    1.5 / float64(size),
		Rotation: -25 * math.Pi / 180,
		Width:    2 * size,
		Height:   size,
	}

	renderer := fractal.Renderer{
		Viewport: viewport,
		Generator: fractal.Julia(complex(-0.8, 0.156), 1600),
		ColorWheel: color_wheel,
		Threads: threads,
		Multisampling: 4,
	}

	renderer.Render()

	t := time.Now()
	for renderer.Rendering() {
		prog_str := strconv.FormatFloat(100*renderer.Progress(), 'f', 1, 64)
		fmt.Print("\r", prog_str, "%")
	}

	img := renderer.GetImage()

	fmt.Print("\rFinished rendering in ", time.Since(t), " on ", cpus, " CPUs with ", threads, " threads\n")

	imgFile, _ := os.Create("image.png")
	png.Encode(imgFile, img)
}
