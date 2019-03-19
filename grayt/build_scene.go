package grayt

import "github.com/peterstace/grayt/protocol"

func buildScene(proto protocol.Scene) Scene {
	var scene Scene

	// TODO: This doesn't make much sense. Should instead store the fully
	// formed camera object on the scene?
	scene.Camera.Location = proto.Camera.Location
	scene.Camera.LookingAt = proto.Camera.LookingAt
	scene.Camera.UpDirection = proto.Camera.UpDirection
	scene.Camera.FieldOfViewInRadians = proto.Camera.FieldOfViewInRadians
	scene.Camera.FocalLength = proto.Camera.FocalLength
	scene.Camera.FocalRatio = proto.Camera.FocalRatio
	scene.Camera.AspectWide = proto.Camera.AspectWide
	scene.Camera.AspectHigh = proto.Camera.AspectHigh

	scene.Objects = make([]Object, len(proto.Objects))
	for i := range proto.Objects {
		// TODO: Should just use the material object in the scene (until a
		// point where we need some additional fields).
		scene.Objects[i].Material.Colour = proto.Objects[i].Material.Colour
		scene.Objects[i].Material.Emittance = proto.Objects[i].Material.Emittance
		scene.Objects[i].Material.Mirror = proto.Objects[i].Material.Mirror
		scene.Objects[i].Surface = buildSurface(proto.Objects[i].Surface)
	}
	return scene
}

func buildSurface(proto interface{}) surface {
	switch o := proto.(type) {
	case protocol.Triangle:
		return newTriangle(o.A, o.B, o.C)
	case protocol.AlignedBox:
		return newAlignedBox(o.CornerA, o.CornerB)
	case protocol.Sphere:
		return &sphere{Center: o.Center, Radius: o.Radius}
	case protocol.AlignXSquare:
		return &alignXSquare{o.X, o.Y1, o.Y2, o.Z1, o.Z2}
	case protocol.AlignYSquare:
		return &alignYSquare{o.X1, o.X2, o.Y, o.Z1, o.Z2}
	case protocol.AlignZSquare:
		return &alignZSquare{o.X1, o.X2, o.Y1, o.Y2, o.Z}
	case protocol.Disc:
		return &disc{Center: o.Center, RadiusSq: o.Radius * o.Radius, UnitNorm: o.UnitNorm}
	case protocol.Pipe:
		return &pipe{C1: o.EndpointA, C2: o.EndpointB, R: o.Radius}
	default:
		// TODO: Handle this a bit better
		panic("unknown type")
	}
}
