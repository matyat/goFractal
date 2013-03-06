package fractal

import (
	"image"
	"image/color"
)

// Split an image into a number of virtal strips
func splitImage(img *Paletted_uint32, n int) []image.Image {
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

func Render(view Viewport, gen Generator, palette color.Palette,
	monitor *Monitor, multisampling, threads int) *Paletted_uint32 {

	output_img := NewPaletted_uint32(
		image.Rect(0, 0, view.Width, view.Height), palette)

	monitor.Channels = make([]chan bool, threads)
	monitor.MaxPix = len(output_img.Pix)
	sub_images := splitImage(output_img, threads)

	for i := range monitor.Channels {
		monitor.Channels[i] = make(chan bool, 512)

		go func(img *Paletted_uint32, view Viewport,
			gen Generator, channel chan bool) {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
					img.SetColorIndex(x, y, gen.EscapeAt(view.ComplexAt(x, y)))
					channel <- true
				}
			}
			close(channel)
		}(sub_images[i].(*Paletted_uint32), view, gen, monitor.Channels[i])
	}

	return output_img
}
