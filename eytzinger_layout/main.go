package eytzingerlayout

import (
	"fmt"
)

// Eytzinger Layout is a cache-optimized binary search tree layout that stores
// nodes in a specific order to improve memory access patterns.
//
// Advantages over standard binary search:
//  1. Better cache locality: Nodes accessed during search are stored closer together
//     in memory, reducing cache misses and improving performance by 15-30%
//  2. Predictable memory access: The layout ensures sequential memory access patterns
//     which are more CPU-friendly than random access in standard binary search
//  3. Branch prediction: The traversal pattern is more predictable, allowing CPUs
//     to better optimize branch prediction
//  4. SIMD-friendly: The layout can be optimized for vectorized operations in some cases
//
// The trade-off is that insertion/deletion operations are more expensive due to
// the need to maintain the specific layout, making it ideal for read-heavy workloads.
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
