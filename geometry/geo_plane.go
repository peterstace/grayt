package geometry

import (
	"math"

	"github.com/peterstace/grayt"
)

type Plane struct {
	N grayt.Vector
	X grayt.Vector
}

func (p Plane) MakeSurfaces() []grayt.Surface {
	switch {
	case p.N.Y == 0 && p.N.Z == 0:
		return []grayt.Surface{&alignXPlane{x: p.X.X, d: math.Copysign(1, p.N.X)}}
	case p.N.Z == 0 && p.N.X == 0:
		return []grayt.Surface{&alignYPlane{y: p.X.Y, d: math.Copysign(1, p.N.Y)}}
	case p.N.X == 0 && p.N.Y == 0:
		return []grayt.Surface{&alignZPlane{z: p.X.Z, d: math.Copysign(1, p.N.Z)}}
	default:
		return []grayt.Surface{&plane{n: p.N.Unit(), x: p.X}}
	}
}

type plane struct {
	n grayt.Vector // Unit normal out of the plane.
	x grayt.Vector // Any point on the plane.
}

func (p *plane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := p.n.Dot(p.x.Sub(r.Start)) / p.n.Dot(r.Dir)
	return grayt.Intersection{p.n, t}, t > 0
}

type alignXPlane struct {
	x float64
	d float64
}

func (p *alignXPlane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := (p.x - r.Start.X) / r.Dir.X
	return grayt.Intersection{grayt.Vector{p.d, 0, 0}, t}, t > 0
}

type alignYPlane struct {
	y float64
	d float64
}

func (p *alignYPlane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := (p.y - r.Start.Y) / r.Dir.Y
	return grayt.Intersection{grayt.Vector{0, p.d, 0}, t}, t > 0
}

type alignZPlane struct {
	z float64
	d float64
}

func (p *alignZPlane) Intersect(r grayt.Ray) (grayt.Intersection, bool) {
	t := (p.z - r.Start.Z) / r.Dir.Z
	return grayt.Intersection{grayt.Vector{0, 0, p.d}, t}, t > 0
}
