package grayt

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"math"
)

type Camera struct {
	Location             Vector
	ViewDirection        Vector
	UpDirection          Vector
	FieldOfViewInDegrees float64
	FocalLength          float64 // Distance to the focus plane.
	FocalRatio           float64 // Ratio between the focal length and the diameter of the eye.
}

func (c Camera) String() string {
	return fmt.Sprintf("Location=%v ViewDir=%v UpDir=%v FOV=%v FocalLength=%v FocalRation=%v",
		c.Location, c.ViewDirection, c.UpDirection, c.FieldOfViewInDegrees, c.FocalLength, c.FocalRatio)
}

func DefaultCamera() Camera {
	return Camera{
		Location:             Vector{},
		ViewDirection:        Vect(0, 0, -1),
		UpDirection:          Vect(0, 1, 0),
		FieldOfViewInDegrees: 100,
		FocalLength:          10,
		FocalRatio:           math.Inf(+1),
	}
}

type Triangle struct {
	A, B, C   Vector
	Colour    Colour
	Emittance float64
}

func (t Triangle) String() string {
	return fmt.Sprintf("A=%v B=%v C=%v Colour=%v Emittance=%v",
		t.A, t.B, t.C, t.Colour, t.Emittance)
}

func JoinTriangles(ts ...[]Triangle) []Triangle {
	var all []Triangle
	for _, t := range ts {
		all = append(all, t...)
	}
	return all
}

type Scene struct {
	Camera    Camera
	Triangles []Triangle
}

func (s Scene) hash() string {

	// Calculate hash.
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	h.Write([]byte(s.Camera.String()))
	for _, t := range s.Triangles {
		h.Write([]byte(t.String()))
	}

	// Calculate base64 encoded hash.
	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.RawURLEncoding, &buf)
	binary.Write(
		enc,
		binary.LittleEndian,
		h.Sum64(),
	)
	enc.Close()
	return buf.String()
}
