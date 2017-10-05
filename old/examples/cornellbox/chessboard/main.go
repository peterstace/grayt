package main

import (
	"flag"
	"fmt"
	"strconv"

	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

var (
	focalDepth = flag.Float64("focal-depth", 0.5, "focal depth")
	focalRatio = flag.Float64("focal-ratio", 20, "focal ratio")
)

func init() {
}

func main() {
	flag.Parse()
	c := Cam(1.3)
	c = c.With(FocalLengthAndRatio(1.3+*focalDepth, *focalRatio))
	Run(
		fmt.Sprintf(
			"cornellbox_chessboard_%v_%v",
			strconv.FormatFloat(*focalDepth, 'f', -1, 64),
			strconv.FormatFloat(*focalRatio, 'f', -1, 64),
		),
		Scene{
			Camera: c,
			Objects: Group(
				ShortBlock(),
				TallBlock(),
				chessboard(Vect(0, 0, 0), Vect(0, 1, -1), Red, Yellow),
				chessboard(Vect(1, 0, 0), Vect(1, 1, -1), Cyan, Blue),
				chessboard(Vect(0, 0, 0), Vect(1, 0, -1), Magenta, Green),
				chessboard(Vect(0, 0, -1), Vect(1, 1, -1), Magenta, Green),
				Ceiling,
				LeftWall.With(ColourRGB(Red)),
				RightWall.With(ColourRGB(Green)),
				CeilingLight().With(Emittance(5.0)),
			),
		})
}

func chessboard(v1, v2 Vector, c1, c2 uint32) ObjectList {

	const divisions = 8

	var a, b, c func(*Vector) *float64

	switch {
	case v1.X == v2.X:
		a = func(v *Vector) *float64 { return &v.X }
		b = func(v *Vector) *float64 { return &v.Y }
		c = func(v *Vector) *float64 { return &v.Z }
	case v1.Y == v2.Y:
		b = func(v *Vector) *float64 { return &v.X }
		a = func(v *Vector) *float64 { return &v.Y }
		c = func(v *Vector) *float64 { return &v.Z }
	case v1.Z == v2.Z:
		b = func(v *Vector) *float64 { return &v.X }
		c = func(v *Vector) *float64 { return &v.Y }
		a = func(v *Vector) *float64 { return &v.Z }
	default:
		panic(false)
	}

	var objList ObjectList
	for i := 0; i < divisions; i++ {
		for j := 0; j < divisions; j++ {

			uMin := float64(i) / divisions
			uMax := float64(i+1) / divisions
			vMin := float64(j) / divisions
			vMax := float64(j+1) / divisions

			var corner1 Vector
			*a(&corner1) = *a(&v1)
			*b(&corner1) = *b(&v1)*(1-uMin) + *b(&v2)*uMin
			*c(&corner1) = *c(&v1)*(1-vMin) + *c(&v2)*vMin

			var corner2 Vector
			*a(&corner2) = *a(&v1)
			*b(&corner2) = *b(&v1)*(1-uMax) + *b(&v2)*uMax
			*c(&corner2) = *c(&v1)*(1-vMax) + *c(&v2)*vMax

			objList = Group(objList,
				AlignedSquare(corner1, corner2).With(
					ColourRGB([]uint32{c1, c2}[(i+j)%2]),
				),
			)
		}
	}
	return objList
}
