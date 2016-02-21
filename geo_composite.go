package grayt

import "encoding/json"

const squareT = "square"

type Square struct {
	V1, V2, V3, V4 Vect
}

func (s Square) MarshalJSON() ([]byte, error) {
	type alias Square
	return json.Marshal(struct {
		Type string
		alias
	}{squareT, alias(s)})
}

func (s Square) MakeSurfaces() []Surface {
	return append(
		Triangle{s.V1, s.V2, s.V3}.MakeSurfaces(),
		Triangle{s.V3, s.V4, s.V1}.MakeSurfaces()...,
	)
}
