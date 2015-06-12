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
	return []Surface{
		NewTriangle(s.V1, s.V2, s.V3),
		NewTriangle(s.V3, s.V4, s.V1),
	}
}
