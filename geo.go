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

	// Bounding box finds the smallest axis-aligned box that bounds the
	// surface. Two opposing vertices of the box are returned, with the value
	// of a coordinate in min always being less than or equal to the
	// corresponding coordinate in max.
	BoundingBox() (min, max Vect)

	// Sample returns a uniformly distributed random point on the surface. If a
	// uniform distribution doesn't exist, then some other "well behaved but
	// undefined" distribution is used instead.
	Sample() Vect
}

type Emitter struct {
	Surface
	Colour    Colour
	Intensity float64
}

type Reflector struct {
	Surface
	Material Material
}

type Material struct {
	Colour Colour
	// Other properties e.g. refractive index, reflectance etc go here.
}

type Scene struct {
	Camera     Camera
	Emitters   []Emitter
	Reflectors []Reflector
}
