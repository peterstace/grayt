package grayt

import "image"

type accumulator struct {
	acc  []Colour
	wide int
	high int
}

func newAccumulator(wide, high int) accumulator {
	return accumulator{
		acc:  make([]Colour, wide*high),
		wide: wide,
		high: high,
	}
}

func (a accumulator) add(x, y int, c Colour) {
	i := y*a.wide + x
	a.acc[i] = a.acc[i].Add(c)
}

func (a accumulator) get(x, y int) Colour {
	i := y*a.wide + x
	return a.acc[i]
}

func (a accumulator) mean() float64 {
	var sum float64
	for _, c := range a.acc {
		sum += c.R + c.G + c.B
	}
	return sum / float64(len(a.acc)) / 3.0
}

func (a accumulator) dimensions() (pxHigh, pxWide int) {
	return a.wide, a.high
}

// ToImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a accumulator) toImage(exposure float64) image.Image {
	const gamma = 2.2
	mean := a.mean()
	img := image.NewNRGBA(image.Rect(0, 0, a.wide, a.high))
	for x := 0; x < a.wide; x++ {
		for y := 0; y < a.high; y++ {
			img.Set(x, y, a.get(x, y).
				Scale(0.5*exposure/mean).
				Pow(1.0/gamma).
				ToNRGBA())
		}
	}
	return img
}

func sum(c1, c2, c3, c4 Colour) Colour {
	return c1.Add(c2).Add(c3).Add(c4)
}
