package pairwise

import "github.com/geofpwhite/numberPuzzleSolver/graph"

type qnode struct {
	distance int
	node     graph.Node
}

func Pairwise(size, a, b int) map[string]int {
	start := GeneratePairwiseSolvedState(size, a, b)
	distances := make(map[string]int)
	distances[start.String()] = 0
	queue := []qnode{{0, start}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		distances[cur.node.String()] = cur.distance
		neighbors := cur.node.DetermineNeighbors()
		for _, neighbor := range neighbors {
			if neighbor.String() == start.String() {
				continue
			}
			if distances[neighbor.String()] == 0 || distances[neighbor.String()] > cur.distance+1 {
				queue = append(queue, qnode{cur.distance + 1, neighbor})
			}
		}
	}
	return distances
}

func PairwiseSolutions(size int) map[graph.Coords]map[string]int {
	m := make(map[graph.Coords]map[string]int)
	for i := range size * size {
		for j := i + 1; j < size*size; j++ {
			m[graph.Coords{X: i, Y: j}] = Pairwise(size, i, j)
		}
	}
	return m
}

func GeneratePairwiseSolvedState(size, a, b int) graph.Node {
	n := graph.Node{State: make([][]int, size)}
	for i := range n.State {
		n.State[i] = make([]int, size)
	}
	cur := 1
	for i := range n.State {
		for j := range n.State[i] {
			n.State[i][j] = 0
			cur = (cur + 1) % (size * size)
			if cur == a {
				n.State[i][j] = a
			}
			if cur == b {
				n.State[i][j] = b
			}
		}
	}
	return n
}
