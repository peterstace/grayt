package grayt

import (
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
	got := m1.mul(m2)
	if want != got {
		t.Errorf("Want=%v Got=%v", want, got)
	}
}

func TestMatrixInv(t *testing.T) {
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
		m3 := m1.mul(m2)
		_ = m3
	}
}
