package grayt

import (
	"image"
	"sort"
)

func createImage(a *accumulator) image.Image {

	intensityMap := buildIntensityMap(a)

	const gamma = 2.2
	img := image.NewNRGBA(image.Rect(0, 0, a.wide, a.high))
	for x := 0; x < a.wide; x++ {
		for y := 0; y < a.high; y++ {
			c := a.get(x, y)
			oldIntensity := intensity(c)
			newIntensity, ok := intensityMap[oldIntensity]
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

func buildIntensityMap(a *accumulator) map[float64]float64 {

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
		center := boxCenter(float64(len(intensities)), float64(startI), float64(i))
		newIntensity := center / float64(len(intensities))
		oldIntensity := intensities[startI]
		m[oldIntensity] = newIntensity
	}

	return m
}

func boxCenter(boxMin, boxMax, maxVal float64) float64 {
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

	return maxVal * boxMin / (maxVal + boxMin - boxMax)
}
