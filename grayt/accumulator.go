package grayt

import (
	"image"
)

type pixelGrid struct {
	pixels []Colour
	wide   int
	high   int
}

func (g *pixelGrid) set(x, y int, c Colour) {
	i := y*g.wide + x
	g.pixels[i] = c
}

type accumulator struct {
	pixelGrid
	count int
}

func (a *accumulator) merge(g *pixelGrid) {
	a.count++
	for i, c := range a.pixels {
		a.pixels[i] = c.add(g.pixels[i])
	}
}

func (a *accumulator) mean() float64 {
	var sum float64
	for _, c := range a.pixels {
		sum += c.R + c.G + c.B
	}
	return sum / float64(len(a.pixels)) / 3.0
}

// ToImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a *accumulator) toImage(exposure float64) image.Image {
	const gamma = 2.2
	mean := a.mean()
	img := image.NewNRGBA(image.Rect(0, 0, a.wide, a.high))
	for x := 0; x < a.wide; x++ {
		for y := 0; y < a.high; y++ {
			i := y*a.wide + x
			img.Set(x, y, a.pixels[i].
				scale(0.5*exposure/mean).
				pow(1.0/gamma).
				toNRGBA())
		}
	}
	return img
}

func (a *accumulator) load() error {
	// TODO
	return nil
}

func (a *accumulator) save() error {
	// TODO
	return nil
}
