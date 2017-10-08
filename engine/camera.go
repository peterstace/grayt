package engine

import (
	"math"
	"math/rand"
)

type camera struct {
	screen, eye struct {
		// X vectors go from the center of the screen or eye to the right of
		// the screen or eye.  Y vectors go from the center of the screen or
		// eye towards the top of the screen or eye.  Loc is the location of
		// the center of the screen or eye.
		loc, x, y vect3
	}
}

func newCamera(a *API) camera {

	cam := camera{}

	upDirection := a.camUp.unit()
	viewDirection := a.camLook.sub(a.camLoc).unit()

	cam.screen.x = viewDirection.cross(upDirection)
	cam.screen.y = cam.screen.x.cross(viewDirection)

	cam.eye.x = cam.screen.x.scale(a.camFocalLen / a.camFocalRatio)
	cam.eye.y = cam.screen.y.scale(a.camFocalLen / a.camFocalRatio)
	cam.eye.loc = a.camLoc

	halfScreenWidth := math.Tan(a.camRadFOV/2) * a.camFocalLen
	cam.screen.x = cam.screen.x.scale(halfScreenWidth)
	cam.screen.y = cam.screen.y.scale(halfScreenWidth)
	cam.screen.loc = cam.eye.loc.add(viewDirection.scale(a.camFocalLen))

	return cam
}

func (c *camera) makeRay(x, y float64, rng *rand.Rand) (e, d vect3) {
	e = c.eye.loc.
		add(c.eye.x.scale(2*rng.Float64() - 1.0)).
		add(c.eye.y.scale(2*rng.Float64() - 1.0))
	scr := c.screen.loc.
		add(c.screen.x.scale(x)).
		add(c.screen.y.scale(y))
	d = scr.sub(e)
	return
}
