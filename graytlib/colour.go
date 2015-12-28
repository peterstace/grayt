package graytlib

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
)

type Colour struct {
	R, G, B float64
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

func (c Colour) Square() Colour {
	return Colour{
		c.R * c.R,
		c.G * c.G,
		c.B * c.B,
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

func (c *Colour) UnmarshalJSON(p []byte) error {
	var record []float64
	if err := json.Unmarshal(p, &record); err != nil {
		return err
	}
	if len(record) != 3 {
		return fmt.Errorf("cannot marshal array of size %d into Go value of type grayt.Colour",
			len(record))
	}
	c.R = record[0]
	c.G = record[1]
	c.B = record[2]
	return nil
}

func (c Colour) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%g,%g,%g]", c.R, c.G, c.B)), nil
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
