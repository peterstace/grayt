package main

import (
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/peterstace/grayt"
)
import "flag"

func main() {

	var out string

	flag.StringVar(&out, "o", "", "output file (must end in .jpeg, .jpg, or .png)")
	flag.Parse()

	if getOutType(out) == outTypeUnknown {
		flag.Usage()
		return
	}

	// Cornell Box
	//
	// (0,y,-1) +---------+ (1,y,-1)
	//          | BB      |
	//          | BB  BB  |
	//          |     BB  |
	// (0,y,0)  +         + (1,y,0)
	//            \     /
	//             \   /
	//              \ /
	//               C (0.5,0.5,D)
	//
	//
	// SINE(0.5*T) = 0.5 / SQRT(0.5^2 + D)
	// T = 2*ARCSINE(0.5/SQRT(0.25 + D))

	const (
		D = 1.3
	)

	var (
		up    = grayt.Vect{0.0, 1.0, 0.0}
		down  = grayt.Vect{0.0, -1.0, 0.0}
		left  = grayt.Vect{-1.0, 0.0, 0.0}
		right = grayt.Vect{1.0, 0.0, 0.0}
		back  = grayt.Vect{0.0, 0.0, 1.0}
		zero  = grayt.Vect{0.0, 0.0, 0.0}
		one   = grayt.Vect{1.0, 1.0, -1.0}
	)

	scene := grayt.Scene{
		Camera: grayt.NewRectilinearCamera(grayt.CameraConfig{
			Location:      grayt.Vect{0.5, 0.5, 1.0},
			ViewDirection: grayt.Vect{0.0, 0.0, -1.0},
			UpDirection:   up,
			FieldOfView:   2 * math.Asin(0.5/math.Sqrt(0.25+D)),
			FocalLength:   1.5,
			FocalRatio:    math.Inf(+1),
		}),
		Geometries: []grayt.Geometry{
			grayt.NewPlane(up, zero),
			grayt.NewPlane(down, one),
			grayt.NewPlane(right, zero),
			grayt.NewPlane(left, one),
			grayt.NewPlane(back, one),
		},
		Lights: []grayt.Light{
			grayt.Light{
				Location:  zero.Add(one).Extended(0.5),
				Intensity: 0.3,
			},
		},
	}

	img := grayt.RayTracer(scene)

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	switch getOutType(out) {
	case outTypePNG:
		if err := png.Encode(f, img); err != nil {
			panic(err)
		}
	case outTypeJPEG:
		panic("unimplemented") // XXX
	case outTypeUnknown:
		panic("unknown outType should have been rejected during flag validation")
	}
}

type outType int

const (
	outTypeUnknown outType = iota
	outTypePNG
	outTypeJPEG
)

func getOutType(out string) outType {
	switch {
	case strings.HasSuffix(out, ".jpeg") || strings.HasSuffix(out, ".jpg"):
		return outTypeJPEG
	case strings.HasSuffix(out, ".png"):
		return outTypePNG
	default:
		return outTypeUnknown
	}
}
