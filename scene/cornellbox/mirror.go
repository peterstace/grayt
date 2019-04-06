package cornellbox

import (
	"github.com/peterstace/grayt/scene"
	. "github.com/peterstace/grayt/scene/dsl"
)

func Mirror() scene.Scene {
	const (
		d = 1.3
		e = 0.05
	)
	return scene.Scene{
		Camera: CornellCam(d),
		Objects: []scene.Object{
			scene.Object{
				Surface: MergeSurfaces(
					// Ceiling
					AlignedSquare(Vect(0, 1, d), Vect(1, 1, -1)),
					// Floor
					AlignedSquare(Vect(0, 0, d), Vect(1, 0, -1)),
					// Back wall
					AlignedSquare(Vect(0, 0, -1), Vect(e, 1, -1)),
					AlignedSquare(Vect(1, 0, -1), Vect(1-e, 1, -1)),
					AlignedSquare(Vect(e, 1-e, -1), Vect(1-e, 1, -1)),
					AlignedSquare(Vect(e, 0, -1), Vect(1-e, e, -1)),
					// Front wall
					AlignedSquare(Vect(0, 0, d), Vect(1, 1, d)),
					// Front side walls
					AlignedSquare(Vect(0, 0, 0), Vect(0, 1, d)),
					AlignedSquare(Vect(1, 0, 0), Vect(1, 1, d)),
				),
				Material: scene.Material{
					Colour: White,
				},
			},
			scene.Object{
				Surface: CornellCeilingLight(),
				Material: scene.Material{
					Colour:    White,
					Emittance: 5.0,
				},
			},
			scene.Object{
				Surface: MergeSurfaces(
					AlignedSquare(Vect(e, e, -1), Vect(1-e, 1-e, -1)),
					AlignedSquare(Vect(0, e, -e), Vect(0, 1-e, -1+e)),
					AlignedSquare(Vect(1, e, -e), Vect(1, 1-e, -1+e)),
					Sphere(Vect(0.2, 0.1, -0.2), 0.1),
					Sphere(Vect(0.85, 0.1, -0.85), 0.1),
					Sphere(Vect(0.33, 0.6+0.075, -0.64), 0.075),
				),
				Material: scene.Material{
					Mirror: true,
				},
			},
			scene.Object{
				Surface: MergeSurfaces(
					AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -e)),
					AlignedSquare(Vect(0, 0, -1+e), Vect(0, 1, -1)),
					AlignedSquare(Vect(0, 0, -e), Vect(0, e, -1+e)),
					AlignedSquare(Vect(0, 1-e, -e), Vect(0, 1, -1+e)),
				),
				Material: scene.Material{
					Colour: Red,
				},
			},
			scene.Object{
				Surface: MergeSurfaces(
					AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -e)),
					AlignedSquare(Vect(1, 0, -1+e), Vect(1, 1, -1)),
					AlignedSquare(Vect(1, 0, -e), Vect(1, e, -1+e)),
					AlignedSquare(Vect(1, 1-e, -e), Vect(1, 1, -1+e)),
				),
				Material: scene.Material{
					Colour: Green,
				},
			},
			scene.Object{
				Surface:  CornellShortBlock(),
				Material: scene.Material{Colour: Hex(0xd684ff)},
			},
			scene.Object{
				Surface:  CornellTallBlock(),
				Material: scene.Material{Colour: Hex(0x3ae3ff)},
			},
		},
	}
}
