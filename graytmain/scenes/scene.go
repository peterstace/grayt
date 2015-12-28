package scenes

import "github.com/peterstace/grayt/graytlib"

var SceneFactories = []SceneFactory{
	CornellBox{},
}

type SceneFactory interface {
	Name() string
	Scene() graytlib.Scene
}
