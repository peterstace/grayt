package graytlib

import (
	"encoding/json"
	"math"
)

const alignedBoxT = "aligned_box"

type AlignedBox struct {
	Corner1, Corner2 Vect
}

func (b AlignedBox) MarshalJSON() ([]byte, error) {
	type alias AlignedBox
	return json.Marshal(struct {
		Type string
		alias
	}{alignedBoxT, alias(b)})
}

func (b AlignedBox) MakeSurfaces() []Surface {
	return []Surface{
		&alignedBox{
			min: b.Corner1.Min(b.Corner2),
			max: b.Corner1.Max(b.Corner2),
		},
	}
}

type alignedBox struct {
	min, max Vect
}

func (b *alignedBox) Intersect(r Ray) (Intersection, bool) {

	tx1 := (b.min.X - r.Start.X) / r.Dir.X
	tx2 := (b.max.X - r.Start.X) / r.Dir.X
	ty1 := (b.min.Y - r.Start.Y) / r.Dir.Y
	ty2 := (b.max.Y - r.Start.Y) / r.Dir.Y
	tz1 := (b.min.Z - r.Start.Z) / r.Dir.Z
	tz2 := (b.max.Z - r.Start.Z) / r.Dir.Z

	tmin, tmax := math.Inf(-1), math.Inf(+1)
	var nMin Vect
	var nMax Vect

	if math.Min(tx1, tx2) > tmin {
		if tx1 < tx2 {
			tmin = tx1
			nMin = Vect{-1, 0, 0}
		} else {
			tmin = tx2
			nMin = Vect{1, 0, 0}
		}
	}
	if math.Max(tx1, tx2) < tmax {
		if tx1 > tx2 {
			tmax = tx1
			nMax = Vect{-1, 0, 0}
		} else {
			tmax = tx2
			nMax = Vect{1, 0, 0}
		}
	}

	if math.Min(ty1, ty2) > tmin {
		if ty1 < ty2 && ty1 > 0 {
			tmin = ty1
			nMin = Vect{0, -1, 0}
		} else {
			tmin = ty2
			nMin = Vect{0, 1, 0}
		}
	}
	if math.Max(ty1, ty2) < tmax {
		if ty1 > ty2 {
			tmax = ty1
			nMax = Vect{0, -1, 0}
		} else {
			tmax = ty2
			nMax = Vect{0, 1, 0}
		}
	}

	if math.Min(tz1, tz2) > tmin {
		if tz1 < tz2 && tz1 > 0 {
			tmin = tz1
			nMin = Vect{0, 0, -1}
		} else {
			tmin = tz2
			nMin = Vect{0, 0, 1}
		}
	}
	if math.Max(tz1, tz2) < tmax {
		if tz1 > tz2 {
			tmax = tz1
			nMax = Vect{0, 0, -1}
		} else {
			tmax = tz2
			nMax = Vect{0, 0, 1}
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
