package grayt

func Plane(unitNormal Vector) Surface {
	return &plane{unitNormal, Vector{}}
}

type plane struct {
	n Vector // Unit normal out of the plane.
	x Vector // Any point on the plane.
}

func (p *plane) Translate(x, y, z float64) Surface {
	p.x = p.x.Add(Vect(x, y, z))
	return p
}

func (p *plane) Intersect(r Ray) (Intersection, bool) {
	t := p.n.Dot(p.x.Sub(r.Start)) / p.n.Dot(r.Dir)
	return Intersection{UnitNormal: p.n, Distance: t}, t > 0
}

func XPlane() Surface {
	return &alignXPlane{0}
}

type alignXPlane struct {
	x float64
}

func (p *alignXPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.x - r.Start.X) / r.Dir.X
	return Intersection{UnitNormal: Vect(+1, 0, 0), Distance: t}, t > 0
}

func (p *alignXPlane) Translate(x, y, z float64) Surface {
	p.x += x
	return p
}

func YPlane() Surface {
	return &alignYPlane{0}
}

type alignYPlane struct {
	y float64
}

func (p *alignYPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.y - r.Start.Y) / r.Dir.Y
	return Intersection{UnitNormal: Vect(0, +1, 0), Distance: t}, t > 0
}

func (p *alignYPlane) Translate(x, y, z float64) Surface {
	p.y += y
	return p
}

func ZPlane() Surface {
	return &alignZPlane{0}
}

type alignZPlane struct {
	z float64
}

func (p *alignZPlane) Intersect(r Ray) (Intersection, bool) {
	t := (p.z - r.Start.Z) / r.Dir.Z
	return Intersection{UnitNormal: Vect(0, 0, +1), Distance: t}, t > 0
}

func (p *alignZPlane) Translate(x, y, z float64) Surface {
	p.z += z
	return p
}
