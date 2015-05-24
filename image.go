package grayt

import (
	"image"
	"image/color"
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

func (a Accumulator) ToImage() image.Image {
	img := image.NewGray(image.Rect(0, 0, a.wide, a.high))
	for pxX := 0; pxX < a.wide; pxX++ {
		for pxY := 0; pxY < a.high; pxY++ {
			var colourUint8 uint8
			colour := a.get(pxX, pxY)
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
