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

func (v Vector) Min(u Vector) Vector {
	return Vector{
		math.Min(v.X, u.X),
		math.Min(v.Y, u.Y),
		math.Min(v.Z, u.Z),
	}
}

func (v Vector) Max(u Vector) Vector {
	return Vector{
		math.Max(v.X, u.X),
		math.Max(v.Y, u.Y),
		math.Max(v.Z, u.Z),
	}
}

func (v Vector) sign() Vector {
	return Vector{
		math.Copysign(1, v.X),
		math.Copysign(1, v.Y),
		math.Copysign(1, v.Z),
	}
}

func (v Vector) floor() Vector {
	return Vector{
		math.Floor(v.X),
		math.Floor(v.Y),
		math.Floor(v.Z),
	}
}

func (v Vector) abs() Vector {
	return Vector{
		math.Abs(v.X),
		math.Abs(v.Y),
		math.Abs(v.Z),
	}
}

func (v Vector) addULPs(ulps int64) Vector {
	return Vector{
		addULPs(v.X, ulps),
		addULPs(v.Y, ulps),
		addULPs(v.Z, ulps),
	}
}

type ray struct {
	start, dir Vector
}

func (r ray) at(t float64) Vector {
	return r.start.Add(r.dir.Scale(t))
}

type triple struct {
	x, y, z int
}

func truncate(v Vector) triple {
	return triple{
		int(v.X),
		int(v.Y),
		int(v.Z),
	}
}

func (v triple) asVector() Vector {
	return Vector{
		float64(v.x),
		float64(v.y),
		float64(v.z),
	}
}

func (v triple) min(u triple) triple {
	return triple{
		intMin(v.x, u.x),
		intMin(v.y, u.y),
		intMin(v.z, u.z),
	}
}

func (v triple) max(u triple) triple {
	return triple{
		intMax(v.x, u.x),
		intMax(v.y, u.y),
		intMax(v.z, u.z),
	}
}

func (v triple) sub(u triple) triple {
	return triple{
		v.x - u.x,
		v.y - u.y,
		v.z - u.z,
	}
}

func intMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func intMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}
