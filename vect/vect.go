package vect

import "math"

type V struct {
	X, Y, Z float64
}

func New(x, y, z float64) V {
	return V{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v V) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v V) Length2() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v V) Unit() V {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return V{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v V) Extended(mul float64) V {
	return V{
		X: v.X * mul,
		Y: v.Y * mul,
		Z: v.Z * mul,
	}
}

func Add(v1, v2 V) V {
	return V{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func Sub(v1, v2 V) V {
	return V{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

func Dot(lhs, rhs V) float64 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z
}

func Cross(rhs, lhs V) V {
	return V{
		X: lhs.Y*rhs.Z - lhs.Z*rhs.Y,
		Y: lhs.Z*rhs.X - lhs.X*rhs.Z,
		Z: lhs.X*rhs.Y - lhs.Y*rhs.X,
	}
}
