package cornellbox

import (
	"math/rand"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/protocol"
	"github.com/peterstace/grayt/scenelib/dsl"
	"github.com/peterstace/grayt/xmath"
)

func Splitbox() protocol.Scene {
	cam := CornellCam(1.3)
	cam.LookingAt = xmath.Vect(0.5, initialBoxRadius.Y+0.03, -0.5)
	cam.FieldOfViewInRadians *= 0.5
	cam.AspectWide = 16
	cam.AspectHigh = 10

	redObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 0, 0}},
		Surface:  CornellLeftWall,
	}
	greenObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{0, 1, 0}},
		Surface:  CornellRightWall,
	}
	lights := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}, Emittance: 5},
		Surface:  CornellCeilingLight(),
	}
	whiteObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}},
		Surface: dsl.MergeSurfaces(
			CornellFloor,
			CornellCeiling,
			CornellBackWall,
			splitbox(),
		),
	}
	return protocol.Scene{
		Camera:  cam,
		Objects: []protocol.Object{redObjs, greenObjs, lights, whiteObjs},
	}
}

const numMovements = 100

var initialBoxRadius = xmath.Vect(0.22, 0.1, 0.1)

type box struct {
	min, max xmath.Vector
}

func splitbox() protocol.Surface {
	v1 := xmath.Vect(0.5-initialBoxRadius.X, 0, -0.5+initialBoxRadius.Z)
	v2 := xmath.Vect(0.5+initialBoxRadius.X, 2*initialBoxRadius.Y, -0.5-initialBoxRadius.Z)
	v1, v2 = v1.Min(v2), v1.Max(v2)
	boxes := []box{{v1, v2}}

	rnd := rand.New(rand.NewSource(0))
	for i := 0; i < numMovements; i++ {
		var newBoxes []box
		for _, box := range boxes {

			kind := rnd.Intn(6)
			fn := movements[kind]

			var splitLocation float64
			switch kind {
			case 0, 4:
				splitLocation = v1.X + (v2.X-v1.X)*rnd.Float64()
			case 2, 3:
				splitLocation = v1.Y + (v2.Y-v1.Y)*rnd.Float64()
			case 1, 5:
				splitLocation = v1.Z + (v2.Z-v1.Z)*rnd.Float64()
			default:
				panic(false)
			}

			splitAmount := (rnd.Float64() - 0.5) * 0.05
			splitBoxes := fn(splitLocation, splitAmount, box)
			newBoxes = append(newBoxes, splitBoxes...)
		}
		boxes = newBoxes
	}

	var surf protocol.Surface
	for _, b := range boxes {
		surf.AlignedBoxes = append(surf.AlignedBoxes, protocol.AlignedBox{
			CornerA: b.min,
			CornerB: b.max,
		})
	}
	return surf
}

func splitLeftRight(x float64, b box) (box, box) {
	b1 := box{b.min, xmath.Vect(x, b.max.Y, b.max.Z)}
	b2 := box{xmath.Vect(x, b.min.Y, b.min.Z), b.max}
	return b1, b2
}

func splitUpDown(y float64, b box) (box, box) {
	b1 := box{b.min, xmath.Vect(b.max.X, y, b.max.Z)}
	b2 := box{xmath.Vect(b.min.X, y, b.min.Z), b.max}
	return b1, b2
}

func splitFwdBack(z float64, b box) (box, box) {
	b1 := box{b.min, xmath.Vect(b.max.X, b.max.Y, z)}
	b2 := box{xmath.Vect(b.min.X, b.min.Y, z), b.max}
	return b1, b2
}

func heightMovementLeftRight(x float64, amount float64, input box) []box {
	if x < input.min.X || x > input.max.X {
		return []box{input}
	}
	b1, b2 := splitLeftRight(x, input)
	scale := amount / (2 * initialBoxRadius.Y)
	b1.min.Y *= 1 + scale
	b1.max.Y *= 1 + scale
	b2.min.Y *= 1 - scale
	b2.max.Y *= 1 - scale
	return []box{b1, b2}
}

func heightMovementFwdBack(z float64, amount float64, input box) []box {
	if z < input.min.Z || z > input.max.Z {
		return []box{input}
	}
	b1, b2 := splitFwdBack(z, input)
	scale := amount / (2 * initialBoxRadius.Y)
	b1.min.Y *= 1 + scale
	b1.max.Y *= 1 + scale
	b2.min.Y *= 1 - scale
	b2.max.Y *= 1 - scale
	return []box{b1, b2}
}

func layerMovementLeftRight(y float64, amount float64, input box) []box {
	if y < input.min.Y || y > input.max.Y {
		return []box{input}
	}
	b1, b2 := splitUpDown(y, input)
	b1.min.X += amount
	b1.max.X += amount
	b2.min.X -= amount
	b2.max.X -= amount
	return []box{b1, b2}
}

func layerMovementFwdBack(y float64, amount float64, input box) []box {
	if y < input.min.Y || y > input.max.Y {
		return []box{input}
	}
	b1, b2 := splitUpDown(y, input)
	b1.min.Z += amount
	b1.max.Z += amount
	b2.min.Z -= amount
	b2.max.Z -= amount
	return []box{b1, b2}
}

func shearFwdBack(x float64, amount float64, input box) []box {
	if x < input.min.X || x > input.max.X {
		return []box{input}
	}
	b1, b2 := splitLeftRight(x, input)
	b1.min.Z += amount
	b1.max.Z += amount
	b2.min.Z -= amount
	b2.max.Z -= amount
	return []box{b1, b2}
}

func shearLeftRight(z float64, amount float64, input box) []box {
	if z < input.min.Z || z > input.max.Z {
		return []box{input}
	}
	b1, b2 := splitFwdBack(z, input)
	b1.min.X += amount
	b1.max.X += amount
	b2.min.X -= amount
	b2.max.X -= amount
	return []box{b1, b2}
}

var movements = [...]func(float64, float64, box) []box{
	heightMovementLeftRight,
	heightMovementFwdBack,
	layerMovementLeftRight,
	layerMovementFwdBack,
	shearFwdBack,
	shearLeftRight,
}
