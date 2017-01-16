package grayt

import (
	"image"
	"sort"
)

func createImage(a *accumulator) image.Image {

	intensities := make([]float64, len(a.pixels))
	for i, p := range a.pixels {
		intensities[i] = intensity(p.colourSum)
	}
	sort.Float64s(intensities)

	m := map[float64]float64{}
	var i int
	for i < len(intensities) {
		startI := i
		for i < len(intensities) && intensities[i] == intensities[startI] {
			i++
		}
		/*
			x - center
			min, max
			top

			c / (top - c) == (c - min) / (max - c)
			c == (top - c).(c - min)/(max - c)
			c.(max-c) == (top-c).(c-min)
			c.max - c*c == top.c - top.min - c*c + c.min
			c.max == top.c - top.min + c.min
			c.max - top.c - c.min == -top.min
			c.(max - top - min) == -top.min
			c == top.min / (top + min - max)
		*/

		top := float64(len(intensities))
		min := float64(startI)
		max := float64(i)
		c := top * min / (top + min - max)

		newIntensity := c / float64(len(intensities))
		oldIntensity := intensities[startI]
		m[oldIntensity] = newIntensity
	}

	const gamma = 2.2
	img := image.NewNRGBA(image.Rect(0, 0, a.wide, a.high))
	for x := 0; x < a.wide; x++ {
		for y := 0; y < a.high; y++ {
			c := a.get(x, y)
			oldIntensity := intensity(c)
			newIntensity, ok := m[oldIntensity]
			if !ok {
				panic(false)
			}
			c = c.scale(newIntensity / oldIntensity)
			c = c.pow(1.0 / gamma)
			img.Set(x, y, c.toNRGBA())
		}
	}
	return img
}

func intensity(c Colour) float64 {
	return 0.2126*c.R + 0.7152*c.G + 0.0722*c.B
}
