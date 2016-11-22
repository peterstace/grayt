package grayt

import (
	"image/color"
	"math"
)

type Colour struct {
	R, G, B float64
}

func newColourFromRGB(rgb uint32) Colour {
	r := (rgb & 0xff0000) >> 16
	g := (rgb & 0x00ff00) >> 8
	b := (rgb & 0x0000ff)
	return Colour{
		R: float64(r) / 0xff,
		G: float64(g) / 0xff,
		B: float64(b) / 0xff,
	}
}

func (c Colour) add(rhs Colour) Colour {
	return Colour{
		c.R + rhs.R,
		c.G + rhs.G,
		c.B + rhs.B,
	}
}

func (c Colour) scale(f float64) Colour {
	return Colour{
		c.R * f,
		c.G * f,
		c.B * f,
	}
}

func (c Colour) pow(exp float64) Colour {
	return Colour{
		math.Pow(c.R, exp),
		math.Pow(c.G, exp),
		math.Pow(c.B, exp),
	}
}

func (c Colour) mul(r Colour) Colour {
	return Colour{
		c.R * r.R,
		c.G * r.G,
		c.B * r.B,
	}
}

func (c Colour) div(r Colour) Colour {
	return Colour{
		c.R / r.R,
		c.G / r.G,
		c.B / r.B,
	}
}

func (c Colour) toNRGBA() color.NRGBA {
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
