package fractal

import(
	"image"
)

type Point8 struct{
	X, Y int
}

type Point64 struct {
	X, Y float64
}

type Rectangle8 struct {
	Min, Max Point8
}

func Rect8(x0, y0, x1, y1 int) Rectangle8 {
	return Rectangle8{Point8{x0, y0}, Point8{x1, y1}}
}

type Rectangle64 struct {
	Min, Max Point64
}

func Rect64(x0, y0, x1, y1 float64) Rectangle64 {
	return Rectangle64{Point64{x0, y0}, Point64{x1, y1}}
}


// Render the fractal to an Paletted, usually a sub-image
// call as a goroutine
func Render(img *image.Paletted, generator Generator, channel chan bool) {
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			img.SetColorIndex(x, y, generator.At(x, y))
			channel <- true
		}
	}
	close(channel)
}

func Mandelbrot(c complex128) func() complex128 {
	C := c
	Z := complex(0, 0)
	return func() complex128 {
		Z = Z*Z + C
		return Z
	}
}
