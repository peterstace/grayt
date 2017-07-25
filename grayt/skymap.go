package grayt

import (
	"errors"
	"image"
	"math"
)

func CreateSkymap(img image.Image) (Skymap, error) {

	sz := img.Bounds().Size()
	if sz.X != 2*sz.Y {
		return Skymap{}, errors.New("invalid size: width must be twice height")
	}

	s := Skymap{make([][]Colour, sz.X)}
	for x := range s.data {
		s.data[x] = make([]Colour, sz.Y)
	}

	for x := 0; x < sz.X; x++ {
		for y := 0; y < sz.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			s.data[x][y] = Colour{
				float64(r) / float64(a),
				float64(g) / float64(a),
				float64(b) / float64(a),
			}
		}
	}
	return s, nil
}

type Skymap struct {
	data [][]Colour
}

func (s Skymap) intersect(dir Vector) Colour {

	if s.data == nil {
		return Colour{}
	}

	dir = dir.Unit()

	// Between -pi and +pi
	azimuthAngle := math.Atan2(-dir.Z, -dir.X)

	x := int((azimuthAngle + math.Pi) / (2 * math.Pi) * float64(len(s.data)))
	x = intMax(0, intMin(x, len(s.data)-1))

	// Between +pi/2 and -pi/2
	altitudeAngle := math.Asin(dir.Y)

	y := int((altitudeAngle + math.Pi/2) / math.Pi * float64(len(s.data)) / 2)
	y = len(s.data)/2 - 1 - intMax(0, intMin(y, len(s.data)/2-1))

	return s.data[x][y]
}
