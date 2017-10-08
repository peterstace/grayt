package engine

import (
	"image"
	"image/color"
	"math"
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

func (a *accumulator) add(x, y int, c vect3, index int) {
	i := y*a.wide + x
	a.pixels[i].mu.Lock()
	a.pixels[i].add(c, index)
	a.pixels[i].mu.Unlock()
}

func (a *accumulator) get(x, y int) vect3 {
	i := y*a.wide + x
	return a.pixels[i].colourSum
}

func (a *accumulator) mean() float64 {
	var sum float64
	for _, c := range a.pixels {
		c := c.colourSum
		sum += c[0] + c[1] + c[2]
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
			c := a.get(x, y).scale(0.5 * exposure / mean)
			r := math.Pow(c[0], 1.0/gamma)
			g := math.Pow(c[1], 1.0/gamma)
			b := math.Pow(c[2], 1.0/gamma)
			img.Set(x, y, nrgba(r, g, b))
		}
	}
	return img
}

func nrgba(r, g, b float64) color.NRGBA {
	return color.NRGBA{
		R: float64ToUint8(r),
		G: float64ToUint8(g),
		B: float64ToUint8(b),
		A: 0xff,
	}
}

func float64ToUint8(f float64) uint8 {
	switch {
	case f >= 1.0:
		return 0xff
	case f < 0.0:
		return 0x00
	default:
		// Since f >= 0.0 and f < 1.0, this returns a value beween 0x00 and
		// 0xff (inclusive).
		return uint8(f * 0x100)
	}
}

type pixel struct {
	mu        sync.Mutex
	colourSum vect3
	nextIndex int
	pending   []pendingColour
}

type pendingColour struct {
	colour vect3
	index  int
}

func (e *pixel) add(c vect3, index int) {

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
