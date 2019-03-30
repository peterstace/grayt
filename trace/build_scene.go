package trace

import "github.com/peterstace/grayt/scene"

type Scene struct {
	Camera  Camera   `json:"camera"`
	Objects []Object `json:"objects"`
}

func BuildScene(proto scene.Scene) Scene {
	scene := Scene{
		Camera: NewCamera(proto.Camera),
	}
	for _, o := range proto.Objects {
		add := func(s surface) {
			scene.Objects = append(scene.Objects, Object{
				Surface: s,
				Material: material{
					Colour:    o.Material.Colour,
					Emittance: o.Material.Emittance,
					Mirror:    o.Material.Mirror,
				},
			})
		}
		for _, x := range o.Surface.Triangles {
			add(newTriangle(x.A, x.B, x.C))
		}
		for _, x := range o.Surface.AlignedBoxes {
			add(newAlignedBox(x.CornerA, x.CornerB))
		}
		for _, x := range o.Surface.Spheres {
			add(&sphere{Center: x.Center, Radius: x.Radius})
		}
		for _, x := range o.Surface.AlignXSquares {
			add(&alignXSquare{x.X, x.Y1, x.Y2, x.Z1, x.Z2})
		}
		for _, x := range o.Surface.AlignYSquares {
			add(&alignYSquare{x.X1, x.X2, x.Y, x.Z1, x.Z2})
		}
		for _, x := range o.Surface.AlignZSquares {
			add(&alignZSquare{x.X1, x.X2, x.Y1, x.Y2, x.Z})
		}
		for _, x := range o.Surface.Discs {
			add(&disc{Center: x.Center, RadiusSq: x.Radius * x.Radius, UnitNorm: x.UnitNorm})
		}
		for _, x := range o.Surface.Pipes {
			add(&pipe{C1: x.EndpointA, C2: x.EndpointB, R: x.Radius})
		}
	}
	return scene
}
