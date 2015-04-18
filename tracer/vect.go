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

func (v Vect) Add(u Vect) Vect {
	return Vect{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

func (v Vect) Sub(u Vect) Vect {
	return Vect{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

func (v Vect) Dot(u Vect) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func VectCross(rhs, lhs Vect) Vect {
	return Vect{
		X: lhs.Y*rhs.Z - lhs.Z*rhs.Y,
		Y: lhs.Z*rhs.X - lhs.X*rhs.Z,
		Z: lhs.X*rhs.Y - lhs.Y*rhs.X,
	}
}
