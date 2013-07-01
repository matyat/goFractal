package main

import (
	"fmt"
	"goFractal/fractal"
	"image/png"
	"launchpad.net/gnuflag"
	"os"
	"runtime"
	"strconv"
	"time"
)

var img_filename *string= gnuflag.String(
	"output", "", "Output image filename.",
)
var cpus *int = gnuflag.Int(
	"cpus", runtime.NumCPU(), "Number of CPUs to untilise.",
)
var threads *int = gnuflag.Int(
	"threads", runtime.NumCPU(), "Number of rendering threads.",
)

func main() {
	gnuflag.Parse(true)
	xml_filename := gnuflag.Arg(0)

	if *img_filename == "" {
		//strip off the .xml and replace with .png
		*img_filename = xml_filename[:len(xml_filename)-3] + "png"
	}

	runtime.GOMAXPROCS(*cpus)

	renderer, err := fractal.ParseXml(xml_filename)
	if err != nil {
		fmt.Println(err)
	}

	renderer.Render(*threads)

	t := time.Now()
	for renderer.Rendering() {
		prog_str := strconv.FormatFloat(100*renderer.Progress(), 'f', 1, 64)
		fmt.Print("\r", prog_str, "%")
	}

	img := renderer.GetImage()

	fmt.Print("\rFinished rendering in ", time.Since(t), " on ", *cpus, " CPUs with ", *threads, " threads\n")

	imgFile, _ := os.Create(*img_filename)
	png.Encode(imgFile, img)
}
