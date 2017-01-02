package grayt

import (
	"image"
	"sync"
)

type accumulator struct {
	pixels []pixel
	wide   int
	high   int
}

func newAccumulator(wide, high int) *accumulator {
	return &accumulator{
		pixels: make([]pixel, wide*high),
		wide:   wide,
		high:   high,
	}
}

func (a *accumulator) add(x, y int, c Colour, index int) {
	i := y*a.wide + x
	a.pixels[i].mu.Lock()
	a.pixels[i].add(c, index)
	a.pixels[i].mu.Unlock()
}

func (a *accumulator) get(x, y int) Colour {
	i := y*a.wide + x
	return a.pixels[i].colourSum
}

func (a *accumulator) mean() float64 {
	var sum float64
	for _, c := range a.pixels {
		c := c.colourSum
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
			img.Set(x, y, a.get(x, y).
				scale(0.5*exposure/mean).
				pow(1.0/gamma).
				toNRGBA())
		}
	}
	return img
}

type pixel struct {
	mu        sync.Mutex
	colourSum Colour
	nextIndex int
	pending   []pendingColour
}

type pendingColour struct {
	colour Colour
	index  int
}

func (e *pixel) add(c Colour, index int) {

	if e.nextIndex != index {
		e.pending = append(e.pending, pendingColour{c, index})
		return
	}

	e.colourSum = e.colourSum.add(c)
	e.nextIndex++

	for true {
		found := false
		for i := range e.pending {
			if e.pending[i].index == e.nextIndex {
				found = true
				e.colourSum = e.colourSum.add(e.pending[i].colour)
				e.nextIndex++
				e.pending[i] = e.pending[len(e.pending)-1]
				e.pending = e.pending[:len(e.pending)-1]
				break
			}
		}
		if !found {
			break
		}
	}
}
