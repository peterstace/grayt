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

	var (
		white = grayt.Material{grayt.Colour{1, 1, 1}}
		green = grayt.Material{grayt.Colour{0, 1, 0}}
		red   = grayt.Material{grayt.Colour{1, 0, 0}}
		blue  = grayt.Material{grayt.Colour{0, 0, 1}}
	)

	scene := grayt.Scene{
		Camera: grayt.NewRectilinearCamera(grayt.CameraConfig{
			Location:      grayt.Vect{0.5, 0.5, 1.0},
			ViewDirection: grayt.Vect{0.0, 0.0, -1.0},
			UpDirection:   up,
			FieldOfView:   2*math.Asin(0.5/math.Sqrt(0.25+D)) + 0.1, // XXX why are we needing to add 0.1 here?
			FocalLength:   1.5,
			FocalRatio:    math.Inf(+1),
		}),
		Geometries: []grayt.Geometry{
			grayt.NewTriangle(blue, grayt.Vect{0, 0, -1}, grayt.Vect{1, 0, -1}, grayt.Vect{1, 0, 0}),
			grayt.NewPlane(white, up, zero),
			grayt.NewPlane(white, down, one),
			grayt.NewPlane(red, right, zero),
			grayt.NewPlane(green, left, one),
			grayt.NewPlane(white, back, one),
		},
		Lights: []grayt.Light{
			grayt.Light{
				Location:  zero.Add(one).Extended(0.5),
				Intensity: 0.3,
			},
		},
	}

	scene.Geometries = append(scene.Geometries, tallBlock(white)...)
	scene.Geometries = append(scene.Geometries, shortBlock(white)...)

	acc := grayt.NewAccumulator(300, 300)
	grayt.RayTracer(scene, acc)
	img := acc.ToImage(1.0)

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

func shortBlock(m grayt.Material) []grayt.Geometry {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = grayt.Vect{0.76, 0.00, -0.12} // 130,   0,  65
		LBB = grayt.Vect{0.85, 0.00, -0.41} //  82,   0, 225
		RBF = grayt.Vect{0.47, 0.00, -0.21} // 290,   0, 114
		RBB = grayt.Vect{0.56, 0.00, -0.49} // 240,   0, 272
		LTF = grayt.Vect{0.76, 0.30, -0.12} // 130, 165,  65
		LTB = grayt.Vect{0.85, 0.30, -0.41} //  82, 165, 225
		RTF = grayt.Vect{0.47, 0.30, -0.21} // 290, 165, 114
		RTB = grayt.Vect{0.56, 0.30, -0.49} // 240, 165, 272
	)
	var gs []grayt.Geometry
	gs = append(gs, grayt.NewSquare(m, LTF, LTB, RTB, RTF)...)
	gs = append(gs, grayt.NewSquare(m, LBF, RBF, RTF, LTF)...)
	gs = append(gs, grayt.NewSquare(m, LBB, RBB, RTB, LTB)...)
	gs = append(gs, grayt.NewSquare(m, LBF, LBB, LTB, LTF)...)
	gs = append(gs, grayt.NewSquare(m, RBF, RBB, RTB, RTF)...)
	return gs
}

func tallBlock(m grayt.Material) []grayt.Geometry {
	var (
		// Left/Right, Top/Bottom, Front/Back.
		LBF = grayt.Vect{0.52, 0.00, -0.54} // 265,   0, 296
		LBB = grayt.Vect{0.43, 0.00, -0.83} // 314,   0, 456
		RBF = grayt.Vect{0.23, 0.00, -0.45} // 423,   0, 247
		RBB = grayt.Vect{0.14, 0.00, -0.74} // 472,   0, 406
		LTF = grayt.Vect{0.52, 0.60, -0.54} // 265, 330, 296
		LTB = grayt.Vect{0.43, 0.60, -0.83} // 314, 330, 456
		RTF = grayt.Vect{0.23, 0.60, -0.45} // 423, 330, 247
		RTB = grayt.Vect{0.14, 0.60, -0.74} // 472, 330, 406
	)
	var gs []grayt.Geometry
	gs = append(gs, grayt.NewSquare(m, LTF, LTB, RTB, RTF)...)
	gs = append(gs, grayt.NewSquare(m, LBF, RBF, RTF, LTF)...)
	gs = append(gs, grayt.NewSquare(m, LBB, RBB, RTB, LTB)...)
	gs = append(gs, grayt.NewSquare(m, LBF, LBB, LTB, LTF)...)
	gs = append(gs, grayt.NewSquare(m, RBF, RBB, RTB, RTF)...)
	return gs
}
