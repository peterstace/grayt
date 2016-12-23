package grayt

type accelerationStructure interface {
	closestHit(ray) (intersection, material, bool)
}

func newListAccelerationStructure(objs ObjectList) accelerationStructure {
	return listAccelerationStructure{objs}
}

type listAccelerationStructure struct {
	objs []Object
}

func (a listAccelerationStructure) closestHit(r ray) (intersection, material, bool) {
	var closest struct {
		intersection intersection
		material     material
		hit          bool
	}
	for i := range a.objs {
		intersection, hit := a.objs[i].intersect(r)
		if !hit {
			continue
		}
		if !closest.hit || intersection.distance < closest.intersection.distance {
			closest.intersection = intersection
			closest.material = a.objs[i].material
			closest.hit = true
		}
	}
	return closest.intersection, closest.material, closest.hit
}
