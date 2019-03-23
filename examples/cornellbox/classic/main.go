package classic

import (
	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/protocol"
	"github.com/peterstace/grayt/scenelib/cornellbox"
	"github.com/peterstace/grayt/scenelib/dsl"
)

func Scene() protocol.Scene {
	cam := cornellbox.CornellCam(1.3)
	whiteObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}},
		Surface: dsl.MergeSurfaces(
			cornellbox.CornellShortBlock(),
			cornellbox.CornellTallBlock(),
			cornellbox.CornellFloor,
			cornellbox.CornellCeiling,
			cornellbox.CornellBackWall,
		),
	}
	redObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 0, 0}},
		Surface:  cornellbox.CornellLeftWall,
	}
	greenObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{0, 1, 0}},
		Surface:  cornellbox.CornellRightWall,
	}
	lights := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}, Emittance: 5},
		Surface:  cornellbox.CornellCeilingLight(),
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
