package grayt

import "math"

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
		minCoord := truncate(min.Sub(grid.origin).div(grid.stride))
		maxCoord := truncate(max.Sub(grid.origin).div(grid.stride))
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

	delta := g.stride.div(r.dir)
	inc := truncate(delta.sign())

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
	if pos.x < 0 {
		pos.x = -1
	}
	if pos.y < 0 {
		pos.y = -1
	}
	if pos.y < 0 {
		pos.y = -1
	}
	if pos.x >= g.resolution.x {
		pos.x = g.resolution.x
	}
	if pos.y >= g.resolution.y {
		pos.y = g.resolution.y
	}
	if pos.z >= g.resolution.z {
		pos.z = g.resolution.z
	}

loop:
	for true {

		// TODO: Is is numerically stable to keep incrementing next? Could we
		// instead compute it fresh each time? Does it really matter if it's
		// numerically stable?
		switch {
		case next.X < math.Min(next.Y, next.Z):
			pos.x += inc.x
			next.X += delta.X
			if pos.x < 0 || pos.x >= g.resolution.x {
				break loop
			}
		case next.Y < next.Z:
			pos.y += inc.y
			next.Y += delta.Y
			if pos.y < 0 || pos.y >= g.resolution.y {
				break loop
			}
		default:
			pos.z += inc.z
			next.Y += delta.Y
			if pos.z < 0 || pos.z >= g.resolution.z {
				break loop
			}
		}

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

func truncate(v Vector) intVect {
	return intVect{
		int(v.X),
		int(v.Y),
		int(v.Z),
	}
}

type intVect struct {
	x, y, z int
}

func (v intVect) asVector() Vector {
	return Vector{
		float64(v.x),
		float64(v.y),
		float64(v.z),
	}
}
