package main

import "testing"

// BenchmarkInk for benchmark
func BenchmarkInk(b *testing.B) {
	b.ReportAllocs()
	b.StartTimer()
	mdStream()
	b.StopTimer()
}
