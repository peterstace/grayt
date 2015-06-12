package grayt

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
	return []Surface{NewPlane(p.N, p.X)} // XXX
}

// XXX Don't need this factory any more.
func NewPlane(n, x Vect) Surface {
	switch {
	case n.Y == 0 && n.Z == 0:
		return &alignXPlane{x: x.X, d: math.Copysign(1, n.X)}
	case n.Z == 0 && n.X == 0:
		return &alignYPlane{y: x.Y, d: math.Copysign(1, n.Y)}
	case n.Y == 0 && n.Y == 0:
		return &alignZPlane{z: x.Z, d: math.Copysign(1, n.Z)}
	default:
		return &plane{
			n: n.Unit(),
			x: x,
		}
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
