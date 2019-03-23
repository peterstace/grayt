package dsl

import (
	"math"

	"github.com/peterstace/grayt/protocol"
	"github.com/peterstace/grayt/xmath"
)

func DefaultCamera() protocol.Camera {
	return protocol.Camera{
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

type ObjectList []protocol.Object

func MergeSurfaces(surfs ...protocol.Surface) protocol.Surface {
	var all protocol.Surface
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
	var objs []protocol.Object
	for i := range lists {
		for j := range lists[i] {
			objs = append(objs, lists[i][j])
		}
	}
	return objs
}

func Square(a, b, c, d xmath.Vector) protocol.Surface {
	return protocol.Surface{
		Triangles: []protocol.Triangle{
			{a, b, c},
			{c, d, a},
		},
	}
}

func AlignedSquare(a, b xmath.Vector) protocol.Surface {
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
		return protocol.Surface{AlignXSquares: []protocol.AlignXSquare{{a.X, a.Y, b.Y, a.Z, b.Z}}}
	case a.Y == b.Y:
		return protocol.Surface{AlignYSquares: []protocol.AlignYSquare{{a.X, b.X, a.Y, a.Z, b.Z}}}
	case a.Z == b.Z:
		return protocol.Surface{AlignZSquares: []protocol.AlignZSquare{{a.X, b.X, a.Y, b.Y, a.Z}}}
	default:
		panic(false)

	}
}
