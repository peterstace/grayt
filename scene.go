package grayt

type Scene struct {
	Camera     *camera
	Geometries []geometry
	Lights     []Light
}
