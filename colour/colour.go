package colour

import (
	"image/color"
	"math"
)

type Colour struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

func NewColourFromRGB(rgb uint32) Colour {
	r := (rgb & 0xff0000) >> 16
	g := (rgb & 0x00ff00) >> 8
	b := (rgb & 0x0000ff)
	return Colour{
		R: float64(r) / 0xff,
		G: float64(g) / 0xff,
		B: float64(b) / 0xff,
	}
}

func NewColourFromHSL(hue, saturation, lightness float64) Colour {

	if hue < 0 || hue > 360 {
		panic("hue must be from 0 to 360")
	}
	if saturation < 0 || saturation > 1 {
		panic("saturation must be between 0 and 1")
	}
	if lightness < 0 || lightness > 1 {
		panic("lightness must be between 0 and 1")
	}

	c := (1 - math.Abs(2*lightness-1)) * saturation // chroma
	hueAdj := hue / 60
	for hueAdj > 2 {
		hueAdj -= 2
	}
	x := c * (1 - math.Abs(hueAdj-1))

	var r, g, b float64
	switch {
	case hueAdj <= 1:
		r, g, b = c, x, 0
	case hueAdj <= 2:
		r, g, b = x, c, 0
	case hueAdj <= 3:
		r, g, b = 0, c, x
	case hueAdj <= 4:
		r, g, b = 0, x, c
	case hueAdj <= 5:
		r, g, b = x, 0, c
	case hueAdj <= 6:
		r, g, b = c, 0, x
	default:
		panic(false)
	}

	m := lightness - 0.5*c
	r += m
	g += m
	b += m

	if r < 0 || r > 1.0 {
		panic(false)
	}
	if g < 0 || g > 1.0 {
		panic(false)
	}
	if b < 0 || b > 1.0 {
		panic(false)
	}
	return Colour{r, g, b}
}

func (c Colour) Add(rhs Colour) Colour {
	return Colour{
		c.R + rhs.R,
		c.G + rhs.G,
		c.B + rhs.B,
	}
}

func (c Colour) Scale(f float64) Colour {
	return Colour{
		c.R * f,
		c.G * f,
		c.B * f,
	}
}

func (c Colour) Pow(exp float64) Colour {
	return Colour{
		math.Pow(c.R, exp),
		math.Pow(c.G, exp),
		math.Pow(c.B, exp),
	}
}

func (c Colour) Mul(r Colour) Colour {
	return Colour{
		c.R * r.R,
		c.G * r.G,
		c.B * r.B,
	}
}

func (c Colour) Div(r Colour) Colour {
	return Colour{
		c.R / r.R,
		c.G / r.G,
		c.B / r.B,
	}
}

func (c Colour) ToNRGBA() color.NRGBA {
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
