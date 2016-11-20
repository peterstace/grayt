package main

import (
	"math"
	"math/rand"

	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("sphere_tree", Scene{
		Camera: Cam(1.3),
		Objects: Group(
			Tree(),
			Floor,
			Ceiling,
			BackWall,
			LeftWall.With(ColourRGB(Red)),
			RightWall.With(ColourRGB(Green)),
			CeilingLight().With(Emittance(5.0)),
		),
	})
}
func Cam(d float64) Camera {
	loc := Vect(0.5, 0.5, d)
	at := Vect(0.5, 0.25, -0.5)
	dir := at.Sub(loc)
	return Camera{
		Location:             loc,
		ViewDirection:        dir,
		UpDirection:          Vect(0.0, 1.0, 0.0),
		FieldOfViewInDegrees: 0.95 * 2 * math.Asin(0.5/math.Sqrt(0.25+d*d)) * 180 / math.Pi,
		FocalLength:          0.5 + d,
		FocalRatio:           math.Inf(+1),
	}
}

var (
	Floor     = AlignedSquare(Vect(0, 0, 0), Vect(1, 0, -1))
	Ceiling   = AlignedSquare(Vect(0, 1, 0), Vect(1, 1, -1))
	BackWall  = ZPlane(-1)
	LeftWall  = AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -1))
	RightWall = AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -1))
)

func CeilingLight() ObjectList {
	const size = 0.9
	return AlignedBox(
		Vect(size, 1.0, -size),
		Vect(1.0-size, 0.999, -1.0+size),
	)
}

type sphere struct {
	c Vector
	r float64
}

func Tree() ObjectList {

	root := sphere{Vect(0.5, 0, -0.5), 0.2}
	spheres := new([]sphere)
	*spheres = append(*spheres, root)
	recurse(spheres, root, 9)

	var objList ObjectList
	for _, s := range *spheres {
		objList = Group(objList, Sphere(s.c, s.r))
	}
	return objList
}

const radiusScaleDown = 0.7

func recurse(spheres *[]sphere, parent sphere, level int) {

	if level == 0 {
		return
	}

	child1, child2 := findChildren(spheres, parent)

	*spheres = append(*spheres, child1)
	*spheres = append(*spheres, child2)

	recurse(spheres, child1, level-1)
	recurse(spheres, child2, level-1)
}

func findChildren(spheres *[]sphere, parent sphere) (sphere, sphere) {
	var child1, child2 sphere
	for true {
		child1 = createChild(parent)
		child2 = createChild(parent)
		if !isValidChild(child1, parent, spheres) {
			continue
		}
		if !isValidChild(child2, parent, spheres) {
			continue
		}
		if spheresIntersect(child1, child2) {
			continue
		}
		break
	}
	return child1, child2
}

var rnd = rand.New(rand.NewSource(0))

func createChild(parent sphere) sphere {
	rndUnit := Vector{rnd.NormFloat64(), rnd.NormFloat64(), rnd.NormFloat64()}.Unit()
	return sphere{
		parent.c.Add(rndUnit.Scale(parent.r)),
		radiusScaleDown * parent.r,
	}
}

func isValidChild(child, parent sphere, spheres *[]sphere) bool {

	// Check for intersection with other spheres (ignore the parent).
	for _, s := range *spheres {
		if s.c == parent.c && s.r == parent.r {
			continue
		}
		if spheresIntersect(s, child) {
			return false
		}
	}

	// Check for wall/floor/ceiling intersection.
	return true &&
		child.c.X > child.r &&
		child.c.X < 1.0-child.r &&
		child.c.Y > child.r &&
		child.c.Y < 1.0-child.r &&
		child.c.Z < -child.r &&
		child.c.Z > -1.0+child.r
}

func spheresIntersect(s1, s2 sphere) bool {
	return s1.c.Sub(s2.c).Length() < s1.r+s2.r
}
