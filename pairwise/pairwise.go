package pairwise

import (
	"fmt"
	"slices"

	"github.com/geofpwhite/numberPuzzleSolver/graph"
)

type qnode struct {
	distance int
	node     graph.Node
}

func Pairwise(size, a, b int) map[[2]graph.Coords]int {
	start := GeneratePairwiseSolvedState(size, a, b)
	distances := make(map[[2]graph.Coords]int)
	distances[[2]graph.Coords{index(start, a), index(start, b)}] = 0
	queue := []qnode{{0, start}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		distances[[2]graph.Coords{index(cur.node, a), index(cur.node, b)}] = cur.distance
		neighbors := cur.node.DetermineNeighbors()
		for _, neighbor := range neighbors {
			if neighbor.Equals(start) {
				continue
			}
			if distances[[2]graph.Coords{index(neighbor, a), index(neighbor, b)}] == 0 || distances[[2]graph.Coords{index(neighbor, a), index(neighbor, b)}] > cur.distance+1 {
				queue = append(queue, qnode{cur.distance + 1, neighbor})
			}
		}
	}
	return distances
}
func Pairwise2(size, a, b int, start graph.Node) map[[2]graph.Coords]int {
	distances := make(map[[2]graph.Coords]int)
	distances[[2]graph.Coords{index(start, a), index(start, b)}] = 0
	queue := []qnode{{0, start}}
	for len(queue) > 0 {
		// fmt.Println(len(queue), "q")
		cur := queue[0]
		queue = queue[1:]
		curAIndex := index(cur.node, a)
		curBIndex := index(cur.node, b)
		distances[[2]graph.Coords{curAIndex, curBIndex}] = cur.distance
		neighbors := cur.node.DetermineNeighbors()
		// fmt.Println(len(neighbors), "neighbors")
		for _, neighbor := range neighbors {
			if neighbor.Equals(start) {
				continue
			}
			if distances[[2]graph.Coords{index(neighbor, a), index(neighbor, b)}] == 0 || distances[[2]graph.Coords{index(neighbor, a), index(neighbor, b)}] > cur.distance+1 {
				queue = append(queue, qnode{cur.distance + 1, neighbor})
			}
		}
	}
	return distances
}

func PairwiseSolutions(size int) map[graph.Coords]map[[2]graph.Coords]int {
	m := make(map[graph.Coords]map[[2]graph.Coords]int)
	for i := range size * size {
		for j := i + 1; j < size*size; j++ {
			m[graph.Coords{X: i, Y: j}] = Pairwise(size, i, j)
		}
	}
	return m
}
func PairwiseSolutions2(size int, node graph.Node) map[graph.Coords]map[[2]graph.Coords]int {
	m := make(map[graph.Coords]map[[2]graph.Coords]int)
	for i := range size * size {
		for j := i + 1; j < size*size; j++ {
			m[graph.Coords{X: i, Y: j}] = Pairwise2(size, i, j, node)
		}
	}
	return m
}
func index(node graph.Node, num int) graph.Coords {
	for i, row := range node.State {
		for j, val := range row {
			if val == num {
				return graph.Coords{X: i, Y: j}
			}
		}
	}
	return graph.Coords{X: -1, Y: -1}
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
			if cur == a {
				n.State[i][j] = a
			}
			if cur == b {
				n.State[i][j] = b
			}
			cur = (cur + 1) % (size * size)
		}
	}
	return n
}

type hnode struct {
	value     int
	neighbors map[*hnode]int
}

func ComputeHeuristicValue(node graph.Node, solutions map[graph.Coords]map[[2]graph.Coords]int) int {
	hGraph := make(map[int]*hnode)
	num := 1
	for i := range len(node.State) {
		for j := range len(node.State) {
			hGraph[num] = &hnode{node.State[i][j], make(map[*hnode]int)}
			num = (num + 1) % (len(node.State) * len(node.State))
			// println(num)
		}
	}
	for i := range len(node.State) * len(node.State) {
		for j := i + 1; j < len(node.State)*len(node.State); j++ {
			ii := index(node, i)
			ij := index(node, j)
			// fmt.Println(ii, ij, i, j, hGraph[i], hGraph[j])
			hGraph[i].neighbors[hGraph[j]] = solutions[graph.Coords{X: i, Y: j}][[2]graph.Coords{ii, ij}]
			hGraph[j].neighbors[hGraph[i]] = solutions[graph.Coords{X: i, Y: j}][[2]graph.Coords{ii, ij}]
		}
	}
	return MaxSumOfPairwise(hGraph)
}

// given graph of tiles with all pairwise distances between tiles, find max sum of edges such that no 2 edges are adjacent to the same vertex
func MaxSumOfPairwise(hGraph map[int]*hnode) int {
	// fmt.Println(hGraph)
	// for k, v := range hGraph {
	// 	fmt.Println(k, v)
	// }
	LU := make([]int, 0)
	for _, node := range hGraph {
		lu := 0
		for _, distance := range node.neighbors {
			if distance > lu {
				lu = distance
			}
		}
		LU = append(LU, (lu))
	}
	LV := make([]int, len(hGraph))
	matchV := make([]int, len(hGraph))
	for i := range matchV {
		matchV[i] = -1
	}
	for u := range len(hGraph) {
		slack := make([]int, len(hGraph))
		slackx := make([]int, len(hGraph))
		prev := make([]int, len(hGraph))
		S, T := make([]bool, len(hGraph)), make([]bool, len(hGraph))
		for i := range slack {
			slack[i] = 1_000_000_000_000
			slackx[i] = -1
			prev[i] = -1
		}
		queue := make([]int, 1)
		queue[0] = (u)
		S[u] = true
		prevV := -1
		for {
			for len(queue) > 0 {
				i := int(queue[0])
				queue = queue[1:]
				for j := range len(hGraph) {
					if !T[j] {
						// fmt.Println(i, j, hGraph[i], hGraph[j])
						gap := LU[i] + LV[j] - (hGraph[i].neighbors[hGraph[j]])
						if gap == 0 {
							T[j] = true
							if matchV[j] == -1 {
								prevV = j
								break
							} else {
								queue = append(queue, matchV[j])
								S[int(matchV[j])] = true
								prev[int(matchV[j])] = (i)
							}

						} else if slack[j] > (gap) {
							slack[j] = (gap)
							slackx[j] = (i)
							// fmt.Println(slackx)
						}
					}
				}
				if prevV != -1 {
					break
				}
			}
			if prevV != -1 {
				j := prevV
				// fmt.Println("j=", j)
				for j != -1 {
					i := slackx[j]
					tmp := matchV[j]
					matchV[j] = i
					// fmt.Printf("matchV[%d]=%d\n", j, i)

					j = int(tmp)
				}
				// fmt.Println("breaking")
				break
			} else {
				delta := -1
				for i := range slack {
					if !T[i] && (slack[i] < (delta) || delta == -1) {
						delta = slack[i]
					}
				}
				for i := range slack {
					if S[i] {
						LU[i] -= delta
					}
				}
				for i := range slack {
					if T[i] {
						LV[i] += delta
					} else {
						slack[i] -= delta
					}
				}
				for j := range slack {
					if !T[j] && slack[j] == 0 {
						T[j] = true
						if matchV[j] == -1 {
							prevV = j
							break
						} else {
							queue = append(queue, matchV[j])
							S[int(matchV[j])] = true
							prev[int(matchV[j])] = slackx[j]
						}
					}
				}
				if prevV != -1 {
					j := prevV
					// fmt.Println(prevV)
					for j != -1 {
						i := slackx[j]
						tmp := matchV[j]
						matchV[j] = i
						// fmt.Printf("\n matchV[%d]=%d\n", j, i)
						// fmt.Println(matchV)
						j = int(tmp)
					}
					// fmt.Println("breaking")
					break
				}
			}

		}

	}
	// fmt.Println(matchV)
	sum := 0
	for pair1, pair2 := range matchV {
		sum += hGraph[pair1].neighbors[hGraph[int(pair2)]]
		// fmt.Println(pair1, pair2)
	}
	// fmt.Println(sum)
	return sum
}

func IDAstar(root graph.Node) []graph.Node {
	solvedState := make([][]int, len(root.State))
	path := []graph.Node{root}
	index := 1
	for i := range solvedState {
		solvedState[i] = make([]int, len(root.State))
		for j := range solvedState[i] {
			solvedState[i][j] = index
			index++
		}
	}
	solvedState[len(solvedState)-1][len(solvedState)-1] = 0
	goal := graph.Node{solvedState}
	pwiseSolutions := PairwiseSolutions(len(root.State))
	bound := ComputeHeuristicValue(root, pwiseSolutions)
	fmt.Println(pwiseSolutions)
	for {
		t, cost := search(&path, 0, bound, goal, pwiseSolutions)
		if t != nil {
			return t
		}
		if cost == -1 {
			return nil
		}
		bound = cost
		fmt.Println("searched")
		fmt.Println(path, t, cost)
	}
}

func search(path *[]graph.Node, g, bound int, goal graph.Node, pwiseSolutions map[graph.Coords]map[[2]graph.Coords]int) ([]graph.Node, int) {
	// fmt.Println("searched")
	n := (*path)[len((*path))-1]
	f := g + ComputeHeuristicValue(n, pwiseSolutions)
	if f > bound {
		return nil, f
	}
	if n.IsGoal(goal) {
		return (*path), bound
	}
	min := -1
	var t []graph.Node = nil
	for _, succ := range n.DetermineNeighbors() {
		if slices.ContainsFunc((*path), func(n graph.Node) bool { return n.Equals(succ) }) {
			continue
		}
		(*path) = append((*path), succ)
		// x := PairwiseSolutions2(len(goal.State), succ)
		// t, cost := search(path, max(g+ComputeHeuristicValue(n, x), g+n.ManhattanSum(succ)), bound, goal, pwiseSolutions)
		t, cost := search(path, g+n.ManhattanSum(succ)+1, bound, goal, pwiseSolutions)
		if t != nil {
			return t, cost
		}
		if cost < min || min == -1 {
			min = cost
		}
		*path = (*path)[:len(*path)-1]
	}
	return t, min
}

func SolveIter(start graph.Node) graph.QueueNode {
	startStr := start.String()
	solvedState := make([][]int, len(start.State))
	index := 1
	for i := range solvedState {
		solvedState[i] = make([]int, len(start.State))
		for j := range solvedState[i] {
			solvedState[i][j] = index
			index++
		}
	}
	solvedState[len(solvedState)-1][len(solvedState)-1] = 0
	solvedString := graph.Node{State: solvedState}.String()
	// solvedNode := graph.Node{State: solvedState}
	solutions := PairwiseSolutions(len(solvedState))
	// fmt.Println(solutions)
	startNode := graph.QueueNode{N: start, Path: []string{}, Steps: 0, Manhatlen: ComputeHeuristicValue(start, solutions)}
	queue := graph.PriorityQueue{startNode}
	paths := make(map[string]int, graph.Factorial(len(start.State)*len(start.State)))
	minMattanFound := 0
	minState := make([][]int, 0)
	visited := make(map[string]bool)

	// Create a channel to signal program termination

	// Goroutine to handle the signal
	var cur graph.QueueNode
	for len(queue) > 0 {
		// cur := queue[0]
		cur = queue.Pop()
		if visited[cur.N.String()] {
			continue
		}
		visited[cur.N.String()] = true
		// fmt.Println(cur.n.String())
		// fmt.Println(len(queue), minMattanFound, cur.manhatlen)
		// queue = queue[1:]
		// mc := cur.n.Manhattan(solvedNode)
		if cur.N.String() == solvedString {

			// fmt.Println(cur.n.Manhattan(solvedNode))
			return cur
		}
		// queue = queue[1:]
		neighbors := cur.N.DetermineNeighbors()
		slices.SortFunc(neighbors, func(a, b graph.Node) int {
			ma, mb := ComputeHeuristicValue(a, solutions), ComputeHeuristicValue(b, solutions)
			// fmt.Println(ma, mb, solutions)
			switch {
			case ma > mb:
				return 1
			case ma < mb:
				return -1
			default:
				return 0
			}
		})
		for _, neighbor := range neighbors {
			stateStr := neighbor.String()
			newPath := append([]string{}, cur.Path...)
			if (paths[stateStr] == 0 || cur.Steps+1 < paths[stateStr]) && stateStr != startStr {
				paths[stateStr] = cur.Steps + 1
				newPath = append(newPath, stateStr)

				newNode := graph.QueueNode{N: neighbor, Path: newPath, Steps: cur.Steps + 1, Manhatlen: ComputeHeuristicValue(neighbor, solutions)}

				if newNode.Manhatlen < minMattanFound || minMattanFound == 0 {
					minMattanFound = newNode.Manhatlen
					minState = newNode.N.State
				}
				// queue = append(queue, newNode)
				queue.Insert(newNode)
				// fmt.Println(len(queue))
			}
		}
		// fmt.Println(minMattanFound, minState)
	}
	fmt.Println(minMattanFound, minState)
	return graph.QueueNode{}
}
