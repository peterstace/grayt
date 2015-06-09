package grayt

import "math"

func NewPlane(normal, anchor Vect) Surface {
	switch {
	case normal.Y == 0 && normal.Z == 0:
		return &alignXPlane{x: anchor.X, d: math.Copysign(1, normal.X)}
	case normal.Z == 0 && normal.X == 0:
		return &alignYPlane{y: anchor.Y, d: math.Copysign(1, normal.Y)}
	case normal.Y == 0 && normal.Y == 0:
		return &alignZPlane{z: anchor.Z, d: math.Copysign(1, normal.Z)}
	default:
		return &plane{
			unitNormal: normal.Unit(),
			anchor:     anchor,
		}
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

type alignXPlane struct {
	x float64
	d float64
}

func (p *alignXPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.x - r.Start.X) / r.Dir.X
	return Intersection{Vect{p.d, 0, 0}, t}, t > 0
}

type alignYPlane struct {
	y float64
	d float64
}

func (p *alignYPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.y - r.Start.Y) / r.Dir.Y
	return Intersection{Vect{0, p.d, 0}, t}, t > 0
}

type alignZPlane struct {
	z float64
	d float64
}

func (p *alignZPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.z - r.Start.Z) / r.Dir.Z
	return Intersection{Vect{0, 0, p.d}, t}, t > 0
}
