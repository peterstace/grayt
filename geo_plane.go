package grayt

import "math"

func Plane(unitNormal, pointOnPlane Vector) Surface {
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
	n Vector // Unit normal out of the plane.
	x Vector // Any point on the plane.
}

func (p *plane) Intersect(r Ray) (Intersection, bool) {
	t := p.n.Dot(p.x.Sub(r.Start)) / p.n.Dot(r.Dir)
	return Intersection{UnitNormal: p.n, Distance: t}, t > 0
}

type alignXPlane struct {
	x float64
	d float64
}

func (p *alignXPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.x - r.Start.X) / r.Dir.X
	return Intersection{UnitNormal: Vect(p.d, 0, 0), Distance: t}, t > 0
}

type alignYPlane struct {
	y float64
	d float64
}

func (p *alignYPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.y - r.Start.Y) / r.Dir.Y
	return Intersection{UnitNormal: Vect(0, p.d, 0), Distance: t}, t > 0
}

type alignZPlane struct {
	z float64
	d float64
}

func (p *alignZPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.z - r.Start.Z) / r.Dir.Z
	return Intersection{UnitNormal: Vect(0, 0, p.d), Distance: t}, t > 0
}
