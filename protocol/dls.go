package protocol

import (
	"math"

	"github.com/peterstace/grayt/xmath"
)

// TODO: This should go in the scenelib service. But for now, it will go here
// since it will be needed by both grayt.main and scenelib until the scenelib
// is completely out of grayt.main.

func DefaultCamera() Camera {
	return Camera{
		Location:             xmath.Vect(0, 10, 10),
		LookingAt:            xmath.Vect(0, 0, 0),
		UpDirection:          xmath.Vect(0, 1, 0),
		FieldOfViewInRadians: 90 * math.Pi / 180,
		FocalLength:          1.0,
		FocalRatio:           math.MaxFloat64,
		AspectWide:           1,
		AspectHigh:           1,
	}
}

type SurfaceList []interface{}
type ObjectList []Object

func Combine(m Material, surfaces ...SurfaceList) []Object {
	var objs []Object
	for _, s := range MergeSurfaceLists(surfaces...) {
		objs = append(objs, Object{Material: m, Surface: s})
	}
	return objs
}

func MergeObjectLists(lists ...ObjectList) ObjectList {
	var objs []Object
	for i := range lists {
		for j := range lists[i] {
			objs = append(objs, lists[i][j])
		}
	}
	return objs
}

func MergeSurfaceLists(lists ...SurfaceList) SurfaceList {
	var merged SurfaceList
	for i := range lists {
		for j := range lists[i] {
			merged = append(merged, lists[i][j])
		}
	}
	return merged
}

func Square(a, b, c, d xmath.Vector) SurfaceList {
	return SurfaceList{
		Triangle{a, b, c},
		Triangle{c, d, a},
	}
}

func AlignedSquare(a, b xmath.Vector) SurfaceList {
	same := func(a, b float64) int {
		if a == b {
			return 1
		}
		return 0
	}
	if same(a.X, b.X)+same(a.Y, b.Y)+same(a.Z, b.Z) != 1 {
		panic("a and b must have exactly 1 dimension in common")
	}

	a, b = a.Min(b), a.Max(b)

	switch {
	case a.X == b.X:
		return SurfaceList{AlignXSquare{a.X, a.Y, b.Y, a.Z, b.Z}}
	case a.Y == b.Y:
		return SurfaceList{AlignYSquare{a.X, b.X, a.Y, a.Z, b.Z}}
	case a.Z == b.Z:
		return SurfaceList{AlignZSquare{a.X, b.X, a.Y, b.Y, a.Z}}
	default:
		panic(false)

	}
}
