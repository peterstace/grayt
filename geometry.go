package grayt

import "math"

// Intersection between a surface and a ray.
type Intersection struct {
	UnitNormal Vect    // Unit normal (pointing 'away' from the surface, not 'into' it).
	Distance   float64 // Distance along the ray where the intersection occurred.
}

// Surface is a two dimensional surface that can be intersected with a ray.
type Surface interface {

	// Intersect finds the intersection (if it exists) between a ray and the
	// surface.
	Intersect(Ray) (Intersection, bool)
}

type Material struct {
	Colour    Colour
	Emittance float64 // XXX is this the right word?
	// Other properties such refractive index, reflectance, BRDF etc go here.
}

// Entity is a physical object whithin a scene.
type Entity struct {
	Surface  Surface
	Material Material
}

type Scene struct {
	Camera   Camera
	Entities []Entity
}

// Sphere

func NewSphere(centre Vect, radius float64) Surface {
	return &sphere{
		centre: centre,
		radius: radius,
	}
}

type sphere struct {
	centre Vect
	radius float64
}

func (s *sphere) Intersect(r Ray) (Intersection, bool) {

	// Get coeficients to a.x^2 + b.x + c = 0
	emc := r.Start.Sub(s.centre)
	a := r.Dir.Length2()
	b := 2 * emc.Dot(r.Dir)
	c := emc.Length2() - s.radius*s.radius

	// Find discrimenant b*b - 4*a*c
	disc := b*b - 4*a*c
	if disc < 0 {
		return Intersection{}, false
	}

	// Find x1 and x2 using a numerically stable algorithm.
	var signOfB float64
	signOfB = math.Copysign(1.0, b)
	q := -0.5 * (b + signOfB*math.Sqrt(disc))
	x1 := q / a
	x2 := c / q

	var t float64
	if x1 > 0 && x2 > 0 {
		// Both are positive, so take the smaller one.
		t = math.Min(x1, x2)
	} else {
		// At least one is negative, take the larger one (which is either
		// negative or positive).
		t = math.Max(x1, x2)
	}

	return Intersection{r.At(t).Sub(s.centre).Unit(), t}, t > 0
}

// Plane

func NewPlane(normal, anchor Vect) Surface {
	return &plane{
		unitNormal: normal.Unit(),
		anchor:     anchor,
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

// Triangle

func NewTriangle(cornerA, cornerB, cornerC Vect) Surface {
	u := cornerB.Sub(cornerA)
	v := cornerC.Sub(cornerA)
	t := &triangle{
		a:        cornerA,
		u:        u,
		v:        v,
		unitNorm: u.Cross(v).Unit(),
		dotUV:    u.Dot(v),
		dotUU:    u.Dot(u),
		dotVV:    v.Dot(v),
	}
	return t
}

type triangle struct {
	a, u, v             Vect // Corner A, A to B, and A to C.
	unitNorm            Vect
	dotUV, dotUU, dotVV float64 // Precomputed dot products.
}

func (t *triangle) Intersect(r Ray) (Intersection, bool) {

	// Check if there's a hit with the plane.
	h := t.unitNorm.Dot(t.a.Sub(r.Start)) / t.unitNorm.Dot(r.Dir)
	if h < 0 {
		// Hit was behind the camera.
		return Intersection{}, false
	}

	// Find out if the plane hit was inside the triangle. We need to solve the
	// equation w = alpha*u + beta*v for alpha and beta (where alpha beta
	// scalars, and u and v are vectors from a to b and a to c, and w is a
	// vector from a to the hit point).
	//
	// If the sum of alpha and beta is less than 1 and both alpha and beta are
	// positive, then the hit is inside the triangle.
	//
	// alpha = [(u.v)(w.v) - (v.v)(w.u)] / [(u.v)^2 - (u.u)(v.v)]
	// beta  = [(u.v)(w.u) - (u.u)(w.v)] / [(u.v)^2 - (u.u)(v.v)]

	w := r.At(h).Sub(t.a)
	dotWV := w.Dot(t.v)
	dotWU := w.Dot(t.u)
	alpha := t.dotUV*dotWV - t.dotVV*dotWU
	beta := t.dotUV*dotWU - t.dotUU*dotWV
	denom := t.dotUV*t.dotUV - t.dotUU*t.dotVV
	alpha /= denom
	beta /= denom

	if alpha < 0 || beta < 0 || alpha+beta > 1 {
		return Intersection{}, false
	}
	return Intersection{t.unitNorm, h}, true
}

// Square

func NewSquare(v1, v2, v3, v4 Vect) []Surface {
	return []Surface{
		NewTriangle(v1, v2, v3),
		NewTriangle(v3, v4, v1),
	}
}
