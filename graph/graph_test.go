package graph

import (
	"testing"
)

func TestSolveIter(t *testing.T) {
	start := Node{State: [][]int{
		{0, 1, 2},
		{4, 5, 3},
		{7, 8, 6},
	}}
	expected := QueueNode{
		N:         start,
		Path:      []string{},
		Steps:     0,
		Manhatlen: 0,
	}

	result := SolveIter(start)
	if !result.N.Equals(expected.N) {
		t.Errorf("SolveIter failed: expected %v, got %v", expected, result)
	}
}

func TestIDAstar(t *testing.T) {
	root := Node{State: [][]int{
		{0, 1, 2},
		{4, 5, 3},
		{7, 8, 6},
	}}

	result := IDAstar(root)
	if result == nil {
		t.Errorf("IDAstar failed: result is nil")
	}
	if len(result) < 1 {
		t.Errorf("IDAstar failed: expected path length >= 1, got %d", len(result))
	}
}
