package grayt

import "testing"

func BenchmarkTraceImage(b *testing.B, s Scene) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TraceImage(16, s, 1, 1, new(uint64))
	}
}
