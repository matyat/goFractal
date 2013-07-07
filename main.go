package main

import (
	"fmt"
	"goFractal/fractal"
	"image/png"
	"image/jpeg"
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
var img_quality *int = gnuflag.Int(
	"image-quality", 100, "Image quality (where applicable).",
)

func main() {
	gnuflag.Usage = func() {
		fmt.Print(usage)
		gnuflag.PrintDefaults()
	}

	gnuflag.Parse(true)
	xml_filename := gnuflag.Arg(0)

	if xml_filename == "" || xml_filename[len(xml_filename)-3:] != "xml" {
		gnuflag.Usage()
		return
	}

	if *img_filename == "" {
		//strip off the .xml and replace with .png
		*img_filename = xml_filename[:len(xml_filename)-3] + "png"
	}

	img_format := (*img_filename)[len(*img_filename)-3:]

	runtime.GOMAXPROCS(*cpus)

	renderer, render_err := fractal.ParseXml(xml_filename)
	if render_err != nil {
		fmt.Println(render_err)
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

	imgFile, file_err := os.Create(*img_filename)

	if file_err != nil {
		fmt.Println(file_err)
		return
	}

	var encoding_err error
	switch(img_format){
	case "png":
		encoding_err = png.Encode(imgFile, img)
	case "jpg":
		opts := jpeg.Options{
			Quality: *img_quality,
		}
		encoding_err = jpeg.Encode(imgFile, img, &opts)
	default:
		fmt.Println("Image format not supported: ", img_format)
		return
	}

	if encoding_err != nil{
		fmt.Println(encoding_err)
		return
	}
}
