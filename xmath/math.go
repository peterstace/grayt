package xmath

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
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

func (v Vector) Mul(u Vector) Vector {
	return Vector{
		X: v.X * u.X,
		Y: v.Y * u.Y,
		Z: v.Z * u.Z,
	}
}

func (v Vector) Div(u Vector) Vector {
	return Vector{
		X: v.X / u.X,
		Y: v.Y / u.Y,
		Z: v.Z / u.Z,
	}
}

func (v Vector) Dot(u Vector) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v Vector) Cross(u Vector) Vector {
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

func (v Vector) Sign() Vector {
	return Vector{
		math.Copysign(1, v.X),
		math.Copysign(1, v.Y),
		math.Copysign(1, v.Z),
	}
}

func (v Vector) Floor() Vector {
	return Vector{
		math.Floor(v.X),
		math.Floor(v.Y),
		math.Floor(v.Z),
	}
}

func (v Vector) Abs() Vector {
	return Vector{
		math.Abs(v.X),
		math.Abs(v.Y),
		math.Abs(v.Z),
	}
}

func (v Vector) AddULPs(ulps int64) Vector {
	return Vector{
		AddULPs(v.X, ulps),
		AddULPs(v.Y, ulps),
		AddULPs(v.Z, ulps),
	}
}

func (v Vector) Rotate(u Vector, rads float64) Vector {
	cos := math.Cos(rads)
	sin := math.Sin(rads)
	i := Vect(
		cos+u.X*u.X*(1-cos),
		u.X*u.Y*(1-cos)-u.Z*sin,
		u.X*u.Z*(1-cos)+u.Y*sin,
	)
	j := Vect(
		u.Y*u.X*(1-cos)+u.Z*sin,
		cos+u.Y*u.Y*(1-cos),
		u.Y*u.Z*(1-cos)-u.X*sin,
	)
	k := Vect(
		u.Z*u.X*(1-cos)-u.Y*sin,
		u.Z*u.Y*(1-cos)+u.X*sin,
		cos+u.Z*u.Z*(1-cos),
	)
	return Vect(i.Dot(v), j.Dot(v), k.Dot(v))
}

func (v Vector) Proj(unit Vector) Vector {
	return unit.Scale(v.Dot(unit))
}

func (v Vector) Rej(unit Vector) Vector {
	return v.Sub(v.Proj(unit))
}

func (v Vector) X0() Vector { return Vect(0, v.Y, v.Z) }
func (v Vector) Y0() Vector { return Vect(v.X, 0, v.Z) }
func (v Vector) Z0() Vector { return Vect(v.X, v.Y, 0) }

type Ray struct {
	Start, Dir Vector
}

func (r Ray) At(t float64) Vector {
	return r.Start.Add(r.Dir.Scale(t))
}

type Triple struct {
	X, Y, Z int
}

func Truncate(v Vector) Triple {
	return Triple{
		int(v.X),
		int(v.Y),
		int(v.Z),
	}
}

func (v Triple) AsVector() Vector {
	return Vector{
		float64(v.X),
		float64(v.Y),
		float64(v.Z),
	}
}

func (v Triple) Min(u Triple) Triple {
	return Triple{
		IntMin(v.X, u.X),
		IntMin(v.Y, u.Y),
		IntMin(v.Z, u.Z),
	}
}

func (v Triple) Max(u Triple) Triple {
	return Triple{
		IntMax(v.X, u.X),
		IntMax(v.Y, u.Y),
		IntMax(v.Z, u.Z),
	}
}

func (v Triple) Sub(u Triple) Triple {
	return Triple{
		v.X - u.X,
		v.Y - u.Y,
		v.Z - u.Z,
	}
}

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}
