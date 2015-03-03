package tracer

import (
	"github.com/peterstace/grayt/ray"
	"github.com/peterstace/grayt/vect"
)

type hitRec struct {
	distance   float64
	unitNormal vect.V
}

type Geometry interface {
	intersect(r ray.Ray) (hitRec, bool)
}
