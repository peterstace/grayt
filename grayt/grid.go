package grayt

import (
	"log"
	"math"
)

type grid struct {
	minBound Vector
	maxBound Vector

	stride     Vector
	data       []*link
	resolution intVect
}

func newGrid(lambda float64, objs ObjectList) *grid {

	inf := math.Inf(+1)
	minBound, maxBound := Vect(+inf, +inf, +inf), Vect(-inf, -inf, -inf)
	for _, obj := range objs {
		min, max := obj.bound()
		minBound = minBound.Min(min)
		maxBound = maxBound.Max(max)
	}
	boundDiff := maxBound.Sub(minBound)
	volume := boundDiff.X * boundDiff.Y * boundDiff.Z

	resolutionFactor := math.Pow(lambda*float64(len(objs))/volume, 1.0/3.0)
	res := truncate(boundDiff.Scale(resolutionFactor))
	stride := boundDiff.div(res.asVector())

	data := make([]*link, res.x*res.y*res.z)

	grid := &grid{
		minBound,
		maxBound,
		stride,
		data,
		res,
	}

	for _, obj := range objs {
		min, max := obj.bound()
		minCoord := truncate(min.Sub(grid.minBound).div(grid.stride)).min(grid.resolution.sub(intVect{1, 1, 1}))
		maxCoord := truncate(max.Sub(grid.minBound).div(grid.stride)).min(grid.resolution.sub(intVect{1, 1, 1}))
		log.Printf("%#v", obj)
		log.Printf("Min: %v", min)
		log.Printf("Max: %v", max)
		log.Printf("MinC: %v", minCoord)
		log.Printf("MaxC: %v", maxCoord)
		var pos intVect
		for pos.x = minCoord.x; pos.x <= maxCoord.x; pos.x++ {
			for pos.y = minCoord.y; pos.y <= maxCoord.y; pos.y++ {
				for pos.z = minCoord.z; pos.z <= maxCoord.z; pos.z++ {
					idx := grid.dataIndex(pos)
					grid.data[idx] = &link{grid.data[idx], obj}
				}
			}
		}
	}

	return grid
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

func (g *grid) cellCoord(v Vector) Vector {
	return v.
		Sub(g.minBound).
		div(g.stride)

}

func (g *grid) closestHit(r ray) (intersection, material, bool) {

	if debug {
		log.Print("==DEBUG==")
		log.Printf("Ray: %v", r)
	}

	var cellCoord Vector
	if g.insideBoundingBox(r.start) {
		cellCoord = g.cellCoord(r.start)
	} else {
		if distance, hit := g.hitBoundingBox(r); !hit {
			return intersection{}, material{}, false
		} else {
			cellCoord = g.cellCoord(r.at(distance))
		}
	}

	pos := truncate(cellCoord).
		min(g.resolution).
		max(intVect{})

	delta := g.stride.div(r.dir).abs()
	inc := truncate(r.dir.sign())

	next := cellCoord.
		floor().
		Add(r.dir.
			sign().
			Scale(0.5).
			Add(Vect(0.5, 0.5, 0.5)),
		).
		mul(g.stride).
		Sub(r.start.Sub(g.minBound)).
		div(r.dir)

	if debug {
		log.Printf("CellCoord: %v", cellCoord)
		log.Printf("Pos: %v", pos)
		log.Printf("Delta: %v", delta)
		log.Printf("Inc: %v", inc)
		log.Printf("Next: %v", next)
	}

loop:
	for true {

		if debug {
			log.Print("TOP")
		}

		var closest struct {
			intersection intersection
			material     material
			hit          bool
		}
		var head *link
		if pos.x >= 0 && pos.x < g.resolution.x && pos.y >= 0 && pos.y < g.resolution.y && pos.z >= 0 && pos.z < g.resolution.z {
			head = g.data[g.dataIndex(pos)]
		}
		for link := head; link != nil; link = link.next {
			if debug {
				log.Printf("\tLook for hit: %v", link.obj)
			}
			intersection, hit := link.obj.intersect(r)
			if !hit {
				if debug {
					log.Printf("\tNo hit")
				}
				continue
			}
			nextCell := addULPs(math.Min(next.X, math.Min(next.Y, next.Z)), 50)
			if intersection.distance > nextCell {
				if debug {
					log.Printf("\tHit, but outside of cell, intersection distance: %v", intersection.distance)
				}
				continue
			}
			if debug {
				log.Printf("\tWas hit, at distance %v", intersection.distance)
			}
			if !closest.hit || intersection.distance < closest.intersection.distance {
				closest.intersection = intersection
				closest.material = link.obj.material
				closest.hit = true
			}
		}
		if closest.hit {
			return closest.intersection, closest.material, true
		}

		// TODO: Is is numerically stable to keep incrementing next? Could we
		// instead compute it fresh each time? Does it really matter if it's
		// numerically stable?
		switch {
		case next.X < math.Min(next.Y, next.Z):
			pos.x += inc.x
			next.X += delta.X
			if pos.x < 0 && inc.x < 0 || pos.x >= g.resolution.x && inc.x > 0 {
				if debug {
					log.Printf("\tPos %v", pos)
					log.Printf("\tBreak X")
				}
				break loop
			}
		case next.Y < next.Z:
			pos.y += inc.y
			next.Y += delta.Y
			if pos.y < 0 && inc.y < 0 || pos.y >= g.resolution.y && inc.y > 0 {
				if debug {
					log.Printf("\tPos %v", pos)
					log.Printf("\tBreak Y")
				}
				break loop
			}
		default:
			pos.z += inc.z
			next.Z += delta.Z
			if pos.z < 0 && inc.z < 0 || pos.z >= g.resolution.z && inc.z > 0 {
				if debug {
					log.Printf("\tPos %v", pos)
					log.Printf("\tBreak Z")
				}
				break loop
			}
		}

		if debug {
			log.Printf("\tPos: %v", pos)
			log.Printf("\tNext: %v", next)
		}
	}

	return intersection{}, material{}, false
}

func (g *grid) dataIndex(pos intVect) int {
	return pos.x + g.resolution.x*pos.y + g.resolution.x*g.resolution.y*pos.z
}

func sign(f float64) int {
	if f < 0 {
		return -1
	}
	return +1
}

type link struct {
	next *link
	obj  Object
}

type intVect struct {
	x, y, z int
}

func truncate(v Vector) intVect {
	return intVect{
		int(v.X),
		int(v.Y),
		int(v.Z),
	}
}

func (v intVect) asVector() Vector {
	return Vector{
		float64(v.x),
		float64(v.y),
		float64(v.z),
	}
}

func (v intVect) min(u intVect) intVect {
	if v.x > u.x {
		v.x = u.x
	}
	if v.y > u.y {
		v.y = u.y
	}
	if v.z > u.z {
		v.z = u.z
	}
	return v
}

func (v intVect) max(u intVect) intVect {
	if v.x < u.x {
		v.x = u.x
	}
	if v.y < u.y {
		v.y = u.y
	}
	if v.z < u.z {
		v.z = u.z
	}
	return v
}

func (v intVect) sub(u intVect) intVect {
	return intVect{
		v.x - u.x,
		v.y - u.y,
		v.z - u.z,
	}
}
