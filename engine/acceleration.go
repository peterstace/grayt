package engine

type accelList struct {
	surfs []surface
}

func (a accelList) closestHit(e, d vect3) (n, h vect3, illum float64, ok bool) {
	var closest struct {
		n, h   vect3
		illum  float64
		distSq float64
		ok     bool
	}
	for i := range a.surfs {

		trans := a.surfs[i].transform
		transInv, ok := trans.inv()
		if !ok {
			panic("could not invert transformation matrix")
		}

		ePrime := transInv.mulv(e.extend(1)).trunc()
		dPrime := transInv.mulv(d.extend(0)).trunc()

		nPrime, hPrime, ok := triangle(ePrime, dPrime)
		if !ok {
			continue
		}

		n := transInv.transpose().mulv(nPrime.extend(0)).trunc()
		h := trans.mulv(hPrime.extend(1)).trunc()

		distSq := h.sub(e).norm2()

		if !closest.ok || distSq < closest.distSq {
			closest.n = n
			closest.h = h
			closest.illum = a.surfs[i].illumination
			closest.distSq = distSq
			closest.ok = true
		}
	}
	return closest.n, closest.h, closest.illum, closest.ok
}
