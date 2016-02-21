package geometry

import (
	"math"

	"github.com/peterstace/grayt"
)

func Plane(unitNormal, pointOnPlane grayt.Vector) grayt.Surface {
	switch {
	case unitNormal.Y == 0 && unitNormal.Z == 0:
		return &alignXPlane{x: pointOnPlane.X, d: math.Copysign(1, unitNormal.X)}
	case unitNormal.Z == 0 && unitNormal.X == 0:
		return &alignYPlane{y: pointOnPlane.Y, d: math.Copysign(1, unitNormal.Y)}
	case unitNormal.X == 0 && unitNormal.Y == 0:
		return &alignZPlane{z: pointOnPlane.Z, d: math.Copysign(1, unitNormal.Z)}
	default:
		return &plane{n: unitNormal.Unit(), x: pointOnPlane}
	}
}

type plane struct {
	n grayt.Vector // Unit normal out of the plane.
	x grayt.Vector // Any point on the plane.
}

func (p *plane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := p.n.Dot(p.x.Sub(r.Start)) / p.n.Dot(r.Dir)
	return grayt.Intersection{UnitNormal: p.n, Distance: t}, t > 0
}

type alignXPlane struct {
	x float64
	d float64
}

func (p *alignXPlane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := (p.x - r.Start.X) / r.Dir.X
	return grayt.Intersection{UnitNormal: grayt.Vect(p.d, 0, 0), Distance: t}, t > 0
}

type alignYPlane struct {
	y float64
	d float64
}

func (p *alignYPlane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := (p.y - r.Start.Y) / r.Dir.Y
	return grayt.Intersection{UnitNormal: grayt.Vect(0, p.d, 0), Distance: t}, t > 0
}

type alignZPlane struct {
	z float64
	d float64
}

func (p *alignZPlane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := (p.z - r.Start.Z) / r.Dir.Z
	return grayt.Intersection{UnitNormal: grayt.Vect(0, 0, p.d), Distance: t}, t > 0
}
