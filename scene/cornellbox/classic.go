package cornellbox

import (
	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/scene"
	"github.com/peterstace/grayt/scene/dsl"
)

func Classic() scene.Scene {
	cam := CornellCam(1.3)
	whiteObjs := scene.Object{
		Material: scene.Material{Colour: colour.Colour{1, 1, 1}},
		Surface: dsl.MergeSurfaces(
			CornellShortBlock(),
			CornellTallBlock(),
			CornellFloor,
			CornellCeiling,
			CornellBackWall,
		),
	}
	redObjs := scene.Object{
		Material: scene.Material{Colour: colour.Colour{1, 0, 0}},
		Surface:  CornellLeftWall,
	}
	greenObjs := scene.Object{
		Material: scene.Material{Colour: colour.Colour{0, 1, 0}},
		Surface:  CornellRightWall,
	}
	lights := scene.Object{
		Material: scene.Material{Colour: colour.Colour{1, 1, 1}, Emittance: 5},
		Surface:  CornellCeilingLight(),
	}

	return scene.Scene{
		Camera: cam,
		Objects: []scene.Object{
			whiteObjs,
			redObjs,
			greenObjs,
			lights,
		},
	}
}
