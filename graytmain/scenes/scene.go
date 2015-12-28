package scenes

import (
	"github.com/peterstace/grayt/graytlib"
	"github.com/peterstace/grayt/graytmain/scenes/cornellbox"
)

var RegisteredScenes = map[string]func() graytlib.Scene{
	"CornellBox": cornellbox.Scene,
}
