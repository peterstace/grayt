package tracer

import (
	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

type plane struct {
	unitNormal vect.V // Unit normal out of the plane.
	anchor     vect.V // Any point on the plane.
}

// NewPlane creates a plane given a normal to the plane and an anchor point on
// the plane. The normal vector doesn't need to be a unit vector.
func NewPlane(normal, anchor vect.V) Geometry {
	return &plane{unitNormal: normal.Unit(), anchor: anchor}
}

func (p *plane) intersect(r ray.Ray) (hitRec, bool) {
	t := vect.Dot(p.unitNormal, vect.Sub(p.anchor, r.Start)) /
		vect.Dot(p.unitNormal, r.Dir)
	return hitRec{t: t, n: unitNormal}, t > 0
}
