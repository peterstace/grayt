package engine

import (
	"fmt"
	"math"
	"sort"
)

// RegisterScene registers a new scene. Panics if a scene with the same name
// has already been registered.
func RegisterScene(name string, fn func(*API)) {
	if _, ok := registered[name]; ok {
		panic(fmt.Sprintf("%q already registered", name))
	}
	registered[name] = fn
}

// LookupScene searches for a named scene.
func LookupScene(name string) (func(*API), bool) {
	fn, ok := registered[name]
	return fn, ok
}

func SceneList() []string {
	var scenes []string
	for s := range registered {
		scenes = append(scenes, s)
	}
	sort.Strings(scenes)
	return scenes
}

var registered map[string]func(*API) = make(map[string]func(*API))

func newAPI() *API {
	return &API{
		aspectRatio:   vect2{4, 3},
		camRadFOV:     0.5 * math.Pi,
		camFocalLen:   1.0,
		camFocalRatio: math.Inf(+1),
	}
}

type object struct {
	surf  surface
	illum float64
}

type API struct {
	aspectRatio vect2

	// TODO: Add checks to make sure loc, look, and up have been set.
	camLoc        vect3
	camLook       vect3
	camUp         vect3
	camRadFOV     float64
	camFocalLen   float64
	camFocalRatio float64

	illumination float64

	objs []object
}

// CameraLocation sets the cameras location. It must be set.
func (a *API) CameraLocation(x, y, z float64) {
	a.camLoc = vect3{x, y, z}
}

// CameraLookingAt sets the location where the camera looks at. It must be set.
func (a *API) CameraLookingAt(x, y, z float64) {
	a.camLook = vect3{x, y, z}
}

// CameraUpDirection sets the axial orientation of the camera. It must be set.
func (a *API) CameraUpDirection(x, y, z float64) {
	a.camUp = vect3{x, y, z}
}

// CameraDegFieldOfView sets the camera's horizontal field of view in degrees.
// Defaults to 90 degrees.
func (a *API) CameraDegFieldOfView(deg float64) {
	a.camRadFOV = deg * math.Pi / 180
}

// CameraRadFieldOfView sets the camera's horizontal field of view in radians.
// Defaults to 1/2 pi radians.
func (a *API) CameraRadFieldOfView(rad float64) {
	a.camRadFOV = rad
}

// CameraFocalLength sets the camera focal length. Objects exactly this
// distance from the camera will be in perfect focus.
func (a *API) CameraFocalLength(length float64) {
	a.camFocalLen = length
}

// CameraFocalRatio sets the camera focal ratio. This is the ratio between the
// camera's focal length and aperture width. Higher numbers results in a weaker
// depth of field effect, lower numbers result in a stronger depth of field
// effect. Defaults to +inf (no depth of field effect).
func (a *API) CameraFocalRatio(ratio float64) {
	a.camFocalRatio = ratio
}

// SetAspectRatio sets the aspect ratio of the output image to width:height.
// The default aspect ratio is 4:3.
func (a *API) SetAspectRatio(x, y int) {
	a.aspectRatio = vect2{x, y}
}

// func (a *API) Push() {
// }
// func (a *API) Pop() {
// }
// func (a *API) Translate(x, y, z float64) {
// }

// Tri adds a triangle with vertices a, b, c to the scene.
func (a *API) Tri(ax, ay, az, bx, by, bz, cx, cy, cz float64) {
	// TODO: If we already have value at the top of the stack, we would want to incorporate that *somehow*
	m := matrix4{
		vect4{ax, bx, cx, 0},
		vect4{ay, by, cy, 0},
		vect4{az, bz, cz, 0},
		vect4{1, 1, 1, 1},
	}
	a.objs = append(a.objs, object{
		surf:  surface{transform: m, primitive: triangle},
		illum: a.illumination,
	})
}

// Illuminate TODO
func (a *API) Illuminate(level float64) {
	a.illumination = level
}
