package grayt

func newAccelerationStructure(objs ObjectList) accelerationStructure {
	return accelerationStructure{objs}
}

type accelerationStructure struct {
	objs []Object
}

func (a accelerationStructure) closestHit(r ray) (intersection, material, bool) {
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
