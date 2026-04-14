# Bidirectional Dijkstra

## Overview

Bidirectional Dijkstra is an optimized shortest-path algorithm that simultaneously runs two Dijkstra searches — one **forward** from the source `s` and one **backward** from the target `t` — stopping when the two frontiers meet. This approach can dramatically reduce the number of edges examined compared to a single-direction Dijkstra search.

## Motivation

Standard Dijkstra with branching factor `b` and path length `d` explores on the order of `O(b^d)` nodes. By splitting the search in two, each frontier only needs to expand `d/2` hops, giving `O(2 * b^(d/2))` total work — exponentially less for large graphs.

## Algorithm

Two priority queues (`Q_f` and `Q_b`) and two visited sets (`S_f` and `S_b`) are maintained. A running best estimate `μ` tracks the shortest `s → t` path seen so far.

```
μ = ∞

while Q_f and Q_b are not empty:
    u = extract_min(Q_f)     // forward step
    v = extract_min(Q_b)     // backward step

    S_f.add(u), S_b.add(v)

    for each neighbor x of u:
        relax(u → x) in forward direction
        if x ∈ S_b and d_f[u] + w(u,x) + d_b[x] < μ:
            μ = d_f[u] + w(u,x) + d_b[x]

    for each neighbor x of v:
        relax(v → x) in backward direction
        if x ∈ S_f and d_b[v] + w(v,x) + d_f[x] < μ:
            μ = d_b[v] + w(v,x) + d_f[x]

    if d_f[u] + d_b[v] >= μ:
        break   // μ is the true shortest path distance
```

The key subtlety is the **termination condition**. Simply stopping when any node is visited by both searches is *incorrect* — that edge might not lie on the actual shortest path. The algorithm must continue until the sum of the two frontiers' minimum distances meets or exceeds `μ`.

> For directed graphs, the backward search must be run on the **reverse graph** (transpose), following edges in the opposite direction.

## Complexity

| | Time | Space |
|---|---|---|
| Standard Dijkstra | `O((V + E) log V)` | `O(V)` |
| Bidirectional Dijkstra | `O((V + E) log V)` worst case, but **~2× faster** in practice | `O(V)` |

The theoretical worst-case complexity class is the same, but the constant factor improvement from halving the search radius is significant on large real-world graphs (e.g. road networks, social graphs).

## Implementation

This package provides two functions:

- **`NormalDijsktra(n, start, adj)`** — standard single-source Dijkstra returning distances from `start` to all nodes.
- **`BidirectionalDijsktra(n, start, end, adj)`** — bidirectional Dijkstra returning the shortest distance from `start` to `end`.

The graph is represented as an adjacency list of `Pair{To, Weight}` edges using a min-heap backed priority queue.

### Example Usage

```go
adj := make([][]bidirectionldijsktra.Pair, 5)
adj[0] = append(adj[0], bidirectionldijsktra.Pair{To: 1, Weight: 4})
adj[1] = append(adj[1], bidirectionldijsktra.Pair{To: 2, Weight: 3})
// ... build your graph

dist := bidirectionldijsktra.BidirectionalDijsktra(5, 0, 4, adj)
fmt.Println(dist) // shortest distance from node 0 to node 4
```

## Recent Research

A 2024 paper by Haeupler et al. proves that a careful implementation of bidirectional Dijkstra is **instance-optimal** — no correct algorithm can outperform it on any single instance by more than a constant factor. This provides strong theoretical backing for what has long been observed in practice.

## References

- [Bidirectional Search — Wikipedia](https://en.wikipedia.org/wiki/Bidirectional_search?useskin=vector)
- [Bidirectional Dijkstra — Matthew Towers, UCL](https://www.homepages.ucl.ac.uk/~ucahmto/math/2020/05/30/bidirectional-dijkstra.html)
- [Bidirectional Dijkstra's Algorithm is Instance-Optimal — Haeupler et al. (arXiv 2024)](https://arxiv.org/abs/2410.14638)