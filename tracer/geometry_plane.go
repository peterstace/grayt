package tracer

type plane struct {
	unitNormal Vect // Unit normal out of the plane.
	anchor     Vect // Any point on the plane.
}

// NewPlane creates a plane given a normal to the plane and an anchor point on
// the plane. The normal vector doesn't need to be a unit vector.
func NewPlane(anchor, normal Vect) Geometry {
	return &plane{unitNormal: normal.Unit(), anchor: anchor}
}

func (p *plane) intersect(r Ray) (hitRec, bool) {
	t := p.unitNormal.Dot(p.anchor.Sub(r.Start)) /
		p.unitNormal.Dot(r.Dir)
	return hitRec{distance: t, unitNormal: p.unitNormal}, t > 0
}
