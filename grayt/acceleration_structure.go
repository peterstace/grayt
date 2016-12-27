package grayt

import (
	"math"
	"sort"
)

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

func newNode(objs []Object) *node {

	n := len(objs)

	xmax := make([]float64, n)
	for i, obj := range objs {
		_, max := obj.bound()
		xmax[i] = max.X
	}
	sort.Float64s(xmax)
	xmin := make([]float64, n)
	for i, obj := range objs {
		min, _ := obj.bound()
		xmin[i] = min.X
	}
	sort.Float64s(xmin)

	ymax := make([]float64, n)
	for i, obj := range objs {
		_, max := obj.bound()
		ymax[i] = max.Y
	}
	sort.Float64s(ymax)
	ymin := make([]float64, n)
	for i, obj := range objs {
		min, _ := obj.bound()
		ymin[i] = min.Y
	}
	sort.Float64s(ymin)

	zmax := make([]float64, n)
	for i, obj := range objs {
		_, max := obj.bound()
		zmax[i] = max.Z
	}
	sort.Float64s(zmax)
	zmin := make([]float64, n)
	for i, obj := range objs {
		min, _ := obj.bound()
		zmin[i] = min.Z
	}
	sort.Float64s(zmin)

	bound := boundingArea{
		Vect(xmin[0], ymin[0], zmin[0]),
		Vect(xmax[n-1], ymax[n-1], zmax[n-1]),
	}

	switch n {
	case 1:
		return &node{
			bound: bound,
			obj:   objs[0],
		}
	case 2:
		return &node{
			bound: bound,
			children: []*node{
				newNode([]Object{objs[0]}),
				newNode([]Object{objs[1]}),
			},
		}
	default:
		children1 := []Object{}
		children2 := []Object{}
		cutoffXMin := xmin[n/2]
		cutoffXMax := xmax[n/2]
		for _, obj := range objs {
			min, max := obj.surface.bound()
			if max.X < cutoffXMax {
				children1 = append(children1, obj)
			}
			if min.X > cutoffXMin {
				children2 = append(children2, obj)
			}
		}

		return &node{
			bound: bound,
			children: []*node{
				newNode(children1),
				newNode(children2),
			},
		}
	}
}

func newFastAccelerationStructure(objs ObjectList) accelerationStructure {
	return newNode(objs)
}

type node struct {
	bound boundingArea

	// Exactly 1 field populated:
	children []*node
	obj      Object
}

func (a *node) closestHit(r ray) (intersection, material, bool) {
	if !a.bound.hit(r) {
		return intersection{}, material{}, false
	}
	if len(a.children) == 0 {
		intersection, hit := a.obj.intersect(r)
		return intersection, a.obj.material, hit
	}
	for _, child := range a.children {
		intersection, material, hit := child.closestHit(r)
		if hit {
			return intersection, material, true
		}
	}
	return intersection{}, material{}, false
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
