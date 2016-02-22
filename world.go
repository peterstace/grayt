package grayt

type world struct {
	entities []Entity
}

func newWorld(entities []Entity) *world {
	return &world{entities}
}

func (w *world) closestHit(r Ray) (Intersection, *Material) {
	var closest struct {
		Intersection Intersection
		Material     *Material
	}
	for i := range w.entities {
		intersection, hit := w.entities[i].Surface.Intersect(r)
		if !hit {
			continue
		}
		if closest.Material == nil || intersection.Distance < closest.Intersection.Distance {
			closest.Intersection = intersection
			closest.Material = &w.entities[i].Material
		}
	}
	return closest.Intersection, closest.Material
}
