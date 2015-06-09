package grayt

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
	Emittance float64
	// Other properties such refractive index, reflectance, BRDF etc go here.
}

// Entity is a physical object whithin a scene.
type Entity struct {
	Surfaces []Surface
	Material Material
}

type Scene struct {
	Camera   Camera
	Entities []Entity
}
