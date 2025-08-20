package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func main() {
	// size := flag.Int("size", 0, "give size ")
	// flag.Parse()
	// size := 4
	// dp := make(map[string]int, factorial(size*size))
	n := randomNewNode(4)
	// n := node{
	// 	state: [][]int{
	// 		// 		{7, 3, 1}, {0, 2, 8}, {6, 5, 4},
	// 		{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 0, 15},
	// 		// {1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 0, 14, 15},
	// 	},

	// }
	// n := node{
	// 	state: [][]int{
	// 		// 		{7, 3, 1}, {0, 2, 8}, {6, 5, 4},
	// 		{1, 2, 6}, {4, 5, 3}, {7, 8, 0},
	// 	},
	// }
	// n2 := n
	fmt.Println(n)
	end := "1"
	for i := 2; i < len(n.state)*len(n.state); i++ {
		end = fmt.Sprintf("%s %s", end, strconv.Itoa(i))
	}
	// fmt.Println(end)
	fmt.Println(solveIter(n).path)
	// fmt.Println(solve(n, []string{start}, m, start, end, 0))
	// fmt.Println(n2)
}

func factorial(num int) int {
	ret := 1
	for i := 2; i <= num; i++ {
		ret *= i
	}
	return ret
}

type node struct {
	state [][]int
}

type coords struct{ x, y int }

func (n node) String() string {
	var builder strings.Builder
	builder.WriteString("[\n")
	for _, row := range n.state {
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

func (n node) determineNeighbors() []node {
	neighbors := make([]node, 0, 4)
	var c00rds coords
outer:
	for i, row := range n.state {
		for j, num := range row {
			if num == 0 {
				c00rds.x, c00rds.y = i, j
				break outer
			}
		}
	}
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if (i == 0 && j == 0) || (i != 0 && j != 0) {
				continue
			}
			if c00rds.x+i >= 0 && c00rds.x+i < len(n.state) && c00rds.y+j >= 0 && c00rds.y+j < len(n.state) {
				neighborState := make([][]int, len(n.state))
				for i := range neighborState {
					neighborState[i] = make([]int, len(n.state))
					copy(neighborState[i], n.state[i])
				}
				neighborState[c00rds.x+i][c00rds.y+j], neighborState[c00rds.x][c00rds.y] = 0, neighborState[c00rds.x+i][c00rds.y+j]
				neighbors = append(neighbors, node{state: neighborState})
			}
		}
	}
	// for _, row := range n.state {
	// 	fmt.Println(row)
	// }
	return neighbors
}

func randomNewNode(size int) node {
	var n node = node{state: make([][]int, size)}
	for i := range n.state {
		n.state[i] = make([]int, size)
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
			n.state[i][j] = newState[index]
			index++
		}
	}
	return n
}

func (n node) manhattan(other node) int {
	sum := 0
	for i, row := range other.state {
		for j, num := range row {
		outer:
			for k, row := range n.state {
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

type queueNode struct {
	n         node
	path      []string
	steps     int
	manhatlen int
}

func solveIter(start node) queueNode {
	// inversions := 0
	// zeroIndex := 0
	// for i, row := range start.state {
	// 	for j, num := range row {
	// 		for l := j; l < len(start.state); l++ {
	// 			num2 := start.state[i][l]
	// 			if num > num2 {
	// 				inversions++
	// 			}
	// 		}
	// 		for k := i + 1; k < len(start.state); k++ {
	// 			for l := 0; l < len(start.state); l++ {
	// 				num2 := start.state[k][l]
	// 				if num > num2 {
	// 					inversions++
	// 				}
	// 			}
	// 		}
	// 		if num == 0 {
	// 			zeroIndex = len(start.state) - i
	// 		}
	// 	}
	// }
	// fmt.Println(len(start.state)-zeroIndex, zeroIndex, inversions)
	// switch len(start.state) % 2 {
	// case 0:
	// 	if (inversions+zeroIndex)%2 != 0 {
	// 		return queueNode{}
	// 	}
	// case 1:
	// 	if inversions%2 != 0 {
	// 		return queueNode{}
	// 	}
	// }
	startStr := start.String()
	solvedState := make([][]int, len(start.state))
	index := 1
	for i := range solvedState {
		solvedState[i] = make([]int, len(start.state))
		for j := range solvedState[i] {
			solvedState[i][j] = index
			index++
		}
	}
	solvedState[len(solvedState)-1][len(solvedState)-1] = 0
	solvedString := node{solvedState}.String()
	solvedNode := node{state: solvedState}
	startNode := queueNode{start, []string{}, 0, start.manhattan(solvedNode)}
	queue := priorityQueue{startNode}
	paths := make(map[string]int, factorial(len(start.state)*len(start.state)))
	minMattanFound := 0
	minState := make([][]int, 0)
	visited := make(map[string]bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to signal program termination

	// Goroutine to handle the signal
	var cur queueNode
	go func() {
		sig := <-c
		fmt.Printf("\nReceived signal: %s\n", sig)
		// Perform cleanup or other actions here
		fmt.Println("Cleaning up resources...")
		fmt.Println(cur.n.state)
		fmt.Println(minMattanFound)
		os.Exit(1)
	}()
	for len(queue) > 0 {
		// cur := queue[0]
		cur = queue.pop()
		if visited[cur.n.String()] {
			continue
		}
		visited[cur.n.String()] = true
		// fmt.Println(cur.n.String())
		// fmt.Println(len(queue), minMattanFound, cur.manhatlen)
		// queue = queue[1:]
		// mc := cur.n.manhattan(solvedNode)
		if cur.n.String() == solvedString {
			// fmt.Println(cur.n.manhattan(solvedNode))
			return cur
		}
		// queue = queue[1:]
		neighbors := cur.n.determineNeighbors()
		slices.SortFunc(neighbors, func(a, b node) int {
			ma, mb := a.manhattan(solvedNode), b.manhattan(solvedNode)
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
			newPath := append([]string{}, cur.path...)
			if (paths[stateStr] == 0 || cur.steps+1 < paths[stateStr]) && stateStr != startStr {
				paths[stateStr] = cur.steps + 1
				newPath = append(newPath, stateStr)
				newNode := queueNode{n: neighbor, path: newPath, steps: cur.steps + 1, manhatlen: neighbor.manhattan(solvedNode)}
				// fmt.Println("neighbor", newNode.manhatlen)
				if newNode.manhatlen == 0 {
					fmt.Println(newNode.manhatlen)
				}
				if newNode.manhatlen < minMattanFound || minMattanFound == 0 {
					minMattanFound = newNode.manhatlen
					minState = newNode.n.state
				}
				// queue = append(queue, newNode)
				queue.insert(newNode)
				// fmt.Println(len(queue))
			}
		}
	}
	fmt.Println(minMattanFound, minState)
	return queueNode{}
}

type priorityQueue []queueNode

func (pq *priorityQueue) insert(node queueNode) {
	(*pq) = append((*pq), node)
	nodeIndex := len((*pq)) - 1
	for nodeIndex > 0 {
		parentIndex := nodeIndex / 2
		if node.manhatlen > (*pq)[parentIndex].manhatlen {
			break
		}
		(*pq)[parentIndex], (*pq)[nodeIndex] = (*pq)[nodeIndex], (*pq)[parentIndex]
		nodeIndex = parentIndex
	}
}

func (pq *priorityQueue) pop() queueNode {
	ret := (*pq)[0]
	(*pq)[0] = (*pq)[len((*pq))-1]
	(*pq).bubbleDown()
	*pq = (*pq)[:len(*pq)-1]
	return ret
}

func (pq *priorityQueue) bubbleDown() {
	fmt.Println("bubblingdown")
	if len((*pq)) == 1 {
		return
	}
	if len((*pq)) == 2 {
		cur := (*pq)[0]
		child := (*pq)[1]
		mc, mch := cur.manhatlen, child.manhatlen
		if mc > mch {
			(*pq)[0], (*pq)[1] = (*pq)[1], (*pq)[0]
		}
		return
	}
	curIndex := 0
	cur := (*pq)[curIndex]
	for curIndex*2+2 < len((*pq)) {
		child1, child2 := (*pq)[(curIndex*2)+1], (*pq)[(curIndex*2)+2]
		mc, mc1, mc2 := cur.manhatlen, child1.manhatlen, child2.manhatlen
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
