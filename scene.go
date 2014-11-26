package grayt

type Scene struct {
	Camera interface {
		MakeRay(x, y float64) ray
	}
	Objs   []Obj
	Lights []Light
}
