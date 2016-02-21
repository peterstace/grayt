package grayt

import "image"

type Accumulator struct {
	acc  []Colour
	wide int
	high int
}

func NewAccumulator(wide, high int) Accumulator {
	return Accumulator{
		acc:  make([]Colour, wide*high),
		wide: wide,
		high: high,
	}
}

func (a Accumulator) add(x, y int, c Colour) {
	i := y*a.wide + x
	a.acc[i] = a.acc[i].Add(c)
}

func (a Accumulator) get(x, y int) Colour {
	i := y*a.wide + x
	return a.acc[i]
}

func (a Accumulator) mean() float64 {
	var sum float64
	for _, c := range a.acc {
		sum += c.R + c.G + c.B
	}
	return sum / float64(len(a.acc)) / 3.0
}

func (a Accumulator) Dimensions() (pxHigh, pxWide int) {
	return a.wide, a.high
}

// ToImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a Accumulator) ToImage(exposure float64) image.Image {
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

func (a Accumulator) NeighbourCoefficientOfVariation() float64 {
	var mean float64
	for x := 1; x < a.wide; x++ {
		for y := 1; y < a.high; y++ {

			h, i, j, k := a.get(x, y), a.get(x-1, y), a.get(x, y-1), a.get(x-1, y-1)

			m := sum(h, i, j, k).Scale(-0.25)
			v := sum(
				m.Add(h).Square(),
				m.Add(i).Square(),
				m.Add(j).Square(),
				m.Add(k).Square(),
			).Scale(0.25)

			if m.R == 0 || m.G == 0 || m.B == 0 {
				continue
			}

			cv := v.Pow(0.5).Div(m.Scale(-1))

			mean += (cv.R + cv.G + cv.B) / 3
		}
	}
	return mean / float64((a.wide-1)*(a.high-1))
}

func sum(c1, c2, c3, c4 Colour) Colour {
	return c1.Add(c2).Add(c3).Add(c4)
}
