package main

import (
	"image/color"
	"math"
)

var (
	White = colour{1, 1, 1}
	Red   = colour{1, 0, 0}
	Green = colour{0, 1, 0}
	Blue  = colour{0, 0, 1}
	Black = colour{0, 0, 0}
)

type colour struct {
	R, G, B float64
}

func (c colour) Add(rhs colour) colour {
	return colour{
		c.R + rhs.R,
		c.G + rhs.G,
		c.B + rhs.B,
	}
}

func (c colour) Scale(f float64) colour {
	return colour{
		c.R * f,
		c.G * f,
		c.B * f,
	}
}

func (c colour) Pow(exp float64) colour {
	return colour{
		math.Pow(c.R, exp),
		math.Pow(c.G, exp),
		math.Pow(c.B, exp),
	}
}

func (c colour) Mul(r colour) colour {
	return colour{
		c.R * r.R,
		c.G * r.G,
		c.B * r.B,
	}
}

func (c colour) Div(r colour) colour {
	return colour{
		c.R / r.R,
		c.G / r.G,
		c.B / r.B,
	}
}

func (c colour) Square() colour {
	return colour{
		c.R * c.R,
		c.G * c.G,
		c.B * c.B,
	}
}

func (c colour) ToNRGBA() color.NRGBA {
	return color.NRGBA{
		R: float64ToUint8(c.R),
		G: float64ToUint8(c.G),
		B: float64ToUint8(c.B),
		A: 0xff,
	}
}

func float64ToUint8(f float64) uint8 {
	switch {
	case f >= 1.0:
		return 0xff
	case f < 0.0:
		return 0x00
	default:
		// Since f >= 0.0 and f < 1.0, this returns a value beween 0x00 and
		// 0xff (inclusive).
		return uint8(f * 0x100)
	}
}
