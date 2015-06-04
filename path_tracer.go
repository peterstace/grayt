package grayt

import "math/rand"

type Quality struct {
	SamplesPerPixel int
}

func TracerImage(cam Camera, entities []Entity, acc Accumulator, quality Quality) {
	pxPitch := 2.0 / float64(acc.wide)
	for pxX := 0; pxX < acc.wide; pxX++ {
		for pxY := 0; pxY < acc.high; pxY++ {
			//if pxX != 140 || pxY != 40 {
			//	continue
			//}
			for i := 0; i < quality.SamplesPerPixel; i++ {
				x := (float64(pxX-acc.wide/2) + rand.Float64()) * pxPitch
				y := (float64(pxY-acc.high/2) + rand.Float64()) * pxPitch * -1.0
				r := cam.MakeRay(x, y)
				r.Dir = r.Dir.Unit()
				acc.add(pxX, pxY, tracePath(entities, r))
			}
		}
	}
}

func tracePath(entities []Entity, r Ray) Colour {

	intersection, hitEntity := closestHit(entities, r)
	if hitEntity == nil {
		return Colour{0, 0, 0}
	}

	// Since a 50/50 probability is used, don't bother scaling each colour by 2.
	switch rand.Int() % 2 {
	case 0:
		return hitEntity.Material.Colour.Scale(hitEntity.Material.Emittance)
	case 1:

		// Find where the ray hit. Reduce the intersection distance by a small
		// amount so that reflected rays don't intersect with it immediately.
		hitLoc := r.At(addULPs(intersection.Distance, -50))

		// Orient the unit normal towards the ray origin.
		if intersection.UnitNormal.Dot(r.Dir) > 0 {
			intersection.UnitNormal = Intersect.UnitNormal.Extended(-1.0)
		}

		// Create a random vector on the hemisphere towards the normal.
		rnd := Vect{rand.NormFloat64(), rand.NormFloat64(), rand.NormFloat64()}
		rnd = rnd.Unit()
		if rnd.Dot(intersection.UnitNormal) < 0 {
			rnd = rnd.Extended(-1.0)
		}

		// Apply the BRDF (bidirectional reflection distribution function).
		brdf := rnd.Dot(intersection.UnitNormal)

		return tracePath(entities, Ray{Start: hitLoc, Dir: rnd}).
			Scale(brdf).
			Mul(hitEntity.Material.Colour)
	default:
		panic("unexpected default case")
	}
}

func closestHit(entities []Entity, r Ray) (Intersection, *Entity) {
	var closest struct {
		intersection Intersection
		entity       *Entity
	}
	for i := range entities {
		intersection, hit := entities[i].Surface.Intersect(r)
		if !hit {
			continue
		}
		if closest.entity == nil || intersection.Distance < closest.intersection.Distance {
			closest.intersection = intersection
			closest.entity = &entities[i]
		}
	}
	return closest.intersection, closest.entity
}
