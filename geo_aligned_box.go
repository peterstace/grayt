package grayt

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

	//log.Print(b, r)

	var norm Vect
	tmin, tmax := math.Inf(+1), math.Inf(-1)

	tx1 := (b.min.X - r.Start.X) / r.Dir.X
	tx2 := (b.max.X - r.Start.X) / r.Dir.X
	if tx1 < tmin {
		tmin = tx1
		norm = Vect{-1, 0, 0}
	}
	if tx2 < tmin {
		tmin = tx2
		norm = Vect{+1, 0, 0}
	}
	tmax = math.Max(tmax, math.Max(tx1, tx2))

	ty1 := (b.min.Y - r.Start.Y) / r.Dir.Y
	ty2 := (b.max.Y - r.Start.Y) / r.Dir.Y
	if ty1 < tmin {
		tmin = ty1
		norm = Vect{0, -1, 0}
	}
	if ty2 < tmin {
		tmin = ty2
		norm = Vect{0, +1, 0}
	}
	tmax = math.Max(tmax, math.Max(ty1, ty2))

	tz1 := (b.min.Z - r.Start.Z) / r.Dir.Z
	tz2 := (b.max.Z - r.Start.Z) / r.Dir.Z
	if tz1 < tmin {
		tmin = tz1
		norm = Vect{0, 0, -1}
	}
	if tz2 < tmin {
		tmin = tz2
		norm = Vect{0, 0, +1}
	}
	tmax = math.Max(tmax, math.Max(tz1, tz2))

	//log.Print(tmin, tmax)

	var t float64
	if tmin > tmax {
		t = -1.0
	} else if tmin < 0 {
		t = tmax
	} else {
		t = tmin
	}

	return Intersection{
		UnitNormal: norm,
		Distance:   t,
	}, t > 0
}
