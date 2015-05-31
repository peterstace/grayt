package grayt

import "math"

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
