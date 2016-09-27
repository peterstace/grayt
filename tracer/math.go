package main

import (
	"math"

	"github.com/peterstace/grayt/scene"
)

type vector struct {
	x, y, z float64
}

func convertVector(v scene.Vector) vector {
	return vector{v.X, v.Y, v.Z}
}

func (v vector) length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v vector) lengthSq() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v vector) unit() vector {
	length := math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
	return vector{
		x: v.x / length,
		y: v.y / length,
		z: v.z / length,
	}
}

func (v vector) scale(mul float64) vector {
	return vector{
		x: v.x * mul,
		y: v.y * mul,
		z: v.z * mul,
	}
}

func (v vector) add(u vector) vector {
	return vector{
		x: v.x + u.x,
		y: v.y + u.y,
		z: v.z + u.z,
	}
}

func (v vector) sub(u vector) vector {
	return vector{
		x: v.x - u.x,
		y: v.y - u.y,
		z: v.z - u.z,
	}
}

func (v vector) mul(u vector) vector {
	return vector{
		x: v.x * u.x,
		y: v.y * u.y,
		z: v.z * u.z,
	}
}

func (v vector) div(u vector) vector {
	return vector{
		x: v.x / u.x,
		y: v.y / u.y,
		z: v.z / u.z,
	}
}

func (v vector) dot(u vector) float64 {
	return v.x*u.x + v.y*u.y + v.z*u.z
}

func (v vector) cross(u vector) vector {
	return vector{
		x: v.y*u.z - v.z*u.y,
		y: v.z*u.x - v.x*u.z,
		z: v.x*u.y - v.y*u.x,
	}
}

func (v vector) min(u vector) vector {
	return vector{
		math.Min(v.x, u.x),
		math.Min(v.y, u.y),
		math.Min(v.z, u.z),
	}
}

func (v vector) max(u vector) vector {
	return vector{
		math.Max(v.x, u.x),
		math.Max(v.y, u.y),
		math.Max(v.z, u.z),
	}
}

type ray struct {
	start, dir vector
}

func (r ray) at(t float64) vector {
	return r.start.add(r.dir.scale(t))
}
