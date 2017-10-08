package engine

type accelList struct {
	objs []object
}

func (a accelList) closestHit(e, d vect3) (n, h vect3, illum float64, ok bool) {
	var closest struct {
		n, h   vect3
		illum  float64
		distSq float64
		ok     bool
	}
	for i := range a.objs {

		n, h, ok := a.objs[i].surf.intersect(e, d)
		if !ok {
			continue
		}

		distSq := h.sub(e).norm2()

		if !closest.ok || distSq < closest.distSq {
			closest.n = n
			closest.h = h
			closest.illum = a.objs[i].illum
			closest.distSq = distSq
			closest.ok = true
		}
	}
	return closest.n, closest.h, closest.illum, closest.ok
}
