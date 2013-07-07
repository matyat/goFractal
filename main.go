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

var usage string = "Usage: goFractal [OPTIONS]... [XML_RENDER_FILE]...\n" +
	"Render the fractal defined in XML_RENDER_FILE\n"

var img_filename *string = gnuflag.String(
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

	gnuflag.Usage = func() {
		fmt.Print(usage)
		gnuflag.PrintDefaults()
	}

	if xml_filename == "" || xml_filename[len(xml_filename)-3:] != "xml" {
		gnuflag.Usage()
		return
	}

	if *img_filename == "" {
		//strip off the .xml and replace with .png
		*img_filename = xml_filename[:len(xml_filename)-3] + "png"
	}

	runtime.GOMAXPROCS(*cpus)

	renderer, err := fractal.ParseXml(xml_filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// begin the render
	renderer.Render(*threads)

	// wait untill it is completed
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
