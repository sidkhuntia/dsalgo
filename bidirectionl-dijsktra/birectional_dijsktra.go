package bidirectionldijsktra

import (
	"container/heap"
	"math"
)

type Pair struct {
	To     int
	Weight int
}

type Node struct {
	Weight int
	Index  int
}

type IntHeap []*Node

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i].Weight < h[j].Weight
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(*Node))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func NormalDijsktra(n int, start int, adj [][]Pair) []int {
	dist := make([]int, n)
	for i := range n {
		dist[i] = math.MaxInt
	}

	dist[start] = 0

	minHeap := &IntHeap{}
	heap.Init(minHeap)
	heap.Push(minHeap, &Node{Weight: dist[start], Index: start})
	vis := make([]bool, n)

	for minHeap.Len() > 0 {
		node := heap.Pop(minHeap).(*Node)
		vis[node.Index] = true
		for _, p := range adj[node.Index] {
			if !vis[p.To] && dist[node.Index]+p.Weight < dist[p.To] {
				dist[p.To] = dist[node.Index] + p.Weight
				heap.Push(minHeap, &Node{Weight: dist[p.To], Index: p.To})
			}
		}
	}

	return dist
}

func BidirectionalDijsktra(n int, start int, end int, adj [][]Pair) int {
	distFromStart := make([]int, n)
	distFromEnd := make([]int, n)
	visFromStart := make([]bool, n)
	visFromEnd := make([]bool, n)
	for i := range n {
		distFromStart[i] = math.MaxInt
		distFromEnd[i] = math.MaxInt
	}

	distFromStart[start] = 0
	distFromEnd[end] = 0

	minHeapFromStart := &IntHeap{}
	heap.Init(minHeapFromStart)
	heap.Push(minHeapFromStart, &Node{Weight: distFromStart[start], Index: start})
	minHeapFromEnd := &IntHeap{}
	heap.Init(minHeapFromEnd)
	heap.Push(minHeapFromEnd, &Node{Weight: distFromEnd[end], Index: end})
	
	ans := math.MaxInt

	for minHeapFromStart.Len() > 0 && minHeapFromEnd.Len() > 0 {
		nodeS := heap.Pop(minHeapFromStart).(*Node)
		nodeE := heap.Pop(minHeapFromEnd).(*Node)
		visFromStart[nodeS.Index] = true
		visFromEnd[nodeE.Index] = true
		for _, p := range adj[nodeS.Index] {
			if !visFromStart[p.To] && distFromStart[nodeS.Index]+p.Weight < distFromStart[p.To] {
				distFromStart[p.To] = distFromStart[nodeS.Index] + p.Weight
				heap.Push(minHeapFromStart, &Node{Weight: distFromStart[p.To], Index: p.To})
			}
			
			if visFromEnd[p.To] && distFromStart[nodeS.Index] + p.Weight + distFromEnd[p.To] < ans {
				ans = distFromStart[nodeS.Index] + p.Weight + distFromEnd[p.To]
			}
		}
		for _, p := range adj[nodeE.Index] {
			if !visFromEnd[p.To] && distFromEnd[nodeE.Index]+p.Weight < distFromEnd[p.To] {
				distFromEnd[p.To] = distFromEnd[nodeE.Index] + p.Weight
				heap.Push(minHeapFromEnd, &Node{Weight: distFromEnd[p.To], Index: p.To})
			}
			
			if visFromStart[p.To] && distFromEnd[nodeE.Index] + p.Weight + distFromStart[p.To] < ans {
				ans = distFromEnd[nodeE.Index] + p.Weight + distFromStart[p.To]
			}
		}
		
		if distFromEnd[nodeE.Index] + distFromStart[nodeS.Index] >= ans {
			break
		}
		
	}
	
	return ans
}