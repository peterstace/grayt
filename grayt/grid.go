package grayt

import (
	"log"
	"math"
)

type grid struct {
	origin     Vector
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
		stride,
		data,
		res,
	}

	for _, obj := range objs {
		min, max := obj.bound()
		minCoord := truncate(min.Sub(grid.origin).div(grid.stride)).min(grid.resolution.sub(intVect{1, 1, 1}))
		maxCoord := truncate(max.Sub(grid.origin).div(grid.stride)).min(grid.resolution.sub(intVect{1, 1, 1}))
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

func (g *grid) closestHit(r ray) (intersection, material, bool) {

	if debug {
		log.Print("==DEBUG==")
		log.Printf("Ray: %v", r)
	}

	delta := g.stride.div(r.dir).abs()
	inc := truncate(r.dir.sign())

	if debug {
		log.Printf("Delta: %v", delta)
		log.Printf("Inc: %v", inc)
	}

	ogrid := r.start.Sub(g.origin) // ray start relative to the grid origin.
	cellCoord := ogrid.div(g.stride)
	next := cellCoord.
		floor().
		Add(r.dir.
			sign().
			Scale(0.5).
			Add(Vect(0.5, 0.5, 0.5)),
		).
		mul(g.stride).
		Sub(ogrid).
		div(r.dir)

	pos := truncate(cellCoord)

	if debug {
		log.Printf("Ogrid: %v", ogrid)
		log.Printf("CellCoord: %v", cellCoord)
		log.Printf("Next: %v", next)
		log.Printf("Pos: %v", pos)
	}

loop:
	for true {

		if debug {
			log.Print("TOP")
		}

		// TODO: Is is numerically stable to keep incrementing next? Could we
		// instead compute it fresh each time? Does it really matter if it's
		// numerically stable?
		switch {
		case next.X < math.Min(next.Y, next.Z):
			pos.x += inc.x
			next.X += delta.X
			if pos.x < 0 && inc.x < 0 || pos.x > g.resolution.x && inc.x > 0 {
				if debug {
					log.Printf("Break X")
				}
				break loop
			}
		case next.Y < next.Z:
			pos.y += inc.y
			next.Y += delta.Y
			if pos.y < 0 && inc.y < 0 || pos.y > g.resolution.y && inc.y > 0 {
				if debug {
					log.Printf("Break Y")
				}
				break loop
			}
		default:
			pos.z += inc.z
			next.Z += delta.Z
			if pos.z < 0 && inc.z < 0 || pos.z > g.resolution.z && inc.z > 0 {
				if debug {
					log.Printf("Break Z")
				}
				break loop
			}
		}

		if debug {
			log.Printf("Pos: %v", pos)
			log.Printf("Next: %v", next)
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
				log.Printf("Look for hit: %v", link.obj)
			}
			intersection, hit := link.obj.intersect(r)
			if !hit {
				continue
			}
			if intersection.distance > math.Min(next.X, math.Min(next.Y, next.Z)) {
				continue
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

func (v intVect) sub(u intVect) intVect {
	return intVect{
		v.x - u.x,
		v.y - u.y,
		v.z - u.z,
	}
}
