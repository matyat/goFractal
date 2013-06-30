package fractal

import (
	"image"
	"image/color"
	"sync/atomic"
)

// Stuct used to monitor the progress of the threads
type Monitor struct {
	Channels []chan bool
	DonePix  int
	MaxPix   int
}

// Get the current progress. Returns the number of processed pixel, and
// a whether the render is finished or not. Warning, not calling this often
// will cause the render threads to block.
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

// Create a new Monitor
func NewMonitor() *Monitor {
	return new(Monitor)
}

// Begins a fractal render and returns an image. The image should not be accessed
// until the monitor reports the render has ended.
func Render(view Viewport, gen Generator, color_wheel *ColorWheel,
	monitor *Monitor, multisample, threads int) *image.RGBA {

	color_wheel.generate()

	output_img := image.NewRGBA(image.Rect(0, 0, view.Width, view.Height))

	monitor.Channels = make([]chan bool, threads)
	monitor.MaxPix = view.Width * view.Height

	// an array of values to make sure each pixel is only rendered once
	pixel_lock := make([]int64, monitor.MaxPix)

	for i := range monitor.Channels {
		// buffering these channels stops them from being locked when reporting their
		// progress
		monitor.Channels[i] = make(chan bool, 64)

		m := float64(multisample)

		go func(channel chan bool, id int) {
			for y := output_img.Rect.Min.Y; y < output_img.Rect.Max.Y; y++ {
				stride := y * view.Width
				for x := output_img.Rect.Min.X; x < output_img.Rect.Max.X; x++ {
					// if pixel_lock is assigned any other than 0, the pixel has
					// already been processed
					if atomic.CompareAndSwapInt64(
						&pixel_lock[x+stride], 0, int64(id)) {
						var R, G, B uint8

						//sample sub-pixels when multisampling
						for sx := 0; sx < multisample; sx++ {
							for sy := 0; sy < multisample; sy++ {
								x0 := float64(x) + float64(sx)/m
								y0 := float64(y) + float64(sy)/m

								itr := gen.EscapeAt(view.ComplexAt(x0, y0))
								r, g, b, _ := color_wheel.ColorAt(itr).RGBA()

								R += uint8(r)
								G += uint8(g)
								B += uint8(b)
							}
						}

						// Average each colour
						m_sqr := uint8(multisample * multisample)
						R /= m_sqr
						G /= m_sqr
						B /= m_sqr

						output_img.Set(x, y, color.RGBA{R, G, B, 255})
						channel <- true
					}
				}
			}
			close(channel)
		}(monitor.Channels[i], i+1)
	}

	return output_img
}
