package main

func newAccelerationStructure(ts []triangle) accelerationStructure {
	return accelerationStructure{ts}
}

type accelerationStructure struct {
	tris []triangle
}

func (a accelerationStructure) closestHit(r ray) (intersection, bool) {
	var closest struct {
		intersection intersection
		hit          bool
	}
	for i := range a.tris {
		intersection, hit := a.tris[i].intersect(r)
		if !hit {
			continue
		}
		if !closest.hit || intersection.distance < closest.intersection.distance {
			closest.intersection = intersection
			closest.hit = true
		}
	}
	return closest.intersection, closest.hit
}
