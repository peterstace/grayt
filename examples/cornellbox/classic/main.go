package classic

import (
	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/protocol"
)

func Scene() protocol.Scene {
	cam := protocol.CornellCam(1.3)
	whiteObjs := protocol.Combine(
		protocol.Material{Colour: colour.Colour{1, 1, 1}},
		protocol.CornellShortBlock(),
		protocol.CornellTallBlock(),
		protocol.CornellFloor,
		protocol.CornellCeiling,
		protocol.CornellBackWall,
	)
	redObjs := protocol.Combine(
		protocol.Material{Colour: colour.Colour{1, 0, 0}},
		protocol.CornellLeftWall,
	)
	greenObjs := protocol.Combine(
		protocol.Material{Colour: colour.Colour{0, 1, 0}},
		protocol.CornellRightWall,
	)
	lights := protocol.Combine(
		protocol.Material{Colour: colour.Colour{1, 1, 1}, Emittance: 5},
		protocol.CornellCeilingLight(),
	)

	return protocol.Scene{
		Camera: cam,
		Objects: protocol.MergeObjectLists(
			whiteObjs,
			redObjs,
			greenObjs,
			lights,
		),
	}
}
