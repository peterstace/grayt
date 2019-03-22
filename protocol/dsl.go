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

type ObjectList []Object

func MergeSurfaces(surfs ...Surface) Surface {
	var all Surface
	for _, s := range surfs {
		all.Triangles = append(all.Triangles, s.Triangles...)
		all.AlignedBoxes = append(all.AlignedBoxes, s.AlignedBoxes...)
		all.Spheres = append(all.Spheres, s.Spheres...)
		all.AlignXSquares = append(all.AlignXSquares, s.AlignXSquares...)
		all.AlignYSquares = append(all.AlignYSquares, s.AlignYSquares...)
		all.AlignZSquares = append(all.AlignZSquares, s.AlignZSquares...)
		all.Discs = append(all.Discs, s.Discs...)
		all.Pipes = append(all.Pipes, s.Pipes...)
	}
	return all
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

func Square(a, b, c, d xmath.Vector) Surface {
	return Surface{
		Triangles: []Triangle{
			{a, b, c},
			{c, d, a},
		},
	}
}

func AlignedSquare(a, b xmath.Vector) Surface {
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
		return Surface{AlignXSquares: []AlignXSquare{{a.X, a.Y, b.Y, a.Z, b.Z}}}
	case a.Y == b.Y:
		return Surface{AlignYSquares: []AlignYSquare{{a.X, b.X, a.Y, a.Z, b.Z}}}
	case a.Z == b.Z:
		return Surface{AlignZSquares: []AlignZSquare{{a.X, b.X, a.Y, b.Y, a.Z}}}
	default:
		panic(false)

	}
}
