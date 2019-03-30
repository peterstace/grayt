package cornellbox

import (
	"github.com/peterstace/grayt/protocol"
	. "github.com/peterstace/grayt/scenelib/dsl"
	. "github.com/peterstace/grayt/xmath"
)

func Mirror() protocol.Scene {
	const (
		d = 1.3
		e = 0.05
	)
	return protocol.Scene{
		Camera: CornellCam(d),
		Objects: []protocol.Object{
			protocol.Object{
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
				Material: protocol.Material{
					Colour: White,
				},
			},
			protocol.Object{
				Surface: CornellCeilingLight(),
				Material: protocol.Material{
					Colour:    White,
					Emittance: 5.0,
				},
			},
			protocol.Object{
				Surface: MergeSurfaces(
					AlignedSquare(Vect(e, e, -1), Vect(1-e, 1-e, -1)),
					AlignedSquare(Vect(0, e, -e), Vect(0, 1-e, -1+e)),
					AlignedSquare(Vect(1, e, -e), Vect(1, 1-e, -1+e)),
					Sphere(Vect(0.2, 0.1, -0.2), 0.1),
					Sphere(Vect(0.85, 0.1, -0.85), 0.1),
					Sphere(Vect(0.33, 0.6+0.075, -0.64), 0.075),
				),
				Material: protocol.Material{
					Mirror: true,
				},
			},
			protocol.Object{
				Surface: MergeSurfaces(
					AlignedSquare(Vect(0, 0, 0), Vect(0, 1, -e)),
					AlignedSquare(Vect(0, 0, -1+e), Vect(0, 1, -1)),
					AlignedSquare(Vect(0, 0, -e), Vect(0, e, -1+e)),
					AlignedSquare(Vect(0, 1-e, -e), Vect(0, 1, -1+e)),
				),
				Material: protocol.Material{
					Colour: Red,
				},
			},
			protocol.Object{
				Surface: MergeSurfaces(
					AlignedSquare(Vect(1, 0, 0), Vect(1, 1, -e)),
					AlignedSquare(Vect(1, 0, -1+e), Vect(1, 1, -1)),
					AlignedSquare(Vect(1, 0, -e), Vect(1, e, -1+e)),
					AlignedSquare(Vect(1, 1-e, -e), Vect(1, 1, -1+e)),
				),
				Material: protocol.Material{
					Colour: Green,
				},
			},
			protocol.Object{
				Surface:  CornellShortBlock(),
				Material: protocol.Material{Colour: Hex(0xd684ff)},
			},
			protocol.Object{
				Surface:  CornellTallBlock(),
				Material: protocol.Material{Colour: Hex(0x3ae3ff)},
			},
		},
	}
}
