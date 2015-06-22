package grayt

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Material         Material
	SurfaceFactories []SurfaceFactory
}

func (e *Entity) UnmarshalJSON(p []byte) error {
	type record struct {
		Material         Material
		SurfaceFactories []json.RawMessage
	}
	rec := new(record)
	if err := json.Unmarshal(p, &rec); err != nil {
		return err
	}
	e.Material = rec.Material

	for _, raw := range rec.SurfaceFactories {
		t, err := getType(raw)
		if err != nil {
			return err
		}
		switch t {
		case "plane":
			var obj Plane
			if err := json.Unmarshal(raw, &obj); err != nil {
				return err
			}
			e.SurfaceFactories = append(e.SurfaceFactories, obj)
		case "sphere":
			var obj Sphere
			if err := json.Unmarshal(raw, &obj); err != nil {
				return err
			}
			e.SurfaceFactories = append(e.SurfaceFactories, obj)
		case "square":
			var obj Square
			if err := json.Unmarshal(raw, &obj); err != nil {
				return err
			}
			e.SurfaceFactories = append(e.SurfaceFactories, obj)
		case alignedBoxT:
			var obj AlignedBox
			if err := json.Unmarshal(raw, &obj); err != nil {
				return err
			}
			e.SurfaceFactories = append(e.SurfaceFactories, obj)
		default:
			return errors.New("unknown type " + t)
		}
	}

	return nil
}

func getType(raw json.RawMessage) (string, error) {
	record := struct{ Type string }{}
	if err := json.Unmarshal(raw, &record); err != nil {
		return "", err
	}
	if record.Type == "" {
		return "", fmt.Errorf("'Type' field is missing in %q", raw)
	}
	return record.Type, nil
}
