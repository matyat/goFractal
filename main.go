package main

import (
	"fmt"
	"goFractal/fractal"
	"image"
	"image/jpeg"
	"image/png"
	"launchpad.net/gnuflag"
	"os"
	"path/filepath"
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
	runtime.GOMAXPROCS(*cpus)

	gnuflag.Usage = func() {
		fmt.Print(usage)
		gnuflag.PrintDefaults()
	}

	gnuflag.Parse(true)
	configFileName := gnuflag.Arg(0)
	configFormat := filepath.Ext(configFileName)

	var parse func(string) (fractal.Renderer, error)
	switch configFormat{
	case ".xml":
		parse = fractal.ParseXml
	default:
		fmt.Println("Sorry %v input is not supported.\n", configFormat)
		return
	}

	renderer, renderErr := parse(configFileName)
	if renderErr != nil {
		fmt.Println(renderErr)
		return
	}

	if *imgFileName == "" {
		// Strip off the config extension and replace with .png and use
		// that as an output if an output file is not specified.
		*imgFileName = configFileName[:strings.LastIndex(configFileName, ".")] + "png"
	}
	imgFormat := filepath.Ext(*imgFileName)

	imgFile, fileErr := os.Create(*imgFileName)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	// set up the encode before we start rendering so we can
	// bail out if there are any errors
	var img image.Image
	var encode func() error
	switch imgFormat {
	case ".png":
		encode = func() error {
			return png.Encode(imgFile, img)
		}
	case ".jpg", ".jpeg":
		encode = func() error {
			opts := jpeg.Options{
				Quality: *imgQuality,
			}
			return jpeg.Encode(imgFile, img, &opts)
		}
	default:
		fmt.Printf("Sorry, image format %v is not supported.\n", imgFormat)
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
	img = renderer.GetImage()

	fmt.Print("\rFinished rendering in ", time.Since(t), " on ", *cpus, " CPUs with ", *threads, " threads\n")

	encodingErr := encode()
	if encodingErr != nil {
		fmt.Println(encodingErr)
		return
	}
}
