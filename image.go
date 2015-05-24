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

func (a Accumulator) distribution() (mean, stddev float64) {
	for _, v := range a.acc {
		mean += v
	}
	mean /= float64(len(a.acc))
	for _, v := range a.acc {
		stddev += (mean - v) * (mean - v)
	}
	stddev /= float64(len(a.acc))
	stddev = math.Sqrt(stddev)
	return
}

func (a Accumulator) ToImage() image.Image {

	// Number of stdandard deviations between the mean and the minimum/maiximum
	// pixel intensity. In the diagram below, this value is x.
	//
	//    WHITE        _         BLACK
	//      |       ,./ \.,        |
	//      | ,-----       ------, |
	//      +----------+-----------+
	// mean-x*stddev  mean  mean+x*stddev

	const numStdDevs = 3

	mean, stddev := a.distribution()

	img := image.NewGray(image.Rect(0, 0, a.wide, a.high))
	for pxX := 0; pxX < a.wide; pxX++ {
		for pxY := 0; pxY < a.high; pxY++ {

			colour := a.get(pxX, pxY)
			colour -= mean
			colour /= stddev * numStdDevs
			colour += 1.0
			colour /= 2.0

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
