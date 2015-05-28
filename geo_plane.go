package grayt

func NewPlane(material Material, normal, anchor Vect) Geometry {
	return &plane{
		material:   material,
		unitNormal: normal.Unit(),
		anchor:     anchor,
	}
}

type plane struct {
	material   Material
	unitNormal Vect // Unit normal out of the plane.
	anchor     Vect // Any point on the plane.
}

func (p *plane) Intersect(r Ray) (Intersection, bool) {
	t := p.unitNormal.Dot(p.anchor.Sub(r.Start)) / p.unitNormal.Dot(r.Dir)
	return Intersection{
		Distance:   t,
		UnitNormal: p.unitNormal,
		Material:   p.material,
	}, t > 0
}
