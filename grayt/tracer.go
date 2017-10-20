package grayt

import (
	"fmt"
	"image"
	"math/rand"
	"sync"
	"sync/atomic"
)

var aoeu = false

func TraceImage(pxWide int, scene Scene, quality, numWorkers int, completed *uint64) image.Image {

	pxHigh := scene.Camera.pxHigh(pxWide)

	cam := newCamera(scene.Camera)
	accum := newAccumulator(pxWide, pxHigh)
	accel := newGrid(4, scene.Objects)

	var wg sync.WaitGroup
	wg.Add(quality)

	sem := make(chan struct{}, numWorkers)

	pxPitch := 2.0 / float64(pxWide)
	for q := 0; q < quality; q++ {
		sem <- struct{}{}
		go func(q int) {
			tr := tracer{
				accel: accel,
				sky:   scene.Sky,
				rng:   rand.New(rand.NewSource(int64(q))),
			}
			for pxY := 0; pxY < pxHigh; pxY++ {
				for pxX := 0; pxX < pxWide; pxX++ {

					aoeu = false
					if pxX == 424 && pxY == 308 {
						fmt.Println()
						fmt.Println("DARK")
						fmt.Println()
						aoeu = true

					}
					if pxX == 419 && pxY == 308 {
						fmt.Println()
						fmt.Println("LIGHT")
						fmt.Println()
						aoeu = true
					}

					x := (float64(pxX-pxWide/2) + tr.rng.Float64()) * pxPitch
					y := (float64(pxY-pxHigh/2) + tr.rng.Float64()) * pxPitch * -1.0
					r := cam.makeRay(x, y, tr.rng)

					r.dir = r.dir.Unit()

					var c Colour
					if !*normals {
						c = tr.tracePath(r)
					} else {
						c = tr.traceNormal(r)
					}
					accum.add(pxX, pxY, c, q)
					atomic.AddUint64(completed, 1)
				}
			}
			<-sem
			wg.Done()
		}(q)
	}
	wg.Wait()

	if !*normals {
		return accum.toImage(1.0)
	} else {
		return accum.toRawImage()
	}
}

type tracer struct {
	sky   func(Vector) Colour
	accel accelerationStructure
	rng   *rand.Rand
}

func (t *tracer) tracePath(r ray) Colour {
	if aoeu {
		fmt.Println()
		fmt.Println("r:", r)
	}

	intersection, material, hit := t.accel.closestHit(r)
	if !hit {
		if t.sky == nil {
			return Colour{0, 0, 0}
		}
		assertUnit(r.dir)
		return t.sky(r.dir)
	}
	assertUnit(intersection.unitNormal)

	if aoeu {
		fmt.Println("inters:", intersection)
		fmt.Println("material:", material)
		fmt.Println("hit:", hit)
	}

	// Calculate probability of emitting.
	pEmit := 0.1
	if material.emittance != 0 {
		pEmit = 1.0
	}

	// Handle emit case.
	if t.rng.Float64() < pEmit {
		return material.colour.scale(material.emittance / pEmit)
	}

	// Find where the ray hit. Reduce the intersection distance by a small
	// amount so that reflected rays don't intersect with it immediately.
	hitLoc := r.at(addULPs(intersection.distance, -ulpFudgeFactor))

	//light := Vect(0, 10, 0)
	//toLight := light.Sub(hitLoc).Unit()
	//str := toLight.Dot(intersection.unitNormal)
	//str = math.Pow(str, 6)
	//return Colour{str, str, str}

	// Orient the unit normal towards the ray origin.
	if intersection.unitNormal.Dot(r.dir) > 0 {
		intersection.unitNormal = intersection.unitNormal.Scale(-1.0)
	}

	if material.mirror {

		reflected := r.dir.Sub(intersection.unitNormal.Scale(2 * intersection.unitNormal.Dot(r.dir)))
		return t.tracePath(ray{start: hitLoc, dir: reflected})

	} else {

		// Create a random vector on the hemisphere towards the normal.
		rnd := Vector{t.rng.NormFloat64(), t.rng.NormFloat64(), t.rng.NormFloat64()}
		rnd = rnd.Unit()
		if rnd.Dot(intersection.unitNormal) < 0 {
			rnd = rnd.Scale(-1.0)
		}
		if aoeu {
			fmt.Println("rnd:", rnd)
		}

		// Apply the BRDF (bidirectional reflection distribution function).
		brdf := rnd.Dot(intersection.unitNormal)

		return t.tracePath(ray{start: hitLoc, dir: rnd}).
			scale(brdf / (1 - pEmit)).
			mul(material.colour)
	}
}

func (t *tracer) traceNormal(r ray) Colour {
	intersection, _, hit := t.accel.closestHit(r)
	if !hit {
		return Colour{0.5, 0.5, 0.5}
	}
	norm := intersection.unitNormal.Add(Vect(1, 1, 1)).Scale(0.5)
	return Colour{norm.X, norm.Y, norm.Z}
}
