package engine

import (
	"fmt"
	"image"
	"math/rand"
	"sync/atomic"
)

type Status struct {
	// Accessed atomically
	Done  int64
	Total int64
}

func TraceImage(pxWide int, scene func(*API), quality int, status *Status) image.Image {
	api := newAPI()
	scene(api)

	for _, o := range api.objs {
		fmt.Println(o)
	}

	pxHigh := pxWide * api.aspectRatio[1] / api.aspectRatio[0]

	cam := newCamera(api)
	accum := newAccumulator(pxWide, pxHigh)
	accel := accelList{api.objs}

	pxPitch := 2.0 / float64(pxWide)
	for q := 0; q < quality; q++ {
		rng := rand.New(rand.NewSource(int64(q)))
		for pxY := 0; pxY < pxHigh; pxY++ {
			for pxX := 0; pxX < pxWide; pxX++ {
				fmt.Printf("%d,%d\n", pxX, pxY)
				x := (float64(pxX-pxWide/2) + rng.Float64()) * pxPitch
				y := (float64(pxY-pxHigh/2) + rng.Float64()) * pxPitch * -1.0
				e, d := cam.makeRay(x, y, rng)
				d = d.unit()
				c := tracePath(&accel, e, d, rng)
				fmt.Println("  c", c)
				accum.add(pxX, pxY, c, q)
				atomic.AddInt64(&status.Done, 1)
			}
		}
	}
	return accum.toImage(1.0)
}

func tracePath(accel *accelList, e vect3, d vect3, rng *rand.Rand) vect3 {
	fmt.Println("  e", e)
	fmt.Println("  d", d)
	n, h, illum, hit := accel.closestHit(e, d)
	if !hit {
		return vect3{}
	}
	fmt.Println("  HIT")

	// Calculate probability of emitting.
	pEmit := 0.1
	if illum != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if rng.Float64() < pEmit {
		lvl := illum / pEmit
		return vect3{lvl, lvl, lvl}
	}

	// TODO: Offset the hit location by a small multiple of the normal so that
	// reflected rays don't intersect with it immediately.

	// Orient the unit normal towards the ray origin.
	if n.dot(d) > 0 {
		n = n.scale(-1)
	}

	// Create a random vector on the hemisphere towards the normal.
	rnd := vect3{
		rng.NormFloat64(),
		rng.NormFloat64(),
		rng.NormFloat64(),
	}
	rnd = rnd.unit()
	if rnd.dot(n) < 0 {
		rnd = rnd.scale(-1)
	}

	// Apply the BRDF (bidirectional reflection distribution function).
	brdf := rnd.dot(n)
	return tracePath(accel, h, rnd, rng).scale(brdf / (1 - pEmit))
}
