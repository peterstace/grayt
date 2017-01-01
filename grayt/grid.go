package grayt

import (
	"math"
)

type grid struct {
	minBound Vector
	maxBound Vector

	stride     Vector
	data       []*link
	resolution triple
}

func newGrid(lambda float64, objs ObjectList) *grid {

	minBound, maxBound := bounds(objs)
	boundDiff := maxBound.Sub(minBound)
	volume := boundDiff.X * boundDiff.Y * boundDiff.Z

	resolutionFactor := math.Pow(lambda*float64(len(objs))/volume, 1.0/3.0)
	resolution := truncate(boundDiff.Scale(resolutionFactor))
	stride := boundDiff.div(resolution.asVector())
	data := make([]*link, resolution.x*resolution.y*resolution.z)

	grid := &grid{
		minBound,
		maxBound,
		stride,
		data,
		resolution,
	}
	grid.populate(objs)

	return grid
}

func bounds(objs ObjectList) (Vector, Vector) {
	inf := math.Inf(+1)
	minBound, maxBound := Vect(+inf, +inf, +inf), Vect(-inf, -inf, -inf)
	for _, obj := range objs {
		min, max := obj.bound()
		minBound = minBound.Min(min)
		maxBound = maxBound.Max(max)
	}
	return minBound, maxBound
}

func (g *grid) populate(objs ObjectList) {
	for _, obj := range objs {
		min, max := obj.bound()
		minCoord := truncate(min.Sub(g.minBound).div(g.stride)).min(g.resolution.sub(triple{1, 1, 1}))
		maxCoord := truncate(max.Sub(g.minBound).div(g.stride)).min(g.resolution.sub(triple{1, 1, 1}))
		var pos triple
		for pos.x = minCoord.x; pos.x <= maxCoord.x; pos.x++ {
			for pos.y = minCoord.y; pos.y <= maxCoord.y; pos.y++ {
				for pos.z = minCoord.z; pos.z <= maxCoord.z; pos.z++ {
					idx := g.dataIndex(pos)
					g.data[idx] = &link{g.data[idx], obj}
				}
			}
		}
	}
}

func (g *grid) closestHit(r ray) (intersection, material, bool) {

	var distance float64
	if !g.insideBoundingBox(r.start) {
		var hit bool
		distance, hit = g.hitBoundingBox(r)
		if !hit {
			return intersection{}, material{}, false
		}
	}

	cellCoordsFloat := g.cellCoordsFloat(r.at(distance))
	pos := g.cellCoordsInt(cellCoordsFloat)
	delta := g.delta(r)
	inc := g.inc(r)
	next := g.next(cellCoordsFloat, r)

	for true {

		if intersection, material, hit := g.findHitInCell(pos, next, r); hit {
			return intersection, material, true
		}

		var exitGrid bool
		next, pos, exitGrid = g.nextCell(next, delta, pos, inc)
		if exitGrid {
			break
		}
	}

	return intersection{}, material{}, false
}

func (g *grid) insideBoundingBox(v Vector) bool {
	return true &&
		v.X >= g.minBound.X && v.X <= g.maxBound.X &&
		v.Y >= g.minBound.Y && v.Y <= g.maxBound.Y &&
		v.Z >= g.minBound.Z && v.Z <= g.maxBound.Z
}

func (g *grid) hitBoundingBox(r ray) (float64, bool) {

	tx1 := (g.minBound.X - r.start.X) / r.dir.X
	tx2 := (g.maxBound.X - r.start.X) / r.dir.X
	ty1 := (g.minBound.Y - r.start.Y) / r.dir.Y
	ty2 := (g.maxBound.Y - r.start.Y) / r.dir.Y
	tz1 := (g.minBound.Z - r.start.Z) / r.dir.Z
	tz2 := (g.maxBound.Z - r.start.Z) / r.dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)

	tmin = math.Max(tmin, math.Min(tx1, tx2))
	tmax = math.Min(tmax, math.Max(tx1, tx2))
	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))
	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmin, tmin <= tmax && tmin >= 0
}

func (g *grid) cellCoordsFloat(v Vector) Vector {
	return v.
		Sub(g.minBound).
		div(g.stride)
}

func (g *grid) cellCoordsInt(cellCoordsFloat Vector) triple {
	return truncate(cellCoordsFloat).
		min(g.resolution.sub(triple{1, 1, 1})).
		max(triple{})
}

func (g *grid) delta(r ray) Vector {
	return g.stride.
		div(r.dir).
		abs()
}

func (g *grid) inc(r ray) triple {
	return truncate(
		r.dir.sign(),
	)
}

func (g *grid) next(cellCoordsFloat Vector, r ray) Vector {
	return g.cellCoordsInt(cellCoordsFloat).asVector().
		Add(r.dir.
			sign().
			Scale(0.5).
			Add(Vect(0.5, 0.5, 0.5)),
		).
		mul(g.stride).
		Sub(r.start.Sub(g.minBound)).
		div(r.dir)
}

func (g *grid) nextCell(next, delta Vector, pos, inc triple) (Vector, triple, bool) {

	// TODO: Is is numerically stable to keep incrementing next? Could we
	// instead compute it fresh each time? Does it really matter if it's
	// numerically stable?
	var exitGrid bool
	switch {
	case next.X < math.Min(next.Y, next.Z):
		pos.x += inc.x
		next.X += delta.X
		exitGrid = pos.x < 0 && inc.x < 0 || pos.x >= g.resolution.x && inc.x > 0
	case next.Y < next.Z:
		pos.y += inc.y
		next.Y += delta.Y
		exitGrid = pos.y < 0 && inc.y < 0 || pos.y >= g.resolution.y && inc.y > 0
	default:
		pos.z += inc.z
		next.Z += delta.Z
		exitGrid = pos.z < 0 && inc.z < 0 || pos.z >= g.resolution.z && inc.z > 0
	}
	return next, pos, exitGrid
}

func (g *grid) dataIndex(pos triple) int {
	return pos.x + g.resolution.x*pos.y + g.resolution.x*g.resolution.y*pos.z
}

func (g *grid) findHitInCell(pos triple, next Vector, r ray) (intersection, material, bool) {

	var closest struct {
		intersection intersection
		material     material
		hit          bool
	}

	for link := g.data[g.dataIndex(pos)]; link != nil; link = link.next {
		intersection, hit := link.obj.intersect(r)
		if !hit {
			continue
		}
		nextCell := addULPs(math.Min(next.X, math.Min(next.Y, next.Z)), ulpFudgeFactor)
		if intersection.distance > nextCell {
			continue
		}
		if !closest.hit || intersection.distance < closest.intersection.distance {
			closest.intersection = intersection
			closest.material = link.obj.material
			closest.hit = true
		}
	}

	return closest.intersection, closest.material, closest.hit
}

type link struct {
	next *link
	obj  Object
}
