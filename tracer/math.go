package main

import "math"

type Vector struct {
	X, Y, Z float64
}

// Vect creates a new vector with elements x, y, and z. This is a convenience
// function to help users of the grayt package create Vectors in a shorthand
// form.
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

type Ray struct {
	Start, Dir Vector
}

func (r Ray) At(t float64) Vector {
	return r.Start.Add(r.Dir.Scale(t))
}
