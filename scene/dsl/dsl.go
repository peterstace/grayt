package dsl

import (
	"math"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/xmath"
)

var (
	White = colour.Colour{1, 1, 1}
	Red   = colour.Colour{1, 0, 0}
	Green = colour.Colour{0, 1, 0}
	Blue  = colour.Colour{0, 0, 1}
)

func Vect(x, y, z float64) xmath.Vector {
	return xmath.Vector{x, y, z}
}

func Hex(col int32) colour.Colour {
	r := float64((col&0xff0000)>>0x10) / 0xff
	g := float64((col&0x00ff00)>>0x08) / 0xff
	b := float64((col&0x0000ff)>>0x00) / 0xff
	return colour.Colour{r, g, b}
}

func DefaultCamera() scene.Camera {
	return scene.Camera{
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

type ObjectList []scene.Object

func MergeSurfaces(surfs ...scene.Surface) scene.Surface {
	var all scene.Surface
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
	var objs []scene.Object
	for i := range lists {
		for j := range lists[i] {
			objs = append(objs, lists[i][j])
		}
	}
	return objs
}

func Square(a, b, c, d xmath.Vector) scene.Surface {
	return scene.Surface{
		Triangles: []scene.Triangle{
			{a, b, c},
			{c, d, a},
		},
	}
}

func AlignedSquare(a, b xmath.Vector) scene.Surface {
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
		return scene.Surface{AlignXSquares: []scene.AlignXSquare{{a.X, a.Y, b.Y, a.Z, b.Z}}}
	case a.Y == b.Y:
		return scene.Surface{AlignYSquares: []scene.AlignYSquare{{a.X, b.X, a.Y, a.Z, b.Z}}}
	case a.Z == b.Z:
		return scene.Surface{AlignZSquares: []scene.AlignZSquare{{a.X, b.X, a.Y, b.Y, a.Z}}}
	default:
		panic(false)

	}
}

func Sphere(center xmath.Vector, radius float64) scene.Surface {
	return scene.Surface{Spheres: []scene.Sphere{{center, radius}}}
}
