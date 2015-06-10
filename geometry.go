package grayt

import (
	"encoding/json"
	"errors"
)

// Intersection between a surface and a ray.
type Intersection struct {
	UnitNormal Vect    // Unit normal (pointing 'away' from the surface, not 'into' it).
	Distance   float64 // Distance along the ray where the intersection occurred.
}

// Surface is a two dimensional surface that can be intersected with a ray.
type Surface interface {

	// Intersect finds the intersection (if it exists) between a ray and the
	// surface.
	Intersect(Ray) (Intersection, bool)
}

type SurfaceFactory interface {
	MakeSurfaces() []Surface
}

type Material struct {
	Colour    Colour
	Emittance float64
	// Other properties such refractive index, reflectance, BRDF etc go here.
}

type Scene struct {
	CameraConfig CameraConfig
	Entities     []Entity
}

// Entity is a physical object whithin a scene.
type Entity struct {
	SurfaceFactories []SurfaceFactory
	Material         Material
}

func (e *Entity) UnmarshalJSON(p []byte) error {
	type record struct {
		RawSurfaces []rawSurface
		Material    Material
	}
	rec := new(record)
	if err := json.Unmarshal(p, &rec); err != nil {
		return err
	}
	e.Material = rec.Material
	for _, s := range rec.RawSurfaces {
		switch s.Type {
		case "Plane":
			var obj Plane
			if err := json.Unmarshal(s.Raw, &obj); err != nil {
				return err
			}
			e.SurfaceFactories = append(e.SurfaceFactories, obj)
		default:
			return errors.New("unknown type " + s.Type)
		}
	}
	return nil
}

type rawSurface struct {
	Type string
	Raw  json.RawMessage
}
