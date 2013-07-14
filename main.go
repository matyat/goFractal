package main

import (
	"errors"
	"fmt"
	"goFractal/fractal"
	"image/jpeg"
	"image/png"
	"launchpad.net/gnuflag"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var usage string = "Usage: goFractal [OPTIONS]... [XML_RENDER_FILE]...\n" +
	"Render the fractal defined in XML_RENDER_FILE\n"

var imgFileName *string = gnuflag.String(
	"output", "", "Output image filename.",
)
var cpus *int = gnuflag.Int(
	"cpus", runtime.NumCPU(), "Number of CPUs to untilise.",
)
var threads *int = gnuflag.Int(
	"threads", runtime.NumCPU(), "Number of rendering threads.",
)
var imgQuality *int = gnuflag.Int(
	"image-quality", 100, "Image quality (where applicable).",
)

func main() {
	gnuflag.Usage = func() {
		fmt.Print(usage)
		gnuflag.PrintDefaults()
	}

	gnuflag.Parse(true)
	xmlFileName := gnuflag.Arg(0)

	if xmlFileName == "" || !strings.HasSuffix(xmlFileName, ".xml") {
		fmt.Println("Sorry, only .xml input files are valid")
		gnuflag.Usage()
		return
	}

	if *imgFileName == "" {
		// strip off the .xml and replace with .png and use that as an output 
		// if an output file is not specified
		*imgFileName = xmlFileName[:strings.LastIndex(*imgFileName, ".")] + "png"
	}

	imgFormat := (*imgFileName)[strings.LastIndex(*imgFileName, ".") + 1:]

	runtime.GOMAXPROCS(*cpus)

	renderer, renderErr := fractal.ParseXml(xmlFileName)
	if renderErr != nil {
		fmt.Println(renderErr)
		return
	}

	if !(imgFormat == "png" || imgFormat == "jpg" || imgFormat == "jpeg"){
		fmt.Println("Sorry, image format", imgFormat, "is not supported")
		return
	}

	// begin the render
	renderer.Render(*threads)

	// wait untill it is completed
	t := time.Now()
	for renderer.Rendering() {
		progStr := strconv.FormatFloat(100*renderer.Progress(), 'f', 1, 64)
		fmt.Print("\r", progStr, "%")
	}

	img := renderer.GetImage()

	fmt.Print("\rFinished rendering in ", time.Since(t), " on ", *cpus, " CPUs with ", *threads, " threads\n")

	imgFile, fileErr := os.Create(*imgFileName)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	var encodingErr error
	switch imgFormat {
	case "png":
		encodingErr = png.Encode(imgFile, img)
	case "jpg", "jpeg":
		opts := jpeg.Options{
			Quality: *imgQuality,
		}
		encodingErr = jpeg.Encode(imgFile, img, &opts)
	default:
		encodingErr = errors.New("Image format not supported: " + imgFormat)
	}

	if encodingErr != nil {
		fmt.Println(encodingErr)
		return
	}
}
