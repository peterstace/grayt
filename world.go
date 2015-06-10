package grayt

type worldEntity struct {
	Surface  Surface
	Material Material
}

type World struct {
	entities []worldEntity
}

func (w *World) AddEntities(entities []Entity) {
	for _, e := range entities {
		for _, f := range e.SurfaceFactories {
			for _, s := range f.MakeSurfaces() {
				w.entities = append(w.entities, worldEntity{s, e.Material})
			}
		}
	}
}

func (w *World) closestHit(r Ray) (Intersection, *Material) {
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
