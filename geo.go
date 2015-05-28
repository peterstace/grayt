package grayt

// Material describes physical properties of a surface.
type Material struct {
	Colour Colour
}

// Intersection between some geometry and a ray.
type Intersection struct {
	Distance   float64 // Distance along the ray where the intersection occurred.
	UnitNormal Vect    // Unit normal (pointing 'away' from the geometry, not 'into' it).
	Material   Material
}

// Geometry implementations represent surfaces that can be intersected with.
type Geometry interface {

	// Intersect finds the intersection (if it exists) between a ray and the
	// geometry.
	Intersect(Ray) (Intersection, bool)
}
