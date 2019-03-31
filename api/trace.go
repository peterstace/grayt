package api

import (
	"encoding/binary"
	"io"
	"math/rand"
	"time"

	"github.com/peterstace/grayt/colour"
	"github.com/peterstace/grayt/trace"
)

func traceLayer(w io.Writer, pxWide, pxHigh, depth int, accel trace.AccelerationStructure, cam trace.Camera) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tr := trace.NewTracer(accel, rng)
	pxPitch := 2.0 / float64(pxWide)
	for pxY := 0; pxY < pxHigh; pxY++ {
		for pxX := 0; pxX < pxWide; pxX++ {
			var c colour.Colour
			for i := 0; i < depth; i++ {
				x := (float64(pxX-pxWide/2) + rng.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rng.Float64()) * pxPitch * -1.0
				cr := cam.MakeRay(x, y, rng)
				cr.Dir = cr.Dir.Unit()
				c = c.Add(tr.TracePath(cr))
			}
			binary.Write(w, binary.BigEndian, c)
		}
	}
}
