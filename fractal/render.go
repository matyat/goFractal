package fractal

import (
	"image"
	"image/color"
)

// Split an image into a number of virtal strips
func splitImage(img *image.RGBA, n int) []image.Image {
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

type Monitor struct {
	Channels []chan bool
	DonePix  int
	MaxPix   int
}

func (mon *Monitor) Progress() (float64, bool) {
	done := true
	for i := range mon.Channels {
		_, open := <-mon.Channels[i]
		if open {
			done = false
			mon.DonePix++
		}
	}
	return float64(mon.DonePix) / float64(mon.MaxPix), done
}

func NewMonitor() *Monitor {
	return new(Monitor)
}

func Render(view Viewport, gen Generator, palette *Palette,
	monitor *Monitor, multisample, threads int) *image.RGBA {
	palette.generate()

	output_img := image.NewRGBA(
		image.Rect(0, 0, view.Width, view.Height))

	monitor.Channels = make([]chan bool, threads)
	monitor.MaxPix = view.Width * view.Height
	sub_images := splitImage(output_img, threads)

	for i := range monitor.Channels {
		monitor.Channels[i] = make(chan bool, 1024)
		m := float64(multisample)
		go func(img *image.RGBA, channel chan bool) {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {

					// sample for each sub-pixel
					var R, G, B uint
					for sx := 0; sx < multisample; sx++ {
						for sy := 0; sy < multisample; sy++ {
							x0 := float64(x)+float64(sx) / m
							y0 := float64(x)+float64(sy) / m
							itr := gen.EscapeAt(view.ComplexAt(x0, y0))
							r, g, b, _ := palette.ColorAt(itr).RGBA()
							R += uint(uint8(r))
							G += uint(uint8(g))
							B += uint(uint8(b))
						}
					}

					R /= uint(multisample * multisample)
					G /= uint(multisample * multisample)
					B /= uint(multisample * multisample)

					img.Set(x, y, color.RGBA{uint8(R), uint8(G), uint8(B), 255})
					channel <- true
				}
			}
			close(channel)
		}(sub_images[i].(*image.RGBA), monitor.Channels[i])
	}

	return output_img
}
