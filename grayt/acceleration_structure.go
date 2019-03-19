package grayt

import "github.com/peterstace/grayt/xmath"

type accelerationStructure interface {
	closestHit(xmath.Ray) (intersection, material, bool)
}

func newListAccelerationStructure(objs ObjectList) accelerationStructure {
	return listAccelerationStructure{objs}
}

type listAccelerationStructure struct {
	objs []Object
}

func (a listAccelerationStructure) closestHit(r xmath.Ray) (intersection, material, bool) {
	var closest struct {
		intersection intersection
		material     material
		hit          bool
	}
	for i := range a.objs {
		intersection, hit := a.objs[i].Surface.intersect(r)
		if !hit {
			continue
		}
		if !closest.hit || intersection.distance < closest.intersection.distance {
			closest.intersection = intersection
			closest.material = a.objs[i].Material
			closest.hit = true
		}
	}
	return closest.intersection, closest.material, closest.hit
}
