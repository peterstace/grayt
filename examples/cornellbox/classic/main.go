package classic

import (
	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/protocol"
)

func Scene() protocol.Scene {
	cam := protocol.CornellCam(1.3)
	whiteObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}},
		Surface: protocol.MergeSurfaces(
			protocol.CornellShortBlock(),
			protocol.CornellTallBlock(),
			protocol.CornellFloor,
			protocol.CornellCeiling,
			protocol.CornellBackWall,
		),
	}
	redObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 0, 0}},
		Surface:  protocol.CornellLeftWall,
	}
	greenObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{0, 1, 0}},
		Surface:  protocol.CornellRightWall,
	}
	lights := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}, Emittance: 5},
		Surface:  protocol.CornellCeilingLight(),
	}

	return protocol.Scene{
		Camera: cam,
		Objects: []protocol.Object{
			whiteObjs,
			redObjs,
			greenObjs,
			lights,
		},
	}
}
