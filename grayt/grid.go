package grayt

import (
	"math"

	"github.com/peterstace/grayt/xmath"
)

type grid struct {
	minBound xmath.Vector
	maxBound xmath.Vector

	stride     xmath.Vector
	data       []*link
	resolution xmath.Triple
}

func newGrid(lambda float64, objs ObjectList) *grid {
	minBound, maxBound := bounds(objs)
	boundDiff := maxBound.Sub(minBound)
	volume := boundDiff.X * boundDiff.Y * boundDiff.Z

	resolutionFactor := math.Pow(lambda*float64(len(objs))/volume, 1.0/3.0)
	resolution := xmath.Truncate(boundDiff.Scale(resolutionFactor)).Max(xmath.Triple{1, 1, 1})
	stride := boundDiff.Div(resolution.AsVector())
	data := make([]*link, resolution.Z*resolution.Z*resolution.Z)

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

func bounds(objs ObjectList) (xmath.Vector, xmath.Vector) {
	inf := math.Inf(+1)
	minBound, maxBound := xmath.Vect(+inf, +inf, +inf), xmath.Vect(-inf, -inf, -inf)
	for _, obj := range objs {
		min, max := obj.Surface.bound()
		minBound = minBound.Min(min)
		maxBound = maxBound.Max(max)
	}
	return minBound, maxBound
}

func (g *grid) populate(objs ObjectList) {
	for _, obj := range objs {
		min, max := obj.Surface.bound()
		minCoord := xmath.Truncate(min.Sub(g.minBound).Div(g.stride)).Min(g.resolution.Sub(xmath.Triple{1, 1, 1}))
		maxCoord := xmath.Truncate(max.Sub(g.minBound).Div(g.stride)).Min(g.resolution.Sub(xmath.Triple{1, 1, 1}))
		var pos xmath.Triple
		for pos.X = minCoord.X; pos.X <= maxCoord.X; pos.X++ {
			for pos.Y = minCoord.Y; pos.Y <= maxCoord.Y; pos.Y++ {
				for pos.Z = minCoord.Z; pos.Z <= maxCoord.Z; pos.Z++ {
					idx := g.dataIndex(pos)
					g.data[idx] = &link{g.data[idx], obj}
				}
			}
		}
	}
}

func (g *grid) closestHit(r xmath.Ray) (intersection, material, bool) {

	var distance float64
	if !g.insideBoundingBox(r.Start) {
		var hit bool
		distance, hit = g.hitBoundingBox(r)
		if !hit {
			return intersection{}, material{}, false
		}
	}

	cellCoordsFloat := g.cellCoordsFloat(r.At(distance))
	initialPos := g.cellCoordsInt(cellCoordsFloat)
	delta := g.delta(r)
	inc := g.inc(r)
	initialNextHitDistance := g.next(cellCoordsFloat, r)

	var pos = initialPos

	for true {

		nextHitDistance := pos.Sub(initialPos).AsVector().Abs().Mul(delta).Add(initialNextHitDistance)

		if intersection, material, hit := g.findHitInCell(pos, nextHitDistance, r); hit {
			return intersection, material, true
		}

		var exitGrid bool
		pos, exitGrid = g.nextCell(nextHitDistance, initialPos, pos, inc)
		if exitGrid {
			break
		}
	}

	return intersection{}, material{}, false
}

func (g *grid) insideBoundingBox(v xmath.Vector) bool {
	return true &&
		v.X >= g.minBound.X && v.X <= g.maxBound.X &&
		v.Y >= g.minBound.Y && v.Y <= g.maxBound.Y &&
		v.Z >= g.minBound.Z && v.Z <= g.maxBound.Z
}

func (g *grid) hitBoundingBox(r xmath.Ray) (float64, bool) {

	tx1 := (g.minBound.X - r.Start.X) / r.Dir.X
	tx2 := (g.maxBound.X - r.Start.X) / r.Dir.X
	ty1 := (g.minBound.Y - r.Start.Y) / r.Dir.Y
	ty2 := (g.maxBound.Y - r.Start.Y) / r.Dir.Y
	tz1 := (g.minBound.Z - r.Start.Z) / r.Dir.Z
	tz2 := (g.maxBound.Z - r.Start.Z) / r.Dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)

	tmin = math.Max(tmin, math.Min(tx1, tx2))
	tmax = math.Min(tmax, math.Max(tx1, tx2))
	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))
	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmin, tmin <= tmax && tmin >= 0
}

func (g *grid) cellCoordsFloat(v xmath.Vector) xmath.Vector {
	return v.
		Sub(g.minBound).
		Div(g.stride)
}

func (g *grid) cellCoordsInt(cellCoordsFloat xmath.Vector) xmath.Triple {
	return xmath.Truncate(cellCoordsFloat).
		Min(g.resolution.Sub(xmath.Triple{1, 1, 1})).
		Max(xmath.Triple{})
}

func (g *grid) delta(r xmath.Ray) xmath.Vector {
	return g.stride.
		Div(r.Dir).
		Abs()
}

func (g *grid) inc(r xmath.Ray) xmath.Triple {
	return xmath.Truncate(
		r.Dir.Sign(),
	)
}

func (g *grid) next(cellCoordsFloat xmath.Vector, r xmath.Ray) xmath.Vector {
	return g.cellCoordsInt(cellCoordsFloat).AsVector().
		Add(r.Dir.
			Sign().
			Scale(0.5).
			Add(xmath.Vect(0.5, 0.5, 0.5)),
		).
		Mul(g.stride).
		Sub(r.Start.Sub(g.minBound)).
		Div(r.Dir)
}

func (g *grid) nextCell(next xmath.Vector, initialPos, pos, inc xmath.Triple) (xmath.Triple, bool) {

	var exitGrid bool
	switch {
	case next.X < math.Min(next.Y, next.Z):
		pos.X += inc.X
		exitGrid = pos.X < 0 && inc.X < 0 || pos.X >= g.resolution.X && inc.X > 0
	case next.Y < next.Z:
		pos.Y += inc.Y
		exitGrid = pos.Y < 0 && inc.Y < 0 || pos.Y >= g.resolution.Y && inc.Y > 0
	default:
		pos.Z += inc.Z
		exitGrid = pos.Z < 0 && inc.Z < 0 || pos.Z >= g.resolution.Z && inc.Z > 0
	}
	return pos, exitGrid
}

func (g *grid) dataIndex(pos xmath.Triple) int {
	return pos.X + g.resolution.X*pos.Y + g.resolution.X*g.resolution.Y*pos.Z
}

func (g *grid) findHitInCell(pos xmath.Triple, next xmath.Vector, r xmath.Ray) (intersection, material, bool) {

	var closest struct {
		intersection intersection
		material     material
		hit          bool
	}

	for link := g.data[g.dataIndex(pos)]; link != nil; link = link.next {
		intersection, hit := link.obj.Surface.intersect(r)
		if !hit {
			continue
		}
		nextCell := xmath.AddULPs(math.Min(next.X, math.Min(next.Y, next.Z)), ulpFudgeFactor)
		if intersection.distance > nextCell {
			continue
		}
		if !closest.hit || intersection.distance < closest.intersection.distance {
			closest.intersection = intersection
			closest.material = link.obj.Material
			closest.hit = true
		}
	}

	return closest.intersection, closest.material, closest.hit
}

type link struct {
	next *link
	obj  Object
}
