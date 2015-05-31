package grayt

import (
	"math"
	"math/rand"
)

func NewPlane(normal, anchor Vect) Surface {
	return &plane{
		unitNormal: normal.Unit(),
		anchor:     anchor,
	}
}

type plane struct {
	unitNormal Vect // Unit normal out of the plane.
	anchor     Vect // Any point on the plane.
}

func (p *plane) Intersect(r Ray) (Intersection, bool) {
	t := p.unitNormal.Dot(p.anchor.Sub(r.Start)) / p.unitNormal.Dot(r.Dir)
	return Intersection{p.unitNormal, t}, t > 0
}

func (p *plane) BoundingBox() (min, max Vect) {
	inf := math.Inf(1)
	max = Vect{inf, inf, inf}
	min = max.Extended(-1)
	return
}

func (p *plane) Sample() Vect {
	// No uniform distribution exists on the plane. So just find a normally
	// distributed random vector, and project it onto the plane.
	rnd := Vect{rand.NormFloat64(), rand.NormFloat64(), rand.NormFloat64()}.Add(p.anchor)
	return rnd.Sub(p.unitNormal.Extended(rnd.Dot(p.unitNormal)))
}
