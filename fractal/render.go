package fractal

import (
	"image"
	"image/color"
	"sync/atomic"
)

type Renderer struct {
	ViewPort         ViewPort
	Generator        Generator
	ColorWheel       ColorWheel
	Multisampling    int
	img              *image.RGBA
	processed_pixels uint64
	total_pixels     uint64
}

func (ren *Renderer) Render(threads int) {
	ren.ColorWheel.generate()

	ren.img = image.NewRGBA(image.Rect(0, 0, ren.ViewPort.Width, ren.ViewPort.Height))

	ren.processed_pixels = 0
	ren.total_pixels = uint64(ren.ViewPort.Width * ren.ViewPort.Height)

	// an array of values to make sure each pixel is only rendered once
	pixel_lock := make([]int64, ren.total_pixels)

	for t := 1; t < threads + 1; t++ {
		m := float64(ren.Multisampling)
		m_sqr := uint64(ren.Multisampling * ren.Multisampling)

		// create a thread
		go func(id int) {
			for y := 0; y < ren.ViewPort.Height; y++ {
				stride := y * ren.ViewPort.Width
				for x := 0; x < ren.ViewPort.Width; x++ {
					// if pixel_lock is assigned any other than 0, the pixel has
					// already been processed
					if atomic.CompareAndSwapInt64(
						&pixel_lock[x+stride], 0, int64(id)) {
						var R, G, B uint64

						//sample sub-pixels when multisampling
						for sx := 0; sx < ren.Multisampling; sx++ {
							for sy := 0; sy < ren.Multisampling; sy++ {
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
						R /= m_sqr
						G /= m_sqr
						B /= m_sqr

						ren.img.Set(x, y, color.RGBA{uint8(R), uint8(G), uint8(B), 255})
						atomic.AddUint64(&ren.processed_pixels, 1)
					}
				}
			}
		}(t)
	}
}

func (ren *Renderer) Rendering() bool {
	return !(ren.processed_pixels == ren.total_pixels)
}

func (ren *Renderer) Progress() float64 {
	return float64(ren.processed_pixels) / float64(ren.total_pixels)
}

func (ren *Renderer) GetImage() *image.RGBA {
	return ren.img
}
