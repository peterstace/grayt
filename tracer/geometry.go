package tracer

import "github.com/peterstace/grayt/ray"

type hitRec struct {
}

type Geometry interface {
	intersect(r ray.Ray) (hitRec, bool)
}
