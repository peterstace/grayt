package cornellbox

import (
	"math/rand"

	"github.com/peterstace/grayt/scene"
	. "github.com/peterstace/grayt/scene/dsl"
	"github.com/peterstace/grayt/xmath"
)

func SphereTree() scene.Scene {
	cam := CornellCam(1.3)
	cam.LookingAt = Vect(0.5, 0.25, -0.5)
	cam.FieldOfViewInRadians *= 0.95
	cam.AspectWide = 2
	cam.AspectHigh = 1
	return scene.Scene{
		Camera: cam,
		Objects: []scene.Object{
			scene.Object{
				Material: scene.Material{Colour: White, Emittance: 5},
				Surface:  CornellCeilingLight(),
			},
			scene.Object{
				Material: scene.Material{Colour: White},
				Surface: MergeSurfaces(
					CornellFloor,
					CornellBackWall,
					CornellCeiling,
					tree(),
				),
			},
			scene.Object{
				Material: scene.Material{Colour: Red},
				Surface:  CornellLeftWall,
			},
			scene.Object{
				Material: scene.Material{Colour: Green},
				Surface:  CornellRightWall,
			},
		},
	}
}

type sphere struct {
	c xmath.Vector
	r float64
}

func tree() scene.Surface {

	root := sphere{xmath.Vect(0.5, 0, -0.5), 0.2}
	spheres := new([]sphere)
	*spheres = append(*spheres, root)
	recurse(spheres, root, 9)

	var surf scene.Surface
	for _, s := range *spheres {
		surf.Spheres = append(surf.Spheres, scene.Sphere{s.c, s.r})
	}
	return surf
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

// TODO: this should be reset at the start of each new scene generation
var rnd = rand.New(rand.NewSource(0))

func createChild(parent sphere) sphere {
	rndUnit := xmath.Vector{rnd.NormFloat64(), rnd.NormFloat64(), rnd.NormFloat64()}.Unit()
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
