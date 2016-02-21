package graytlib

import (
	"encoding/json"
	"math"
)

const planeT = "plane"

type Plane struct {
	N Vect
	X Vect
}

func (p Plane) MarshalJSON() ([]byte, error) {
	type alias Plane
	return json.Marshal(struct {
		Type string
		alias
	}{planeT, alias(p)})
}

func (p Plane) MakeSurfaces() []Surface {
	switch {
	case p.N.Y == 0 && p.N.Z == 0:
		return []Surface{&alignXPlane{x: p.X.X, d: math.Copysign(1, p.N.X)}}
	case p.N.Z == 0 && p.N.X == 0:
		return []Surface{&alignYPlane{y: p.X.Y, d: math.Copysign(1, p.N.Y)}}
	case p.N.Y == 0 && p.N.Y == 0:
		return []Surface{&alignZPlane{z: p.X.Z, d: math.Copysign(1, p.N.Z)}}
	default:
		return []Surface{&plane{n: p.N.Unit(), x: p.X}}
	}
}

type plane struct {
	n Vect // Unit normal out of the plane.
	x Vect // Any point on the plane.
}

func (p *plane) Intersect(r Ray) (Intersection, bool) {
	t := p.n.Dot(p.x.Sub(r.Start)) / p.n.Dot(r.Dir)
	return Intersection{p.n, t}, t > 0
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
