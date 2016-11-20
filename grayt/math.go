package grayt

import "math"

type Vector struct {
	X, Y, Z float64
}

func Vect(x, y, z float64) Vector {
	return Vector{x, y, z}
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector) Unit() Vector {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v Vector) Scale(mul float64) Vector {
	return Vector{
		X: v.X * mul,
		Y: v.Y * mul,
		Z: v.Z * mul,
	}
}

func (v Vector) Add(u Vector) Vector {
	return Vector{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

func (v Vector) Sub(u Vector) Vector {
	return Vector{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

func (v Vector) mul(u Vector) Vector {
	return Vector{
		X: v.X * u.X,
		Y: v.Y * u.Y,
		Z: v.Z * u.Z,
	}
}

func (v Vector) div(u Vector) Vector {
	return Vector{
		X: v.X / u.X,
		Y: v.Y / u.Y,
		Z: v.Z / u.Z,
	}
}

func (v Vector) dot(u Vector) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v Vector) cross(u Vector) Vector {
	return Vector{
		X: v.Y*u.Z - v.Z*u.Y,
		Y: v.Z*u.X - v.X*u.Z,
		Z: v.X*u.Y - v.Y*u.X,
	}
}

func (v Vector) min(u Vector) Vector {
	return Vector{
		math.Min(v.X, u.X),
		math.Min(v.Y, u.Y),
		math.Min(v.Z, u.Z),
	}
}

func (v Vector) max(u Vector) Vector {
	return Vector{
		math.Max(v.X, u.X),
		math.Max(v.Y, u.Y),
		math.Max(v.Z, u.Z),
	}
}

type ray struct {
	start, dir Vector
}

func (r ray) at(t float64) Vector {
	return r.start.Add(r.dir.Scale(t))
}
