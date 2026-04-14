package bidirectionldijsktra_test

import (
	"fmt"
	"testing"

	bidirectionldijsktra "dsalgo/bidirectionl-dijsktra"
)

func TestNormalDijsktra(t *testing.T) {
	fmt.Print("Testing Normal")
	adj := make([][]bidirectionldijsktra.Pair, 0)
	n := 5
	start := 0
	adj = append(adj, []bidirectionldijsktra.Pair{{1, 4}, {2, 8}})
	adj = append(adj, []bidirectionldijsktra.Pair{{0, 4}, {4, 6}, {2, 3}})
	adj = append(adj, []bidirectionldijsktra.Pair{{0, 8}, {3, 2}, {1, 3}})
	adj = append(adj, []bidirectionldijsktra.Pair{{2, 2}, {4, 10}})
	adj = append(adj, []bidirectionldijsktra.Pair{{1, 6}, {3, 10}})
	dist := bidirectionldijsktra.NormalDijsktra(n, start, adj)
	actualDist := []int{0, 4, 7, 9, 10}
	for i, d := range dist {
		if d != actualDist[i] {
			t.Errorf("Expected %d, got %d", actualDist[i], d)
		}
	}

}

func TestBidirectionalDijsktra(t *testing.T) {
	// Base 5-node undirected graph used across multiple sub-tests:
	//   0 --4-- 1 --6-- 4
	//   |       |       |
	//   8       3      10
	//   |       |       |
	//   2 --2-- 3 ------+
	makeBaseAdj := func() [][]bidirectionldijsktra.Pair {
		adj := make([][]bidirectionldijsktra.Pair, 5)
		adj[0] = []bidirectionldijsktra.Pair{{To: 1, Weight: 4}, {To: 2, Weight: 8}}
		adj[1] = []bidirectionldijsktra.Pair{{To: 0, Weight: 4}, {To: 4, Weight: 6}, {To: 2, Weight: 3}}
		adj[2] = []bidirectionldijsktra.Pair{{To: 0, Weight: 8}, {To: 3, Weight: 2}, {To: 1, Weight: 3}}
		adj[3] = []bidirectionldijsktra.Pair{{To: 2, Weight: 2}, {To: 4, Weight: 10}}
		adj[4] = []bidirectionldijsktra.Pair{{To: 1, Weight: 6}, {To: 3, Weight: 10}}
		return adj
	}

	tests := []struct {
		name     string
		n        int
		start    int
		end      int
		adj      [][]bidirectionldijsktra.Pair
		expected int
	}{
		{
			name:     "5-node graph: 0 to 4 via 0->1->4",
			n:        5,
			start:    0,
			end:      4,
			adj:      makeBaseAdj(),
			expected: 10,
		},
		{
			name:     "5-node graph: 0 to 3 via 0->1->2->3",
			n:        5,
			start:    0,
			end:      3,
			adj:      makeBaseAdj(),
			expected: 9,
		},
		{
			name:     "5-node graph: 0 to 2 via 0->1->2",
			n:        5,
			start:    0,
			end:      2,
			adj:      makeBaseAdj(),
			expected: 7,
		},
		{
			name:     "5-node graph: 0 to 1 direct edge",
			n:        5,
			start:    0,
			end:      1,
			adj:      makeBaseAdj(),
			expected: 4,
		},
		{
			name:  "single direct edge between start and end",
			n:     2,
			start: 0,
			end:   1,
			adj: [][]bidirectionldijsktra.Pair{
				{{To: 1, Weight: 5}},
				{{To: 0, Weight: 5}},
			},
			expected: 5,
		},
		{
			name:  "linear chain 0->1->2->3->4 all weight 1",
			n:     5,
			start: 0,
			end:   4,
			adj: [][]bidirectionldijsktra.Pair{
				{{To: 1, Weight: 1}},
				{{To: 0, Weight: 1}, {To: 2, Weight: 1}},
				{{To: 1, Weight: 1}, {To: 3, Weight: 1}},
				{{To: 2, Weight: 1}, {To: 4, Weight: 1}},
				{{To: 3, Weight: 1}},
			},
			expected: 4,
		},
		{
			name:  "two paths: indirect shorter than direct",
			n:     3,
			start: 0,
			end:   1,
			// direct 0->1 costs 10, indirect 0->2->1 costs 2+2=4
			adj: [][]bidirectionldijsktra.Pair{
				{{To: 1, Weight: 10}, {To: 2, Weight: 2}},
				{{To: 0, Weight: 10}, {To: 2, Weight: 2}},
				{{To: 0, Weight: 2}, {To: 1, Weight: 2}},
			},
			expected: 4,
		},
		{
			name:  "triangle: shortest goes through intermediate node",
			n:     3,
			start: 0,
			end:   2,
			// 0->2 direct costs 10, 0->1->2 costs 3+4=7
			adj: [][]bidirectionldijsktra.Pair{
				{{To: 1, Weight: 3}, {To: 2, Weight: 10}},
				{{To: 0, Weight: 3}, {To: 2, Weight: 4}},
				{{To: 1, Weight: 4}, {To: 0, Weight: 10}},
			},
			expected: 7,
		},
		{
			name:  "reverse direction: end closer to start in linear chain",
			n:     5,
			start: 4,
			end:   0,
			adj: [][]bidirectionldijsktra.Pair{
				{{To: 1, Weight: 1}},
				{{To: 0, Weight: 1}, {To: 2, Weight: 1}},
				{{To: 1, Weight: 1}, {To: 3, Weight: 1}},
				{{To: 2, Weight: 1}, {To: 4, Weight: 1}},
				{{To: 3, Weight: 1}},
			},
			expected: 4,
		},
		{
			name:  "graph with varied weights: optimal path skips heavy edges",
			n:     4,
			start: 0,
			end:   3,
			// 0->3 direct costs 100, 0->1->2->3 costs 1+1+1=3
			adj: [][]bidirectionldijsktra.Pair{
				{{To: 1, Weight: 1}, {To: 3, Weight: 100}},
				{{To: 0, Weight: 1}, {To: 2, Weight: 1}},
				{{To: 1, Weight: 1}, {To: 3, Weight: 1}},
				{{To: 2, Weight: 1}, {To: 0, Weight: 100}},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bidirectionldijsktra.BidirectionalDijsktra(tt.n, tt.start, tt.end, tt.adj)
			if got != tt.expected {
				t.Errorf("BidirectionalDijsktra(n=%d, start=%d, end=%d) = %d, want %d",
					tt.n, tt.start, tt.end, got, tt.expected)
			}
		})
	}
}
