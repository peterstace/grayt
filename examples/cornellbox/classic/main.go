package main

import (
	. "github.com/peterstace/grayt/examples/cornellbox"
	. "github.com/peterstace/grayt/grayt"
)

func main() {
	Run("cornellbox_classic", Scene{
		Camera: Cam(1.3),
		Triangles: JoinTriangles(
			JoinTriangles(
				ShortBlock(),
				TallBlock(),
				Floor,
				Ceiling,
				BackWall,
			).SetColour(White),
			LeftWall.SetColour(Red),
			RightWall.SetColour(Green),
			CeilingLight(),
		),
	})
}
