package grayt

type Surface interface {

	// HitSurface checks for a surface hit. If the surface is in front, the
	// first return value is the ray multiplier required to hit the surface. If
	// the surface is behind or cannot be hit, the first return value is
	// negative. If there was a hit, then the second return value is the (not
	// necessarily unit) normal.
	HitSurface(ray) (float64, Vect)
}
