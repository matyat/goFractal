package fractal

import (
	"image"
	"image/color"
	"sync/atomic"
)

type Renderer struct {
	ViewPort                     ViewPort
	Generator                    Generator
	ColorWheel                   ColorWheel
	img                          *image.RGBA
	processedPixels, totalPixels uint64
}

func (ren *Renderer) Render(threads int) {
	ren.ColorWheel.generate()

	ren.img = image.NewRGBA(image.Rect(0, 0, ren.ViewPort.Width, ren.ViewPort.Height))

	ren.processedPixels = 0
	ren.totalPixels = uint64(ren.ViewPort.Width * ren.ViewPort.Height)

	// an array of values to make sure each pixel is only rendered once
	pixelLock := make([]int64, ren.totalPixels)

	for t := 1; t < threads+1; t++ {
		m := float64(ren.ViewPort.Multisampling)
		mSqr := uint64(ren.ViewPort.Multisampling * ren.ViewPort.Multisampling)

		// create a thread
		go func(id int) {
			for y := 0; y < ren.ViewPort.Height; y++ {
				stride := y * ren.ViewPort.Width
				for x := 0; x < ren.ViewPort.Width; x++ {
					// if pixelLock is assigned any other than 0, the pixel has
					// already been processed
					if atomic.CompareAndSwapInt64(
						&pixelLock[x+stride], 0, int64(id)) {
						var R, G, B uint64

						//sample sub-pixels when multisampling
						for sx := 0; sx < ren.ViewPort.Multisampling; sx++ {
							for sy := 0; sy < ren.ViewPort.Multisampling; sy++ {
								x0 := float64(x) + float64(sx)/m
								y0 := float64(y) + float64(sy)/m

								itr := ren.Generator.EscapeAt(ren.ViewPort.ComplexAt(x0, y0))
								r, g, b, _ := ren.ColorWheel.ColorAt(itr).RGBA()

								R += uint64(r >> 8)
								G += uint64(g >> 8)
								B += uint64(b >> 8)
							}
						}

						// Average each colour
						R /= mSqr
						G /= mSqr
						B /= mSqr

						ren.img.Set(x, y, color.RGBA{uint8(R), uint8(G), uint8(B), 255})
						atomic.AddUint64(&ren.processedPixels, 1)
					}
				}
			}
		}(t)
	}
}

func (ren *Renderer) Rendering() bool {
	return !(ren.processedPixels == ren.totalPixels)
}

func (ren *Renderer) Progress() float64 {
	return float64(ren.processedPixels) / float64(ren.totalPixels)
}

func (ren *Renderer) GetImage() *image.RGBA {
	return ren.img
}
