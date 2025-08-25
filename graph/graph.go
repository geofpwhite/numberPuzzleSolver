package graph

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

func Factorial(num int) int {
	ret := 1
	for i := 2; i <= num; i++ {
		ret *= i
	}
	return ret
}

type Node struct {
	State [][]int
}

type Coords struct{ X, Y int }

func GenerateSolvedState(size int) Node {
	n := Node{make([][]int, size)}
	for i := range n.State {
		n.State[i] = make([]int, size)
	}

	cur := 1
	for i := range n.State {
		for j := range n.State[i] {
			n.State[i][j] = cur
			cur = (cur + 1) % (size * size)
		}
	}
	return n
}

func (n Node) String() string {
	var builder strings.Builder
	builder.WriteString("[\n")
	for _, row := range n.State {
		builder.WriteString("\t[")
		builder.Grow(int(unsafe.Sizeof(row)))
		for _, num := range row {
			builder.WriteString(fmt.Sprintf("%s ", strconv.Itoa(num)))
		}
		builder.WriteString("]\n")
	}
	builder.WriteString("]\n")
	s := builder.String()
	return s[:len(s)-1]
}

func (n Node) DetermineNeighbors() []Node {
	neighbors := make([]Node, 0, 4)
	var c00rds Coords
outer:
	for i, row := range n.State {
		for j, num := range row {
			if num == 0 {
				c00rds.X, c00rds.Y = i, j
				break outer
			}
		}
	}
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if (i == 0 && j == 0) || (i != 0 && j != 0) {
				continue
			}
			if c00rds.X+i >= 0 && c00rds.X+i < len(n.State) && c00rds.Y+j >= 0 && c00rds.Y+j < len(n.State) {
				neighborState := make([][]int, len(n.State))
				for i := range neighborState {
					neighborState[i] = make([]int, len(n.State))
					copy(neighborState[i], n.State[i])
				}
				neighborState[c00rds.X+i][c00rds.Y+j], neighborState[c00rds.X][c00rds.Y] = 0, neighborState[c00rds.X+i][c00rds.Y+j]
				neighbors = append(neighbors, Node{State: neighborState})
			}
		}
	}
	// for _, row := range n.state {
	// 	fmt.Println(row)
	// }
	return neighbors
}

func RandomNewNode(size int) Node {
	var n Node = Node{State: make([][]int, size)}
	for i := range n.State {
		n.State[i] = make([]int, size)
	}
	defaul := make([]int, size*size)
	newState := make([]int, 0, size*size)
	for i := range defaul {
		defaul[i] = i
	}
	for len(defaul) > 0 {
		next := rand.Intn(len(defaul))
		newState = append(newState, defaul[next])
		defaul = append(defaul[:next], defaul[next+1:]...)
	}
	index := 0
	for i := range size {
		for j := range size {
			n.State[i][j] = newState[index]
			index++
		}
	}
	return n
}

func (n Node) ManhattanSum(other Node) int {
	sum := 0
	for i, row := range other.State {
		for j, num := range row {
		outer:
			for k, row := range n.State {
				for l, num2 := range row {
					if num == num2 {
						sum += max(i-k, k-i) + max(j-l, l-j)
						break outer
					}
				}
			}
		}
	}
	return sum
}

type QueueNode struct {
	N         Node
	Path      []string
	Steps     int
	Manhatlen int
}

func (n Node) Equals(other Node) bool {
	for i, row := range n.State {
		for j, num := range row {
			if num != other.State[i][j] {
				return false
			}
		}
	}
	return true
}

func SolveIter(start Node) QueueNode {
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
	solvedString := Node{solvedState}.String()
	solvedNode := Node{State: solvedState}
	startNode := QueueNode{start, []string{}, 0, start.ManhattanSum(solvedNode)}
	queue := PriorityQueue{startNode}
	paths := make(map[string]int, Factorial(len(start.State)*len(start.State)))
	minMattanFound := 0
	minState := make([][]int, 0)
	visited := make(map[string]bool)

	// Create a channel to signal program termination

	// Goroutine to handle the signal
	var cur QueueNode
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
		slices.SortFunc(neighbors, func(a, b Node) int {
			ma, mb := a.ManhattanSum(solvedNode), b.ManhattanSum(solvedNode)
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
				newNode := QueueNode{N: neighbor, Path: newPath, Steps: cur.Steps + 1, Manhatlen: neighbor.ManhattanSum(solvedNode)}
				// fmt.Println("neighbor", newNode.manhatlen)
				if newNode.Manhatlen == 0 {
					fmt.Println(newNode.Manhatlen)
				}
				if newNode.Manhatlen < minMattanFound || minMattanFound == 0 {
					minMattanFound = newNode.Manhatlen
					minState = newNode.N.State
				}
				// queue = append(queue, newNode)
				queue.Insert(newNode)
				// fmt.Println(len(queue))
			}
		}
	}
	fmt.Println(minMattanFound, minState)
	return QueueNode{}
}

// IDAstar initializes and orchestrates the iterative deepening A* search.
// It returns the path from the root to the goal node, if one is found.
func IDAstar(root Node) []Node {
	solvedState := make([][]int, len(root.State))
	path := []Node{root}
	index := 1
	for i := range solvedState {
		solvedState[i] = make([]int, len(root.State))
		for j := range solvedState[i] {
			solvedState[i][j] = index
			index++
		}
	}
	solvedState[len(solvedState)-1][len(solvedState)-1] = 0
	goal := Node{solvedState}
	bound := root.ManhattanSum(goal)
	for {
		t, cost := search(&path, 0, bound, goal)
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

// Search is a recursive function that performs a depth-limited Search.
// It explores paths from the current node, pruning branches that exceed the given cost bound.
// It returns the path to the goal if found, and the minimum cost of a path that exceeded the bound.
func search(path *[]Node, g, bound int, goal Node) ([]Node, int) {
	n := (*path)[len((*path))-1]
	f := g + n.ManhattanSum(goal)
	if f > bound {
		return nil, f
	}
	if n.IsGoal(goal) {
		return (*path), bound
	}
	min := -1
	var t []Node = nil
	for _, succ := range n.DetermineNeighbors() {
		if slices.ContainsFunc((*path), func(n Node) bool { return n.String() == succ.String() }) {
			continue
		}
		(*path) = append((*path), succ)
		t, cost := search(path, g+n.ManhattanSum(succ)+1, bound, goal)
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

// IsGoal checks if the current node's state matches the goal state.
func (n Node) IsGoal(goal Node) bool {
	return n.String() == goal.String()
}

type PriorityQueue []QueueNode

func (pq *PriorityQueue) Insert(node QueueNode) {
	(*pq) = append((*pq), node)
	nodeIndex := len((*pq)) - 1
	for nodeIndex > 0 {
		parentIndex := nodeIndex / 2
		if node.Manhatlen > (*pq)[parentIndex].Manhatlen {
			break
		}
		(*pq)[parentIndex], (*pq)[nodeIndex] = (*pq)[nodeIndex], (*pq)[parentIndex]
		nodeIndex = parentIndex
	}
}

func (pq *PriorityQueue) Pop() QueueNode {
	ret := (*pq)[0]
	(*pq)[0] = (*pq)[len((*pq))-1]
	(*pq).BubbleDown()
	*pq = (*pq)[:len(*pq)-1]
	return ret
}

func (pq *PriorityQueue) BubbleDown() {
	if len((*pq)) == 1 {
		return
	}
	if len((*pq)) == 2 {
		cur := (*pq)[0]
		child := (*pq)[1]
		mc, mch := cur.Manhatlen, child.Manhatlen
		if mc > mch {
			(*pq)[0], (*pq)[1] = (*pq)[1], (*pq)[0]
		}
		return
	}
	curIndex := 0
	cur := (*pq)[curIndex]
	for curIndex*2+2 < len((*pq)) {
		child1, child2 := (*pq)[(curIndex*2)+1], (*pq)[(curIndex*2)+2]
		mc, mc1, mc2 := cur.Manhatlen, child1.Manhatlen, child2.Manhatlen
		if mc < mc1 && mc < mc2 {
			return
		}
		larger := 1
		if mc1 != max(mc1, mc2) {
			larger = 2
		}
		(*pq)[curIndex], (*pq)[curIndex*2+larger] = (*pq)[curIndex*2+larger], (*pq)[curIndex]
		curIndex *= 2
		curIndex += larger
	}
}
