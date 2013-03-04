package fractal

import(
	"image"
)

type Pt_f64 struct {
	X, Y float64
}

type Fractal interface {
	At(X, Y int) uint8
}

// Render the fractal to an Paletted, usually a sub-image
// call as a goroutine
func Render(img *image.Paletted, fractal Fractal, channel chan bool) {
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			img.SetColorIndex(x, y, fractal.At(x, y))
			channel <- true
		}
	}
	close(channel)
}
