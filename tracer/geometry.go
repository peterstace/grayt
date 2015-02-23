package tracer

import (
	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

type hitRec struct {
	t float64 // The distance needed to extend the ray to the hit site.
	n vect.V  // The normal at the hit site.
}

type Geometry interface {
	intersect(r ray.Ray) (hitRec, bool)
}
