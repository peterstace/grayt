package tracer

import "math"

type Vect struct {
	X, Y, Z float64
}

func (v Vect) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vect) Length2() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vect) Unit() Vect {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vect{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v Vect) Extended(mul float64) Vect {
	return Vect{
		X: v.X * mul,
		Y: v.Y * mul,
		Z: v.Z * mul,
	}
}

func VectAdd(v1, v2 Vect) Vect {
	return Vect{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func VectSub(v1, v2 Vect) Vect {
	return Vect{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

func VectDot(lhs, rhs Vect) float64 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z
}

func VectCross(rhs, lhs Vect) Vect {
	return Vect{
		X: lhs.Y*rhs.Z - lhs.Z*rhs.Y,
		Y: lhs.Z*rhs.X - lhs.X*rhs.Z,
		Z: lhs.X*rhs.Y - lhs.Y*rhs.X,
	}
}
