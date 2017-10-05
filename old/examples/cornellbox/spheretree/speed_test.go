package main

import (
	"testing"

	"github.com/peterstace/grayt/grayt"
)

func BenchmarkSplitbox(b *testing.B) {
	grayt.BenchmarkTraceImage(b, scene())
}
