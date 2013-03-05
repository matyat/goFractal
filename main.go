package main

import (
	"fmt"
	"goFractal/fractal"
	"image"
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

func ConvertPaletteToRGBA(img *image.Paletted) *image.RGBA {
	output := image.NewRGBA(img.Bounds())
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			fmt.Println(img.ColorIndexAt(x, y), "->", img.At(x, y))
			output.Set(x, y, img.At(x, y))
		}
	}

	return output
}

// Generates a colour palette 
func GenerateColorPalette(levels uint32) color.Palette {
	palette := make([]color.Color, levels)
	for i := uint32(0); i < levels-1; i++ {
		n := float64(i) / float64(levels)
		palette[i] = HSVToRGB(360*n, 0.8, 1.0)
	}
	palette[levels-1] = color.RGBA{0, 0, 0, 255}
	return palette
}

// Split an image into a number of virtal strips
func SplitImage(img *fractal.Paletted_uint32, n int) []image.Image {
	strips := make([]image.Image, n)
	bounds := img.Bounds()
	h_step := int(float64(bounds.Dx()) / float64(n))

	// image width divided by n will not always be an integer
	// so we may have to add/remove a few columns from the last
	// strip 
	offset := bounds.Dx() - n*h_step

	for i := 0; i < n; i++ {
		x0 := bounds.Min.X + i*h_step
		x1 := bounds.Min.X + (i+1)*h_step
		if i == n-1 { // if last strip
			x1 += offset
		}
		strips[i] = img.SubImage(image.Rect(x0, bounds.Min.Y, x1, bounds.Max.Y))
	}
	return strips
}

// function for measuring progress
func ProgressBar(max int) func(int) {
	last_percent_str := ""
	return func(val int) {
		percent_float := 100 * float64(val) / float64(max)
		percent_str := strconv.FormatFloat(percent_float, 'f', 1, 64)
		if percent_str != last_percent_str {
			fmt.Print("\r", percent_str, "%")
		}
		last_percent_str = percent_str
	}
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

	palette := GenerateColorPalette(1000)
	size := 1024.0
	img := fractal.NewPaletted_uint32(image.Rect(0, 0, int(size), int(size*0.4459)), palette)

	fractal_generator := fractal.Generator{
		Domain: fractal.Rect64(0.276185, 0.479000198,
			0.367588933, 0.519762846),
		Size: fractal.Rect8(img.Bounds().Min.X, img.Bounds().Min.Y,
			img.Bounds().Max.X, img.Bounds().Max.Y),
		Bailout:       3.0,
		MaxIterations: uint32(len(palette) - 1),
		Function:      fractal.Mandelbrot,
	}

	// for each CPU create a channel and a strip of the image to render
	sub_images := SplitImage(img, threads)
	channels := make([]chan bool, threads)

	for i := range channels {
		channels[i] = make(chan bool, 512)
		go fractal.Render(sub_images[i].(*fractal.Paletted_uint32), fractal_generator, channels[i])
	}

	//measure the goroutines progress and wait for them to finish
	progress_bar := ProgressBar(len(img.Pix))
	t := time.Now()
	pixels_processed := 0
	done := false
	for !done {
		done = true
		for i := range channels {
			_, open := <-channels[i]
			if open {
				done = false
				pixels_processed++
			}
		}
		progress_bar(pixels_processed)
	}

	fmt.Print("\rFinished rendering in ", time.Since(t), " on ", cpus, " CPUs with ", threads, " threads\n")

	imgFile, _ := os.Create("image.png")
	png.Encode(imgFile, img)
}
