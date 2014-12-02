package grayt

type SceneFactory interface {
	FrameCount() int
	MakeScene(t float64) Scene
}

type Scene struct {
	Camera interface {
		MakeRay(x, y float64) ray
	}
	Objs   []Obj
	Lights []Light
}
