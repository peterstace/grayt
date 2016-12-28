package grayt

import (
	"fmt"
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

func newFastAccelerationStructure(objs ObjectList) accelerationStructure {
	return newNode(objs)
}

type node struct {
	bound boundingArea

	// Either children or obj populated.
	child1, child2 *node
	obj            Object
}

func (n *node) String() string {
	return n.toString("")
}

func (n *node) toString(indent string) string {
	str := fmt.Sprintf(
		indent+"bound: %v\n",
		n.bound,
	)
	assert((n.child1 == nil) == (n.child2 == nil))
	if n.child1 != nil {
		str += n.child1.toString(indent + "\t")
		str += n.child2.toString(indent + "\t")
	} else {
		str += indent + "\t" + n.obj.String() + "\n"
	}
	return str
}

func newNode(objs []Object) *node {

	if len(objs) == 1 {
		min, max := objs[0].bound()
		return &node{
			bound: boundingArea{min, max},
			obj:   objs[0],
		}
	}

	n := len(objs)

	dim := func(v Vector) float64 { return v.X }
	var byCenter byCenter
	byCenter.objs = objs
	byCenter.dimension = dim
	sort.Sort(byCenter)

	inf := math.Inf(+1)
	bound := boundingArea{
		Vect(+inf, +inf, +inf),
		Vect(-inf, -inf, -inf),
	}
	for _, obj := range objs {
		min, max := obj.bound()
		bound.min = bound.min.Min(min)
		bound.max = bound.max.Max(max)
	}

	return &node{
		bound:  bound,
		child1: newNode(objs[:n/2]),
		child2: newNode(objs[n/2:]),
	}
}

type byCenter struct {
	dimension func(Vector) float64
	objs      []Object
}

func (b byCenter) Len() int      { return len(b.objs) }
func (b byCenter) Swap(i, j int) { b.objs[i], b.objs[j] = b.objs[j], b.objs[i] }
func (b byCenter) Less(i, j int) bool {
	minI, maxI := b.objs[i].bound()
	minJ, maxJ := b.objs[j].bound()
	return b.dimension(mid(minI, maxI)) < b.dimension(mid(minJ, maxJ))
}

func mid(u, v Vector) Vector {
	return u.Add(v).Scale(0.5)
}

func (a *node) closestHit(r ray) (intersection, material, bool) {

	if !a.bound.hit(r) {
		return intersection{}, material{}, false
	}

	assert((a.child1 == nil) == (a.child2 == nil))
	if a.child1 == nil {
		intersection, hit := a.obj.intersect(r)
		return intersection, a.obj.material, hit
	}

	intersection1, material1, hit1 := a.child1.closestHit(r)
	intersection2, material2, hit2 := a.child2.closestHit(r)
	switch {
	case hit1 && hit2:
		if intersection1.distance < intersection2.distance {
			return intersection1, material1, true
		} else {
			return intersection2, material2, true
		}
	case hit1:
		return intersection1, material1, true
	case hit2:
		return intersection2, material2, true
	default:
		return intersection{}, material{}, false
	}
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
