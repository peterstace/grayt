package engine

import (
	"testing"
)

func TestTriangle(t *testing.T) {
	for _, test := range []struct {
		e, d vect3
		h    vect3
		ok   bool
	}{
		{
			e:  vect3{5, 5, 5},
			d:  vect3{-invSqrt3, -invSqrt3, -invSqrt3},
			h:  vect3{1 / 3.0, 1 / 3.0, 1 / 3.0},
			ok: true,
		},
	} {
		gotN, gotH, gotOK := triangle(test.e, test.d)
		if gotOK != test.ok {
			t.Errorf("got=%v want=%v", gotOK, test.ok)
		}
		if !test.ok {
			continue
		}
		wantN := vect3{
			invSqrt3,
			invSqrt3,
			invSqrt3,
		}
		if wantN.sub(gotN).norm2() > 0.001 {
			t.Errorf("want=%v got=%v", wantN, gotN)
		}
		if test.h.sub(gotH).norm2() > 0.001 {
			t.Errorf("want=%v got=%v", test.h, gotH)
		}
	}
}

func TestTransformTriangle(t *testing.T) {
	for i, test := range []struct {
		triPos [9]float64
		e, d   vect3
		h, n   vect3
		ok     bool
	}{
		{
			triPos: [9]float64{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			},
			e:  vect3{5, 5, 5},
			d:  vect3{-1, -1, -1},
			h:  vect3{1 / 3.0, 1 / 3.0, 1 / 3.0},
			n:  vect3{invSqrt3, invSqrt3, invSqrt3},
			ok: true,
		},
		//{
		//	triPos: [9]float64{
		//		0, 1, 0,
		//		1, 0, 0,
		//		0, 0, 0,
		//	},
		//	e: vect3{0.5, 0.5, -2},
		//	d: vect3{-0.1, -0.1, -1},
		//	// TODO: Work out what these should be after I've got the intersection working.
		//	//h:  vect3{1 / 3.0, 1 / 3.0, 1 / 3.0},
		//	//n:  vect3{invSqrt3, invSqrt3, invSqrt3},
		//	ok: true,
		//},
	} {
		a := new(API)
		a.Tri(
			test.triPos[0],
			test.triPos[1],
			test.triPos[2],
			test.triPos[3],
			test.triPos[4],
			test.triPos[5],
			test.triPos[6],
			test.triPos[7],
			test.triPos[8],
		)
		s := a.objs[0].surf
		gotN, gotH, gotOK := s.intersect(test.e, test.d)
		if gotOK != test.ok {
			t.Errorf("%d: got=%v want=%v", i, gotOK, test.ok)
		}
		if !test.ok {
			continue
		}
		if test.n.sub(gotN).norm2() > 0.001 {
			t.Errorf("%d: want=%v got=%v", i, test.n, gotN)
		}
		if test.h.sub(gotH).norm2() > 0.001 {
			t.Errorf("%d: want=%v got=%v", i, test.h, gotH)
		}
	}
}
