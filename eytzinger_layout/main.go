package eytzingerlayout

import (
	"fmt"
)

// Eytzinger Layout: A Cache-Optimized Array Layout for Binary Search
//
// LECTURE NOTES:
// ==============
//
//  1. DEFINITION
//     The Eytzinger layout stores n data items in an array by viewing them as nodes
//     in a complete binary search tree and placing values in breadth-first order (BFS).
//     This is the same layout used for heaps, but applied to binary search trees.
//
// 2. WHY EYTZINGER IS FASTER THAN STANDARD BINARY SEARCH
//
//	a) Cache Locality & Prefetching:
//	   - Standard binary search has poor cache behavior: each comparison accesses
//	     memory locations that are far apart (often in different cache lines)
//	   - Eytzinger layout stores nodes accessed during search closer together
//	   - Modern CPU prefetchers can predict and prefetch the next nodes in the
//	     traversal path, effectively hiding memory latency
//	   - Sequential memory access patterns are 2-3x faster than random access
//
//	b) Branch Prediction:
//	   - The traversal pattern in Eytzinger search is more predictable
//	   - CPUs can better optimize branch prediction, reducing pipeline stalls
//	   - Combined with conditional moves (CMOV), branch mispredictions are minimized
//
//	c) Memory Access Patterns:
//	   - Binary search jumps around: indices like n/2, n/4, 3n/4, etc. create
//	     unpredictable memory access patterns
//	   - Eytzinger search follows a more sequential pattern: parent → left child → right child
//	   - This sequential pattern is cache-friendly and prefetcher-friendly
//
// 3. PERFORMANCE CHARACTERISTICS
//
//	According to experimental research (Khuong & Morin, 2017):
//	- For small arrays (n < 2^16): Sorted order with binary search is fastest
//	- For large arrays (n >= 2^16): Eytzinger layout is usually the fastest
//	- Performance improvement: 15-30% faster than standard binary search for large n
//	- The advantage increases with array size due to better cache utilization
//
// 4. TRADE-OFFS
//
//	Advantages:
//	- Faster search queries (especially for large, read-heavy workloads)
//	- Better cache utilization
//	- Predictable memory access patterns
//
//	Disadvantages:
//	- Static data structure: insertion/deletion requires rebuilding the layout
//	- Construction overhead: O(n) time to build the layout from sorted input
//	- Not suitable for dynamic datasets with frequent updates
//
// 5. USE CASES
//
//	Ideal for:
//	- Static or rarely-changing datasets
//	- Read-heavy workloads with many queries
//	- Large datasets where cache performance matters
//	- Real-time systems (e.g., ad bidding engines, browser code)
//
//	Not suitable for:
//	- Dynamic datasets requiring frequent insertions/deletions
//	- Small arrays where overhead doesn't justify the benefit
//
// 6. IMPLEMENTATION NOTES
//
//   - The layout uses 0-based indexing (root at index 0)
//   - For node at index k: left child = 2*k+1, right child = 2*k+2
//   - Search algorithm simulates traversal of the implicit binary search tree
//   - Can be optimized further with branch-free code and explicit prefetching
//
// REFERENCE:
//
//	Khuong, P. V., & Morin, P. (2017). Array Layouts for Comparison-Based Searching.
//	arXiv preprint arXiv:1509.05053. https://arxiv.org/pdf/1509.05053
//
//	Key findings:
//	- After extensive testing on modern hardware, Eytzinger layout is fastest for large n
//	- This conclusion was surprising and contradicted earlier work
//	- The performance advantage comes from better cache behavior and prefetching
//	- Branch-free implementations with conditional moves perform best
//
// . CPP IMPLEMENTATION:
//
//	https://algorithmica.org/en/eytzinger

type EytzingerLayout struct {
	nodes []int
}

func NewEytzingerLayout(input []int) *EytzingerLayout {
	n := len(input)
	eytzinger := make([]int, n)
	it := 0

	var buildEytzinger func(int)
	buildEytzinger = func(k int) {
		if k < n {
			buildEytzinger(2*k + 1)
			eytzinger[k] = input[it]
			it++
			buildEytzinger(2*k + 2)
		}
	}
	buildEytzinger(0)
	return &EytzingerLayout{
		nodes: eytzinger,
	}
}

func (e *EytzingerLayout) Print() {
	for i := 0; i < len(e.nodes); i++ {
		fmt.Printf("%d ", e.nodes[i])
	}
	fmt.Println()
}

// Search can be improved in cpp https://algorithmica.org/en/eytzinger
func (e *EytzingerLayout) Search(x int) int {
	k := 0
	for k < len(e.nodes) {
		if e.nodes[k] < x {
			k = 2*k + 2
		} else if e.nodes[k] > x {
			k = 2*k + 1
		} else {
			return k
		}
	}
	return -1
}
