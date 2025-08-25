package main

import (
	"fmt"
	"strconv"

	"github.com/geofpwhite/numberPuzzleSolver/graph"
	"github.com/geofpwhite/numberPuzzleSolver/pairwise"
)

func main() {
	// size := flag.Int("size", 0, "give size ")
	// flag.Parse()
	// size := 4
	// dp := make(map[string]int, factorial(size*size))
remake:
	n := graph.RandomNewNode(3)
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
	// fmt.Println(n)
	inversions := 0
	zeroIndex := 0
	for i, row := range n.State {
		for j, num := range row {
			for l := j; l < len(n.State); l++ {
				num2 := n.State[i][l]
				if num > num2 {
					inversions++
				}
			}
			for k := i + 1; k < len(n.State); k++ {
				for l := 0; l < len(n.State); l++ {
					num2 := n.State[k][l]
					if num > num2 {
						inversions++
					}
				}
			}
			if num == 0 {
				zeroIndex = len(n.State) - i
			}
		}
	}
	fmt.Println(len(n.State)-zeroIndex, zeroIndex, inversions)
	switch len(n.State) % 2 {
	case 0:
		if (inversions+zeroIndex)%2 != 0 {
			goto remake
		}
	case 1:
		if inversions%2 != 0 {
			goto remake
		}
	}
	end := "1"
	for i := 2; i < len(n.State)*len(n.State); i++ {
		end = fmt.Sprintf("%s %s", end, strconv.Itoa(i))
	}
	// fmt.Println(end)
	// fmt.Println(solveIter(n).path)
	// fmt.Println(graph.IDAstar(n))
	x := pairwise.IDAstar(n)
	fmt.Println(x, len(x))
	// x := pairwise.SolveIter(n)
	// fmt.Println(x, len(x.Path))
	// fmt.Println(solve(n, []string{start}, m, start, end, 0))
	// fmt.Println(n2)
}
