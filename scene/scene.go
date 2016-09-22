package scene

import "io"

type Scene struct {
	Camera    Camera
	Triangles []Triangle
}

func (s Scene) WriteTo(w io.Writer) (n int64, err error) {
	// TODO
	return 0, nil
}

func ReadFrom(r io.Reader) (Scene, error) {
	// TODO
	return Scene{}, nil
}
