package tracer

import "math/rand"

type Light struct {
	Location  Vect
	Radius    float64
	Intensity float64
}

func (l *Light) sampleLocation() Vect {
	var offset Vect
	for {
		offset.X = rand.Float64()*2.0 - 1.0
		offset.Y = rand.Float64()*2.0 - 1.0
		offset.Z = rand.Float64()*2.0 - 1.0
		if offset.Length2() < 1.0 {
			break
		}
	}
	return l.Location.Add(offset.Extended(l.Radius))
}
