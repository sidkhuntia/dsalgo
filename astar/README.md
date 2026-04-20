# A* Search Algorithm — A Complete Guide

> *"A* is Dijkstra with a compass. Dijkstra wanders. A* aims."*

---

## Table of Contents

1. [The Problem A* Solves](#1-the-problem-a-solves)
2. [The Core Idea](#2-the-core-idea)
3. [The Formula: f(n) = g(n) + h(n)](#3-the-formula-fn--gn--hn)
4. [How A* Works — Step by Step](#4-how-a-works--step-by-step)
5. [The Open & Closed Sets](#5-the-open--closed-sets)
6. [Heuristics — The Heart of A*](#6-heuristics--the-heart-of-a)
7. [Heuristic Functions in Detail](#7-heuristic-functions-in-detail)
8. [Choosing a Heuristic](#8-choosing-a-heuristic)
9. [A* on Non-Grid Graphs](#9-a-on-non-grid-graphs)
10. [Heuristics on Random Graphs](#10-heuristics-on-random-graphs)
11. [A* vs Dijkstra](#11-a-vs-dijkstra)
12. [Complexity](#12-complexity)
13. [Pseudocode](#13-pseudocode)

---

## 1. The Problem A* Solves

Given a **weighted graph** (or grid), find the **shortest path** from a start node `S` to a goal node `G`.

```
S ──3── A ──2── B
|               |
5               1
|               |
C ──4── D ──1── G
```

We want the cheapest route from `S` to `G`. There are many possible algorithms — BFS, DFS, Dijkstra, Greedy Best-First. A* is special because it combines the **guarantees of Dijkstra** with the **speed of Greedy search**.

---

## 2. The Core Idea

Every algorithm explores nodes in some order. The question is: **which node do I visit next?**

| Algorithm       | Visits next node by...             | Optimal? | Fast?  |
|-----------------|------------------------------------|----------|--------|
| BFS             | Fewest hops (ignores weights)      | ✗        | Medium |
| Dijkstra        | Cheapest cost from start           | ✓        | Slow   |
| Greedy Best-First | Closest estimated distance to goal | ✗        | Fast   |
| **A\***         | **Cheapest + closest estimate**    | **✓**    | **Fast** |

A* says: **don't just pick the cheapest path so far — also consider how far you still have to go.**

---

## 3. The Formula: f(n) = g(n) + h(n)

This is the entire soul of A*. Every node gets a score:

```
f(n)  =  g(n)  +  h(n)
  │         │        │
total    cost     estimated
score    so far   cost to goal
```

- **`g(n)`** — exact, known cost to reach node `n` from the start
- **`h(n)`** — *heuristic* estimate of cost from `n` to the goal (you design this)
- **`f(n)`** — total estimated cost of the cheapest path through `n`

A* always picks the node with the **lowest `f(n)`** to explore next.

### Visual intuition

```
START ──────────── n ──────────── GOAL
       g(n): known    h(n): guessed
       ◄──────────►  ◄────────────►
            f(n) = g(n) + h(n)
       ◄──────────────────────────►
```

---

## 4. How A* Works — Step by Step

Let's trace through a small example:

```
Grid (S=start, G=goal, █=wall):

S . . . .
. █ █ . .
. █ . . .
. . . █ G
```

### The Loop

```
1. Add START to the open set with f = 0 + h(start)
2. While open set is not empty:
   a. Pick node with lowest f(n)  ← this is the priority queue
   b. If it's GOAL → reconstruct and return path
   c. Move it to closed set
   d. For each neighbour:
      - Skip if in closed set
      - Calculate tentative g = current.g + edge_weight
      - If better than known g → update, set parent, add to open set
3. If open set empties → no path exists
```

### Trace (Manhattan heuristic, goal at bottom-right)

```
Step 1: Pop START (f=0+6=6)
        Discover neighbours, push with f scores

Step 2: Pop best f node
        Update neighbours if cheaper path found

Step 3: ...continue until GOAL is popped

Step N: GOAL popped → walk parent pointers back = PATH
```

### Path Reconstruction

Every node stores a **parent pointer** — when you update a node's `g`, you also update its parent to the current node. When the goal is reached, walk backwards:

```
GOAL → parent → parent → ... → START
then reverse → START → ... → GOAL  ✓
```

---

## 5. The Open & Closed Sets

Two data structures power A*:

### Open Set (Priority Queue)

> "Nodes I've *discovered* but not yet *finalized*"

- Implemented as a **min-heap** on `f(n)`
- Always pop the node with the lowest `f`
- When a cheaper path to an existing node is found → update its `f` and re-sort (`heap.Fix`)
- Also maintained as a **map** (`openMap`) so you can look up and update nodes in O(1)

```
openMap:  node → *Node pointer (for heap.Fix)
heap:     sorted by f(n)
```

Without `openMap`, you'd have to push duplicates into the heap (lazy deletion) — still correct but wastes memory.

### Closed Set (Visited Set)

> "Nodes whose shortest path is already finalized"

- A simple **hash set** of visited nodes
- Once a node is popped from the heap → its `g` is optimal (same proof as Dijkstra)
- Never re-process a closed node

```
closedSet: Set<Node>  →  O(1) lookup
```

### What happens without each?

| Remove        | Consequence                                             |
|---------------|---------------------------------------------------------|
| `closedSet`   | Nodes reprocessed → cycles → **infinite loop**         |
| `openMap`     | Can't do heap.Fix → push duplicates → **wastes memory**|

### Node lifecycle

```
            UNDISCOVERED
                 │
                 │  first time seen as neighbour
                 ▼
             OPEN SET  ◄─── cheaper path found (heap.Fix)
                 │
                 │  popped as lowest f(n)
                 ▼
            CLOSED SET  ──── never touched again
```

---

## 6. Heuristics — The Heart of A*

The heuristic `h(n)` is your **domain knowledge injected into the algorithm**. It answers:

> *"From node `n`, what's the minimum possible cost to reach the goal?"*

### The Admissibility Constraint

For A* to guarantee an optimal path, `h(n)` must be **admissible**:

```
h(n)  ≤  true_cost(n → goal)   for ALL nodes n
```

It must **never overestimate**. It can underestimate (even be 0), but overestimating can skip the optimal path.

### The Consistency Constraint (stronger)

For A* on graphs (not just trees), ideally `h` is also **consistent** (monotone):

```
h(n)  ≤  cost(n → n')  +  h(n')   for every edge n → n'
```

This is the triangle inequality. Consistent heuristics are always admissible. Most natural heuristics (Euclidean, Manhattan) are consistent.

### Heuristic quality affects speed

```
h(n) = 0          →  Dijkstra  (explores everything, guaranteed optimal)
h(n) = perfect    →  explores only nodes on the optimal path (fastest possible)
h(n) = overestimate → may find suboptimal path (not admissible)
```

The better your heuristic estimates the true cost, the fewer nodes A* explores.

---

## 7. Heuristic Functions in Detail

### 7.1 Manhattan Distance

```
h(n) = |n.x - goal.x| + |n.y - goal.y|
```

```
A=(1,1)  →  G=(4,4)

  . . . G
  . . . .
  . . . .
  A . . .

Manhattan = |4-1| + |4-1| = 3 + 3 = 6
```

- ✅ **Best for 4-directional grids** (up/down/left/right only)
- ✅ Admissible and consistent
- ✅ Very fast to compute
- ❌ Overestimates on 8-directional grids (diagonal moves are cheaper than 2 cardinal moves)

---

### 7.2 Euclidean Distance

```
h(n) = √( (n.x - goal.x)² + (n.y - goal.y)² )
```

```
A=(1,1)  →  G=(4,4)

Euclidean = √((4-1)² + (4-1)²) = √18 ≈ 4.24
```

- ✅ **Natural straight-line distance**
- ✅ Admissible for both 4-dir and 8-dir (straight line is always shortest)
- ✅ Works for real-world coordinate graphs (maps, routing)
- ❌ Slightly underestimates on grids (weaker than Manhattan for 4-dir)
- ❌ Involves `sqrt` — slightly slower to compute

---

### 7.3 Chebyshev Distance

```
h(n) = max( |n.x - goal.x|, |n.y - goal.y| )
```

```
A=(1,1)  →  G=(4,4)

Chebyshev = max(3, 3) = 3
```

- ✅ **Best for 8-directional grids** where diagonal cost = cardinal cost = 1
- ✅ Models chess king movement perfectly
- ✅ Admissible and consistent
- ❌ Assumes diagonals cost 1 (not √2) — wrong if diagonal cost is √2

---

### 7.4 Octile Distance

```
D  = 1       (cardinal move cost)
D2 = √2      (diagonal move cost)

dx = |n.x - goal.x|
dy = |n.y - goal.y|

h(n) = D*(dx+dy) + (D2-2D)*min(dx,dy)
     = dx + dy - (2-√2)*min(dx,dy)
```

```
A=(1,1)  →  G=(4,4)

dx=3, dy=3
Octile = 3+3 - (2-√2)*3 = 6 - 0.515*3 ≈ 4.24
```

- ✅ **Most accurate for 8-directional grids with √2 diagonal cost**
- ✅ Tighter than Chebyshev when diagonals cost √2
- ✅ Admissible and consistent
- ✅ No `sqrt` in final formula (uses precomputed `√2`)
- ✅ **Recommended default for game pathfinding**

---

### 7.5 Zero Heuristic (Dijkstra)

```
h(n) = 0
```

- ✅ **Always admissible** (0 never overestimates)
- ✅ Finds optimal path on any graph
- ❌ Explores all nodes radiating from start — **no directional bias**
- ✅ Use when you have no domain knowledge about the graph

---

### 7.6 Landmark / ALT Heuristic

For general graphs with no coordinates, precompute distances from `k` landmark nodes:

```
h(n) = max over all landmarks L of:
         | dist(n, L) - dist(goal, L) |
```

This uses the **triangle inequality**: the distance from `n` to `goal` is at least
`|dist(n,L) - dist(goal,L)|` for any landmark `L`.

```
Precomputation:  k × Dijkstra runs  (offline, once)
Query time:      O(k) per node lookup  (online, fast)
```

- ✅ Works on arbitrary graphs with **no coordinates**
- ✅ Tight, admissible heuristic
- ✅ Used in real GPS routing (OSRM, RoutingKit)
- ❌ Requires offline preprocessing
- ❌ Memory: O(k × |V|) for landmark distances

---

### Heuristic Comparison Summary

```
Movement     Heuristic    Formula                         Quality
─────────────────────────────────────────────────────────────────
4-dir grid   Manhattan    |dx| + |dy|                     ★★★★★
4-dir grid   Euclidean    √(dx²+dy²)                      ★★★★☆
8-dir grid   Octile       dx+dy-(2-√2)*min(dx,dy)         ★★★★★
8-dir grid   Chebyshev    max(|dx|,|dy|)                  ★★★★☆  (cost=1)
Any graph    Euclidean    √(dx²+dy²)  (needs coords)      ★★★☆☆
Any graph    Landmark     |dist(n,L)-dist(goal,L)|         ★★★★☆
Any graph    Zero         0                               ★☆☆☆☆  (= Dijkstra)
```

---

## 8. Choosing a Heuristic

The heuristic comes from **domain knowledge**, not graph structure. Ask yourself:

```
Does your graph have meaning?
│
├── SPATIAL / GEOGRAPHIC
│     Nodes = locations with coordinates
│     → Euclidean or Haversine (for lat/lng)
│     Examples: Maps, game worlds, robot navigation
│
├── STATE-SPACE
│     Nodes = abstract states (puzzle configs, game states)
│     → Domain-specific count
│     Examples: 8-puzzle → misplaced tiles
│               15-puzzle → sum of Manhattan distances of each tile
│               Rubik's cube → pattern databases
│
├── TIME / COST BASED
│     Edges = time or money, not distance
│     → Minimum possible time/cost remaining
│     Examples: Flight routing → minimum flight time remaining
│
├── GENERATED GRAPH (you control it)
│     → Assign (x,y) coordinates at generation time
│     → Enforce edge_weight ≥ Euclidean(a,b)
│     → Now Euclidean is admissible
│
└── PURELY ABSTRACT (no meaning)
      → h = 0  →  Use Dijkstra
      or precompute Landmark heuristic (ALT)
```

### The key question

> *"What do I know about this problem that tells me roughly how far I still am from the goal?"*

If you can answer that, you have a heuristic.

---

## 9. A* on Non-Grid Graphs

A* is **not a grid algorithm**. It works on any weighted graph. The grid is just the most common teaching example because:
- The heuristic is obvious (Euclidean/Manhattan)
- It's easy to visualize

### What changes on a general graph

```go
// Grid: neighbours are the 4/8 adjacent cells
// Graph: neighbours are whatever the adjacency list says

type Graph map[NodeID][]Edge

type Edge struct {
    To     NodeID
    Weight float64
}
```

The A* loop is **identical**. Only the neighbour iteration changes:

```go
// Grid:
for _, dir := range directions {
    neighbour := Point{current.X + dir.X, current.Y + dir.Y}
    cost := 1.0  // uniform
}

// General graph:
for _, edge := range graph[current] {
    neighbour := edge.To
    cost := edge.Weight  // non-uniform
}
```

### Real-world example: Road network

```
Nodes:   Road intersections  (each has lat/lng)
Edges:   Road segments       (each has distance/time as weight)
Heuristic: Euclidean distance between GPS coordinates

A* finds the shortest path from your location to destination,
exploring far fewer intersections than Dijkstra would.
```

---

## 10. Heuristics on Random Graphs

This is the hard case. A "random graph" with no semantic meaning gives you nothing to estimate.

### Option 1 — Assign coordinates (if you generate the graph)

```go
type Node struct {
    ID  int
    X, Y float64  // assign random coordinates
}

// When generating edges, ensure weight >= euclidean distance
// so Euclidean heuristic remains admissible
weight := euclidean(a, b) * (1 + rand.Float64())
```

### Option 2 — Landmark heuristic (if graph is fixed, offline preprocessing ok)

```go
// Pick k random landmark nodes
landmarks := pickLandmarks(graph, k=5)

// Precompute distances from each landmark to all nodes (k Dijkstra runs)
landmarkDist := map[NodeID]map[NodeID]float64{}
for _, l := range landmarks {
    landmarkDist[l] = dijkstra(graph, l)
}

// At query time:
h := func(n, goal NodeID) float64 {
    best := 0.0
    for _, l := range landmarks {
        estimate := math.Abs(landmarkDist[l][n] - landmarkDist[l][goal])
        if estimate > best { best = estimate }
    }
    return best
}
```

### Option 3 — Just use Dijkstra

If no domain knowledge exists and preprocessing isn't feasible:

```go
h := func(_, _ NodeID) float64 { return 0 }
// A* with h=0 is exactly Dijkstra
```

### Decision table

| Scenario                              | Heuristic               | A* advantage |
|---------------------------------------|-------------------------|--------------|
| Grid / spatial graph with coords      | Euclidean / Manhattan   | Large        |
| State-space (puzzle, game)            | Domain-specific count   | Large        |
| Geographic (maps, roads)              | Euclidean / Haversine   | Large        |
| Generated graph (you control)         | Assign coords → Euclidean | Large      |
| Fixed graph, preprocessing ok         | Landmark (ALT)          | Medium       |
| Truly random, no preprocessing        | Zero (= Dijkstra)       | None         |

---

## 11. A* vs Dijkstra

```
Dijkstra:   priority = g(n)
A*:         priority = g(n) + h(n)
```

When `h(n) = 0`: A* **is** Dijkstra. This is the entire relationship.

```
Nodes explored on a 100×100 grid, S=top-left, G=bottom-right:

Dijkstra:  ~5,000 nodes  (expands in all directions equally)
A* (Manhattan): ~200 nodes  (focuses toward goal)

     Dijkstra                    A*
  ████████████                 ....
  ████████████                ......
  ████████████              ....░░░░
  ████████████             ...░░░░░░
  ████████████            ...░░░░░░░
  ███████████G           ...░░░░░░░G

  ████ = explored        .... = explored (much fewer)
```

### When A* is NOT better than Dijkstra

- Graph has no domain knowledge (h=0 anyway)
- You need **all shortest paths** from a source, not just to one goal
- Heuristic is very weak (barely better than 0)

---

## 12. Complexity

| | Time | Space |
|---|---|---|
| Worst case | O((V + E) log V) | O(V) |
| With perfect h | O(path length) | O(path length) |
| With h=0 (Dijkstra) | O((V + E) log V) | O(V) |

A*'s practical speed depends entirely on heuristic quality. A perfect heuristic makes it linear in path length. A zero heuristic makes it identical to Dijkstra.

---

## 13. Pseudocode

```
function AStar(graph, start, goal, h):
    openSet  ← PriorityQueue()
    openMap  ← HashMap()          // node → node pointer (for updates)
    closedSet ← HashSet()

    startNode = {pos: start, g: 0, h: h(start, goal), f: h(start, goal)}
    openSet.push(startNode)
    openMap[start] = startNode

    while openSet is not empty:
        current ← openSet.pop()          // node with lowest f
        openMap.remove(current.pos)
        closedSet.add(current.pos)

        if current.pos == goal:
            return reconstructPath(current)

        for each neighbour of current:
            if neighbour in closedSet: continue

            tentativeG = current.g + edgeCost(current, neighbour)

            if neighbour in openMap:
                existing = openMap[neighbour]
                if tentativeG < existing.g:
                    existing.g = tentativeG
                    existing.f = tentativeG + existing.h
                    existing.parent = current
                    openSet.update(existing)          // heap.Fix
            else:
                n = {pos: neighbour,
                     g: tentativeG,
                     h: h(neighbour, goal),
                     f: tentativeG + h(neighbour, goal),
                     parent: current}
                openSet.push(n)
                openMap[neighbour] = n

    return NO_PATH_FOUND


function reconstructPath(node):
    path ← []
    while node ≠ null:
        path.prepend(node.pos)
        node ← node.parent
    return path
```

---

## Summary

```
┌─────────────────────────────────────────────────────────────────┐
│                        A* in one picture                        │
│                                                                 │
│   f(n) = g(n)  +  h(n)                                         │
│           │          │                                          │
│       "what it    "what it                                      │
│        cost to    might cost                                    │
│        get here"  from here"                                    │
│                                                                 │
│   h = 0          →  Dijkstra  (optimal, slow)                  │
│   h = admissible →  A*        (optimal, fast)                   │
│   h = perfect    →  A*        (optimal, fastest possible)       │
│   h > true cost  →  Greedy    (fast, not optimal)               │
│                                                                 │
│   Works on:  grids, road networks, puzzles, any weighted graph  │
│   Key design: choose h based on domain knowledge                │
└─────────────────────────────────────────────────────────────────┘
```