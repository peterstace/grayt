package grayt

type SceneFactory interface {
	FrameCount() int
	MakeScene(t float64) Scene
}

type Scene struct {
	Camera Camera
	Objs   []Obj
	Lights []Light
}
