package astar

import (
	"container/heap"
	"math"
	"slices"
)

type Point struct {
	X, Y int
}

type Node struct {
	Pos    Point
	G      float64
	H      float64
	F      float64
	Parent *Node
	Index  int
}

// ─────────────────────────────────────────────
// Priority Queue (min-heap)
// ─────────────────────────────────────────────

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].F < pq[j].F }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Push(x any) {
	n := x.(*Node)
	n.Index = len(*pq)
	*pq = append(*pq, n)
}
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := old[len(old)-1]
	old[len(old)-1] = nil
	*pq = old[:len(old)-1]
	return n
}

// ─────────────────────────────────────────────
// Heuristic Functions
// ─────────────────────────────────────────────

type HeuristicFn func(a,b Point) float64

func Manhanttan(a, b Point) float64 {
	return math.Abs(float64(a.X - b.X)) + math.Abs(float64(a.Y - b.Y))
}

func Euclidean(a, b Point) float64 {
	x := float64(a.X - b.X)
	y := float64(a.Y - b.Y)
	return math.Sqrt(x*x + y*y)
}

func Chebyshev(a, b Point) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return math.Max(dx, dy)
}

func Octile(a, b Point) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return dx + dy - (2-math.Sqrt2)*math.Min(dx, dy)
}

func Zero(a,b Point) float64 {return 0}


var directions4 = []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
var directions8 = []Point{
	{0, 1}, {0, -1}, {1, 0}, {-1, 0},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

type AstarAns struct {
	Path []Point
	NodeExp int
	Found bool
}

func Astar(grid [][]int, start, goal Point, h HeuristicFn, diagonalAllowed bool) AstarAns{
	rows := len(grid)
	cols := len(grid[0])
	
	inBounds := func(a Point) bool { return a.X >= 0 && a.X < cols && a.Y >= 0 && a.Y < rows}
	
	canGo := func(a Point) bool { return grid[a.X][a.Y] == 0}
	
	startNode := &Node{ Pos: start, G: 0, H: h(start, goal)}
	startNode.F = startNode.G + startNode.H
	
	directions := directions4
	if diagonalAllowed {
		directions = directions8
	}
	
	pq := &PriorityQueue{startNode}
	heap.Init(pq)
	
	vis := make(map[Point]bool)
	
	nodesExp := 0
	
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*Node)
		vis[curr.Pos] = true
		nodesExp++
		
		if curr.Pos == goal {
			return AstarAns{
				Path: reconstructPath(curr),
				NodeExp: nodesExp,
				Found: true,
			}
		}
		
		for _, dir := range directions {
			next := Point{curr.Pos.X + dir.X, curr.Pos.Y + dir.Y}
			if !inBounds(next) || !canGo(next) || vis[next]{ continue}
			
			cost := 1.0
			
			newNode := &Node{
				Pos: next,
				G: curr.G + cost,
				H: h(next, goal),
				Parent: curr,
			}
			
			newNode.F = newNode.G + newNode.H
			heap.Push(pq, newNode)
		}
	}
	
	return AstarAns{
		Found: false,
		NodeExp: nodesExp,
	}
}


func reconstructPath(curr *Node) []Point{
	var path []Point

	for curr != nil {	
		path = append(path, curr.Pos)
		curr = curr.Parent
	}
	
	slices.Reverse(path)
	return path
}