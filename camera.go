package grayt

// Cameras produce rays that go from an eye to a virtual screen. The rays
// produced are specified via a coordiate system on the virtual screen.  The
// left side of the virtual screen has x coordinate -1, the right side of the
// virtual screen has coordinate +1. The top of the virtual screen has y
// coordinate v and the bottom of the virtual screen has y coordinate -v (where
// the value of v depends on the aspect ratio of the screen).
type Camera interface {
	MakeRay(x, y float64) ray
}

func pixelCoordsToCameraCoords(pxX, pxY int) (x, y float64) {
	return 0.0, 0.0
}

type RectilinearCamera struct {
	screenX, screenY Vect
	screenLoc        Vect
	eyeLoc           Vect
}

func NewRectilinearCamera() Camera {
	return nil
}

func (c *RectilinearCamera) MakeRay(x, y float64) ray {
	return ray{}
}
