package grayt

type RectilinearCamera struct {
	screenX, screenY Vect
	screenLoc        Vect
	eyeLoc           Vect
}

func NewRectilinearCamera() *RectilinearCamera {
	return nil
}

func (c *RectilinearCamera) MakeRay(x, y float64) ray {
	return ray{}
}
