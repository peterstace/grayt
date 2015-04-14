package tracer

type hitRec struct {
	distance   float64
	unitNormal Vect
}

type Geometry interface {
	intersect(r Ray) (hitRec, bool)
}
