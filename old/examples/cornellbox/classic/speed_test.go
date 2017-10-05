package main

import (
	"testing"

	"github.com/peterstace/grayt/grayt"
)

func BenchmarkClassic(b *testing.B) {
	grayt.BenchmarkTraceImage(b, scene())
}
