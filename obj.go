package grayt

type Surface interface {

	// HitSurface checks for a surface hit. If the surface is in front, the
	// first return value is the ray multiplier required to hit the surface. If
	// the surface is behind or cannot be hit, the first return value is
	// negative. The second return value is the unit normal at the surface hit.
	HitSurface(ray) (float64, Vect)
}
