package grayt

import (
	"math"
	"math/rand"
	"testing"
)

func TestMatrixMul(t *testing.T) {
	m1 := translation(2, 5, 1)
	m2 := translation(4, 2, 4)
	want := matrix4{
		vect4{1, 0, 0, 6},
		vect4{0, 1, 0, 7},
		vect4{0, 0, 1, 5},
		vect4{0, 0, 0, 1},
	}
	got := m1.mulm(m2)
	if want != got {
		t.Errorf("Want=%v Got=%v", want, got)
	}
}

func TestMatrixInvKnown(t *testing.T) {
	m := matrix4{
		vect4{2, 2, 3, 4},
		vect4{5, 6, 5, 12},
		vect4{9, 11, 11, 12},
		vect4{13, 14, 15, 16},
	}
	want := matrix4{
		vect4{-68, 12, -96, 80},
		vect4{-56, 4, 68, -40},
		vect4{92, -28, 24, -20},
		vect4{18, 13, -4, -5},
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			want[i][j] /= 100
		}
	}
	got, ok := m.inv()
	if !ok {
		t.Errorf("couldn't invert")
	}
	if got != want {
		t.Errorf("Want=%v Got=%v", want, got)
	}
}

func TestMatrixInvReverse(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	var m matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m[i][j] = rnd.Float64()
		}
	}
	minv, ok := m.inv()
	if !ok {
		t.Errorf("Could not invert")
	}
	got := minv.mulm(m)
	const eps = 0.001
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				if math.Abs(1-got[i][j]) > eps {
					t.Errorf("Bad value")
				}
			} else {
				if math.Abs(got[i][j]) > eps {
					t.Errorf("Bad value")
				}
			}
		}
	}
}

func TestMatrixMulVector(t *testing.T) {
	m := matrix4{
		vect4{1, 2, 3, 4},
		vect4{5, 6, 7, 8},
		vect4{9, 0, 1, 2},
		vect4{3, 4, 5, 6},
	}
	v := vect4{1, 2, 3, 4}
	got := m.mulv(v)
	want := vect4{30, 70, 20, 50}
	if want != got {
		t.Errorf("Got=%v Want=%v", got, want)
	}
}

func BenchmarkMatrixMul(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	var m1, m2 matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m1[i][j] = rnd.Float64()
			m2[i][j] = rnd.Float64()
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m3 := m1.mulm(m2)
		_ = m3
	}
}

func BenchmarkMatrixInv(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	var m matrix4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m[i][j] = rnd.Float64()
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv, _ := m.inv()
		_ = inv
	}
}
