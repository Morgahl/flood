package graph_test

import (
	"testing"

	"github.com/Morgahl/flood/internal/fixtures"
	"github.com/Morgahl/flood/pkg/graph"
)

func BenchmarkNew_ones(b *testing.B) { benchmarkNew(b, fixtures.Ones) }

func BenchmarkNew_test(b *testing.B) { benchmarkNew(b, fixtures.Test) }

func BenchmarkNew_test_layered(b *testing.B) { benchmarkNew(b, fixtures.TestLayered) }

func BenchmarkNew_first(b *testing.B) { benchmarkNew(b, fixtures.First) }

func BenchmarkNew_fifth(b *testing.B) { benchmarkNew(b, fixtures.Fifth) }

func BenchmarkNormalize_ones(b *testing.B) { benchmarkNormalize(b, fixtures.Ones) }

func BenchmarkNormalize_test(b *testing.B) { benchmarkNormalize(b, fixtures.Test) }

func BenchmarkNormalize_test_layered(b *testing.B) { benchmarkNormalize(b, fixtures.TestLayered) }

func BenchmarkNormalize_first(b *testing.B) { benchmarkNormalize(b, fixtures.First) }

func BenchmarkNormalize_fifth(b *testing.B) { benchmarkNormalize(b, fixtures.Fifth) }

func BenchmarkRunCopy_ones(b *testing.B) {
	benchmarkRunCopy(b, fixtures.Ones)
}

func BenchmarkRunCopy_test(b *testing.B) {
	benchmarkRunCopy(b, fixtures.Test)
}

func BenchmarkRunCopy_test_layered(b *testing.B) {
	benchmarkRunCopy(b, fixtures.TestLayered)
}

func BenchmarkRunCopy_first(b *testing.B) {
	benchmarkRunCopy(b, fixtures.First)
}

func BenchmarkRunCopy_fifth(b *testing.B) {
	benchmarkRunCopy(b, fixtures.Fifth)
}

func BenchmarkRunCopyTwice_ones(b *testing.B) {
	benchmarkRunCopyTwice(b, fixtures.Ones)
}

func BenchmarkRunCopyTwice_test(b *testing.B) {
	benchmarkRunCopyTwice(b, fixtures.Test)
}

func BenchmarkRunCopyTwice_test_layered(b *testing.B) {
	benchmarkRunCopyTwice(b, fixtures.TestLayered)
}

func BenchmarkRunCopyTwice_first(b *testing.B) {
	benchmarkRunCopyTwice(b, fixtures.First)
}

func BenchmarkRunCopyTwice_fifth(b *testing.B) {
	benchmarkRunCopyTwice(b, fixtures.Fifth)
}

func BenchmarkRunSolution_test(b *testing.B) {
	benchmarkRunSolution(b, fixtures.Test, fixtures.TestSolution)
}

func BenchmarkRunSolution_test_layered(b *testing.B) {
	benchmarkRunSolution(b, fixtures.TestLayered, fixtures.TestLayeredSolution)
}

func BenchmarkRunSolution_first(b *testing.B) {
	benchmarkRunSolution(b, fixtures.First, fixtures.FirstSolution)
}

func BenchmarkRunSolution_fifth(b *testing.B) {
	benchmarkRunSolution(b, fixtures.Fifth, fixtures.FifthSolution)
}

func benchmarkNew(b *testing.B, base [19][19]uint8) {
	for b.Loop() {
		graph.New(base)
	}
}

func benchmarkNormalize(b *testing.B, base [19][19]uint8) {
	g := graph.New(base)
	var ng *graph.Graph
	for b.Loop() {
		ng = g.Copy()
		ng.Normalize()
	}
	ng.Normalize()
}

func benchmarkRunCopy(b *testing.B, base [19][19]uint8) {
	g := graph.New(base)
	g.Normalize()
	var ng *graph.Graph
	for b.Loop() {
		ng = g.Copy()
	}
	ng.Normalize()
}

func benchmarkRunCopyTwice(b *testing.B, base [19][19]uint8) {
	g := graph.New(base)
	g.Normalize()
	g = g.Copy()
	var ng *graph.Graph
	for b.Loop() {
		ng = g.Copy()
	}
	ng.Normalize()
}

func benchmarkRunSolution(b *testing.B, base [19][19]uint8, solution []uint8) {
	g := graph.New(base)
	g.Normalize()
	for b.Loop() {
		ng := g.Copy()
		for _, v := range solution {
			ng.Flood(v)
		}
	}
}
