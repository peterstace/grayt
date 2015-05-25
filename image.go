package grayt

import (
	"image"
	"image/color"
	"math"
)

type Accumulator struct {
	acc  []float64
	wide int
	high int
}

func NewAccumulator(wide, high int) Accumulator {
	return Accumulator{
		acc:  make([]float64, wide*high),
		wide: wide,
		high: high,
	}
}

func (a Accumulator) add(x, y int, v float64) {
	a.acc[y*a.wide+x] += v
}

func (a Accumulator) get(x, y int) float64 {
	return a.acc[y*a.wide+x]
}

func (a Accumulator) mean() float64 {
	var sum float64
	for _, v := range a.acc {
		sum += v
	}
	return sum / float64(len(a.acc))
}

// ToImage converts the accumulator into an image. Exposure controls how bright
// the arithmetic mean brightness in the image is. A value of 1.0 results in a
// mean brightness half way between black and white.
func (a Accumulator) ToImage(exposure float64) image.Image {

	const gamma = 2.2

	mean := a.mean()

	img := image.NewGray(image.Rect(0, 0, a.wide, a.high))
	for pxX := 0; pxX < a.wide; pxX++ {
		for pxY := 0; pxY < a.high; pxY++ {

			colour := a.get(pxX, pxY)
			colour /= mean
			colour *= 0.5 * exposure

			colour = math.Pow(colour, 1.0/gamma)

			var colourUint8 uint8
			if colour >= 1.0 {
				colourUint8 = 0xff
			} else if colour < 0.0 {
				colourUint8 = 0x00
			} else {
				colourUint8 = uint8(colour * 256)
			}
			img.Set(pxX, pxY, color.Gray{Y: colourUint8})
		}
	}
	return img
}
