package grayt

import "math"

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

func newFastAccelerationStructure(objs ObjectList) accelerationStructure {
	return fastAccelerationStructure{} // TODO
}

type fastAccelerationStructure struct {
	// TODO

	// Bounding Area
	// Child structures.
}

func (a fastAccelerationStructure) closestHit(r ray) (intersection, material, bool) {
	return intersection{}, material{}, false // TODO
}

type boundingArea struct {
	min, max Vector
}

func (b *boundingArea) hit(r ray) bool {

	tx1 := (b.min.X - r.start.X) / r.dir.X
	tx2 := (b.max.X - r.start.X) / r.dir.X
	ty1 := (b.min.Y - r.start.Y) / r.dir.Y
	ty2 := (b.max.Y - r.start.Y) / r.dir.Y
	tz1 := (b.min.Z - r.start.Z) / r.dir.Z
	tz2 := (b.max.Z - r.start.Z) / r.dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)

	if math.Min(tx1, tx2) > tmin {
		tmin = math.Min(tx1, tx2)
	}
	if math.Max(tx1, tx2) < tmax {
		tmax = math.Max(tx1, tx2)
	}

	if math.Min(ty1, ty2) > tmin {
		tmin = math.Min(ty1, ty2)
	}
	if math.Max(ty1, ty2) < tmax {
		tmax = math.Max(ty1, ty2)
	}

	if math.Min(tz1, tz2) > tmin {
		tmin = math.Min(tz1, tz2)
	}
	if math.Max(tz1, tz2) < tmax {
		tmax = math.Max(tz1, tz2)
	}

	return tmin <= tmax && tmax > 0
}
