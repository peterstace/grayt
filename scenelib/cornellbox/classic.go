package cornellbox

import (
	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/protocol"
	"github.com/peterstace/grayt/scenelib/dsl"
)

func Classic() protocol.Scene {
	cam := CornellCam(1.3)
	whiteObjs := protocol.Object{
		Material: protocol.Material{Colour: colour.Colour{1, 1, 1}},
		Surface: dsl.MergeSurfaces(
			CornellShortBlock(),
			CornellTallBlock(),
			CornellFloor,
			CornellCeiling,
			CornellBackWall,
		),
	}
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
