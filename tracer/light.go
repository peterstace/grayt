package tracer

import (
	"math/rand"

	"github.com/peterstace/grayt/vect"
)

type Light struct {
	Location  vect.V
	Radius    float64
	Intensity float64
}

func (l *Light) sampleLocation() vect.V {
	var offset vect.V
	for {
		offset.X = rand.Float64()*2.0 - 1.0
		offset.Y = rand.Float64()*2.0 - 1.0
		offset.Z = rand.Float64()*2.0 - 1.0
		if offset.Norm2() < 1.0 {
			break
		}
	}
	return vect.Add(l.Location, offset.Extended(l.Radius))
}
