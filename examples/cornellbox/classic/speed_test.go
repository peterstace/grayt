package main

import (
	"testing"

	. "github.com/peterstace/grayt/grayt"
)

func BenchmarkSpeed(b *testing.B) {
	s := scene()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TraceImage(16, s, 1, 1, new(uint64))
	}
}
