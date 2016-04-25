package grayt

import "math"

func AlignedBox(corner1, corner2 Vector) Surface {
	return &alignedBox{
		min: corner1.Min(corner2),
		max: corner1.Max(corner2),
	}
}

type alignedBox struct {
	min, max Vector
}

func (b *alignedBox) Intersect(r Ray) (Intersection, bool) {

	tx1 := (b.min.X - r.Start.X) / r.Dir.X
	tx2 := (b.max.X - r.Start.X) / r.Dir.X
	ty1 := (b.min.Y - r.Start.Y) / r.Dir.Y
	ty2 := (b.max.Y - r.Start.Y) / r.Dir.Y
	tz1 := (b.min.Z - r.Start.Z) / r.Dir.Z
	tz2 := (b.max.Z - r.Start.Z) / r.Dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)
	var nMin Vector
	var nMax Vector

	if math.Min(tx1, tx2) > tmin {
		if tx1 < tx2 {
			tmin = tx1
			nMin = Vect(-1, 0, 0)
		} else {
			tmin = tx2
			nMin = Vect(1, 0, 0)
		}
	}
	if math.Max(tx1, tx2) < tmax {
		if tx1 > tx2 {
			tmax = tx1
			nMax = Vect(-1, 0, 0)
		} else {
			tmax = tx2
			nMax = Vect(1, 0, 0)
		}
	}

	if math.Min(ty1, ty2) > tmin {
		if ty1 < ty2 && ty1 > 0 {
			tmin = ty1
			nMin = Vect(0, -1, 0)
		} else {
			tmin = ty2
			nMin = Vect(0, 1, 0)
		}
	}
	if math.Max(ty1, ty2) < tmax {
		if ty1 > ty2 {
			tmax = ty1
			nMax = Vect(0, -1, 0)
		} else {
			tmax = ty2
			nMax = Vect(0, 1, 0)
		}
	}

	if math.Min(tz1, tz2) > tmin {
		if tz1 < tz2 && tz1 > 0 {
			tmin = tz1
			nMin = Vect(0, 0, -1)
		} else {
			tmin = tz2
			nMin = Vect(0, 0, 1)
		}
	}
	if math.Max(tz1, tz2) < tmax {
		if tz1 > tz2 {
			tmax = tz1
			nMax = Vect(0, 0, -1)
		} else {
			tmax = tz2
			nMax = Vect(0, 0, 1)
		}
	}

	if tmin > tmax || tmax <= 0 {
		return Intersection{}, false
	}

	if tmin > 0 {
		return Intersection{Distance: tmin, UnitNormal: nMin}, true
	} else {
		return Intersection{Distance: tmax, UnitNormal: nMax}, true
	}
}

func (b *alignedBox) Translate(x, y, z float64) Surface {
	b.min = b.min.Add(Vect(x, y, z))
	b.max = b.max.Add(Vect(x, y, z))
	return b
}
