package grayt

import (
	"encoding/json"
	"fmt"
	"math"
)

type Vect struct {
	X, Y, Z float64
}

func (v Vect) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vect) Length2() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vect) Unit() Vect {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vect{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v Vect) Extended(mul float64) Vect {
	return Vect{
		X: v.X * mul,
		Y: v.Y * mul,
		Z: v.Z * mul,
	}
}

func (v Vect) Add(u Vect) Vect {
	return Vect{
		X: v.X + u.X,
		Y: v.Y + u.Y,
		Z: v.Z + u.Z,
	}
}

func (v Vect) Sub(u Vect) Vect {
	return Vect{
		X: v.X - u.X,
		Y: v.Y - u.Y,
		Z: v.Z - u.Z,
	}
}

func (v Vect) Dot(u Vect) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v Vect) Cross(u Vect) Vect {
	return Vect{
		X: v.Y*u.Z - v.Z*u.Y,
		Y: v.Z*u.X - v.X*u.Z,
		Z: v.X*u.Y - v.Y*u.X,
	}
}

func (v *Vect) UnmarshalJSON(p []byte) error {
	var record []float64
	if err := json.Unmarshal(p, &record); err != nil {
		return err
	}
	if len(record) != 3 {
		return fmt.Errorf("cannot marshal array of size %d into Go value of type grayt.Vect",
			len(record))
	}
	v.X = record[0]
	v.Y = record[1]
	v.Z = record[2]
	return nil
}

func (v Vect) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%g,%g,%g]", v.X, v.Y, v.Z)), nil
}

type Ray struct {
	Start, Dir Vect
}

func (r Ray) At(t float64) Vect {
	return r.Start.Add(r.Dir.Extended(t))
}
