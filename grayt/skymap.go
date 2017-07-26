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

	x := int(azimuth(dir) * float64(len(s.data)))
	x = intMax(0, intMin(x, len(s.data)-1))

	y := int(altitude(dir) * float64(len(s.data)/2))
	y = intMax(0, intMin(y, len(s.data)/2-1))

	return s.data[x][y]
}

// Calculates the azimuth as a fraction (0 <= azimuth < 1).
func azimuth(unitDir Vector) float64 {
	azimuthAngle := 1.5*math.Pi - math.Atan2(-unitDir.Z, unitDir.X) // 0.5*pi to 2.5*pi
	if azimuthAngle > 2*math.Pi {
		azimuthAngle -= 2 * math.Pi // 0 to 2*pi
	}
	return azimuthAngle / (2 * math.Pi)
}

// Calculates the altitude as a fraction (0 <= alt < 1)
func altitude(unitDir Vector) float64 {
	alt := math.Asin(unitDir.Y) // -pi/2 to +pi/2
	return 1.0 - (alt+math.Pi/2)/math.Pi
}
