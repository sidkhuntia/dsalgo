package astar

import (
	"math"
	"testing"
)

// ─────────────────────────────────────────────
// Test grids                 (all grid[X][Y] — column-major)
// ─────────────────────────────────────────────

// smallGrid — 5×5.
//
//	  X: 0  1  2  3  4
//	Y=0: .  .  .  .  .
//	Y=1: .  #  #  #  .
//	Y=2: .  #  ■  #  .   ← (2,2) open but enclosed, unreachable in 4-dir
//	Y=3: .  #  #  #  .
//	Y=4: .  .  .  .  .
var smallGrid = [][]int{
	{0, 0, 0, 0, 0}, // X=0
	{0, 1, 1, 1, 0}, // X=1
	{0, 1, 0, 1, 0}, // X=2
	{0, 1, 1, 1, 0}, // X=3
	{0, 0, 0, 0, 0}, // X=4
}

// bigGrid — 10×10.
//
// A vertical barrier at X=4 has exactly ONE gap at Y=5, so every 4-dir path
// crossing from the left half (X<4) to the right half (X>4) must pass through
// the cell (4,5).
//
// (7,7) is an open cell whose four cardinal neighbours are all walls:
//
//	(6,7), (8,7), (7,6), (7,8) → walls
//
// This makes (7,7) unreachable in 4-dir mode, but reachable via the diagonal
// neighbours (6,6), (8,6), (6,8), (8,8) in 8-dir mode.
//
//	  X: 0  1  2  3  4  5  6  7  8  9
//	Y=0: .  .  .  .  #  .  .  .  .  .
//	Y=1: .  .  .  .  #  .  .  .  .  .
//	Y=2: .  .  .  .  #  .  .  .  .  .
//	Y=3: .  .  .  .  #  .  .  .  .  .
//	Y=4: .  .  .  .  #  .  .  .  .  .
//	Y=5: .  .  .  .  .  .  .  .  .  .   ← only gap
//	Y=6: .  .  .  .  #  .  .  #  .  .
//	Y=7: .  .  .  .  #  .  #  ■  #  .   ← (7,7) enclosed (4-dir)
//	Y=8: .  .  .  .  #  .  .  #  .  .
//	Y=9: .  .  .  .  #  .  .  .  .  .
//
// Known optimal costs (unit step):
//
//	4-dir (0,0)→(9,9) = 18   (must route through the gap)
//	4-dir (0,5)→(9,5) =  9   (straight line through gap)
//	8-dir (0,0)→(9,9) = 10
var bigGrid = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // X=0  all open
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // X=1  all open
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // X=2  all open
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // X=3  all open
	{1, 1, 1, 1, 1, 0, 1, 1, 1, 1}, // X=4  barrier — gap at Y=5
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // X=5  all open
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // X=6  (6,7) wall
	{0, 0, 0, 0, 0, 0, 1, 0, 1, 0}, // X=7  (7,6) wall · (7,7) open+enclosed · (7,8) wall
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // X=8  (8,7) wall
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // X=9  all open
}

// ─────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// dijkstraCost runs a plain BFS (every move costs 1) and returns the minimum
// number of steps from start to goal, or -1 when the goal is unreachable.
// This acts as the ground-truth oracle for path-optimality assertions.
func dijkstraCost(grid [][]int, start, goal Point, diagonal bool) int {
	if start == goal {
		return 0
	}
	rows := len(grid)
	cols := len(grid[0])
	inBounds := func(p Point) bool { return p.X >= 0 && p.X < cols && p.Y >= 0 && p.Y < rows }
	canGo := func(p Point) bool { return grid[p.X][p.Y] == 0 }

	dirs := directions4
	if diagonal {
		dirs = directions8
	}

	type state struct {
		pos  Point
		dist int
	}

	visited := map[Point]bool{start: true}
	queue := []state{{start, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			next := Point{cur.pos.X + d.X, cur.pos.Y + d.Y}
			if !inBounds(next) || !canGo(next) || visited[next] {
				continue
			}
			if next == goal {
				return cur.dist + 1
			}
			visited[next] = true
			queue = append(queue, state{next, cur.dist + 1})
		}
	}
	return -1 // unreachable
}

// validatePath asserts the structural correctness of a path returned by Astar:
//  1. First element == start, last element == goal.
//  2. Every cell is in-bounds and walkable (grid value == 0).
//  3. Every consecutive pair of cells is a legal single move for the given mode.
func validatePath(t *testing.T, grid [][]int, path []Point, start, goal Point, diagonal bool) {
	t.Helper()
	if len(path) == 0 {
		t.Fatal("path is empty")
		return
	}

	rows := len(grid)
	cols := len(grid[0])

	if path[0] != start {
		t.Errorf("path[0] = %v, want start %v", path[0], start)
	}
	if path[len(path)-1] != goal {
		t.Errorf("path[last] = %v, want goal %v", path[len(path)-1], goal)
	}

	for i, p := range path {
		if p.X < 0 || p.X >= cols || p.Y < 0 || p.Y >= rows {
			t.Errorf("path[%d] = %v out of bounds (cols=%d rows=%d)", i, p, cols, rows)
			continue
		}
		if grid[p.X][p.Y] != 0 {
			t.Errorf("path[%d] = %v is a wall (grid[%d][%d] = %d)", i, p, p.X, p.Y, grid[p.X][p.Y])
		}
		if i == 0 {
			continue
		}
		prev := path[i-1]
		dx := iabs(p.X - prev.X)
		dy := iabs(p.Y - prev.Y)
		if diagonal {
			if dx > 1 || dy > 1 || (dx == 0 && dy == 0) {
				t.Errorf("path[%d→%d]: invalid 8-dir step %v→%v (Δx=%d, Δy=%d)", i-1, i, prev, p, dx, dy)
			}
		} else {
			if dx+dy != 1 {
				t.Errorf("path[%d→%d]: invalid 4-dir step %v→%v (Δx=%d, Δy=%d)", i-1, i, prev, p, dx, dy)
			}
		}
	}
}

func requireFound(t *testing.T, ans AstarAns, start, goal Point) {
	t.Helper()
	if !ans.Found {
		t.Fatalf("expected a path from %v to %v but none was found", start, goal)
	}
}

// allHeuristics collects every named HeuristicFn exported by the package.
var allHeuristics = []struct {
	name string
	fn   HeuristicFn
}{
	{"Manhattan", Manhanttan},
	{"Euclidean", Euclidean},
	{"Chebyshev", Chebyshev},
	{"Octile", Octile},
	{"Zero", Zero},
}

// ─────────────────────────────────────────────
// Heuristic unit-property tests
// ─────────────────────────────────────────────

// TestHeuristicZeroAtGoal: h(p, p) == 0 for every heuristic.
func TestHeuristicZeroAtGoal(t *testing.T) {
	points := []Point{{0, 0}, {3, 7}, {9, 9}, {5, 0}, {0, 5}}
	for _, h := range allHeuristics {
		for _, p := range points {
			if v := h.fn(p, p); v != 0 {
				t.Errorf("%s(%v, %v) = %g, want 0", h.name, p, p, v)
			}
		}
	}
}

// TestHeuristicNonNegative: h(a, b) >= 0 for all pairs.
func TestHeuristicNonNegative(t *testing.T) {
	pairs := [][2]Point{
		{{0, 0}, {9, 9}},
		{{9, 9}, {0, 0}},
		{{3, 7}, {8, 2}},
		{{0, 5}, {5, 0}},
		{{4, 4}, {4, 4}},
	}
	for _, h := range allHeuristics {
		for _, pp := range pairs {
			if v := h.fn(pp[0], pp[1]); v < 0 {
				t.Errorf("%s(%v, %v) = %g, want >= 0", h.name, pp[0], pp[1], v)
			}
		}
	}
}

// TestHeuristicSymmetry: h(a, b) == h(b, a).
func TestHeuristicSymmetry(t *testing.T) {
	pairs := [][2]Point{
		{{0, 0}, {9, 9}},
		{{3, 7}, {8, 2}},
		{{0, 5}, {5, 0}},
		{{1, 3}, {4, 1}},
	}
	for _, h := range allHeuristics {
		for _, pp := range pairs {
			fwd := h.fn(pp[0], pp[1])
			bwd := h.fn(pp[1], pp[0])
			if math.Abs(fwd-bwd) > 1e-9 {
				t.Errorf("%s: h(%v,%v)=%g != h(%v,%v)=%g",
					h.name, pp[0], pp[1], fwd, pp[1], pp[0], bwd)
			}
		}
	}
}

// TestHeuristicMonotonicallyIncreasing: as X-distance grows (Y fixed),
// every heuristic must be non-decreasing.
func TestHeuristicMonotonicallyIncreasing(t *testing.T) {
	origin := Point{0, 0}
	for _, h := range allHeuristics {
		prev := h.fn(origin, origin)
		for x := 1; x <= 9; x++ {
			cur := h.fn(origin, Point{x, 0})
			if cur < prev-1e-9 {
				t.Errorf("%s: h({0,0},{%d,0})=%g < h({0,0},{%d,0})=%g; must be non-decreasing",
					h.name, x, cur, x-1, prev)
			}
			prev = cur
		}
	}
}

// TestHeuristicKnownValues spot-checks analytically derived values for every
// heuristic with a=(0,0), b=(3,4).
//
//	Manhattan  = |3|+|4|          = 7
//	Euclidean  = √(9+16)          = 5
//	Chebyshev  = max(3,4)         = 4
//	Octile     = |3−4| + √2·min(3,4)
//	           = 1 + 3√2          ≈ 5.2426  (equivalent: 4 + 3(√2−1))
//	Zero       = 0
func TestHeuristicKnownValues(t *testing.T) {
	a, b := Point{0, 0}, Point{3, 4}

	cases := []struct {
		name string
		fn   HeuristicFn
		want float64
	}{
		{"Manhattan", Manhanttan, 7},
		{"Euclidean", Euclidean, 5},
		{"Chebyshev", Chebyshev, 4},
		{"Octile", Octile, 1 + 3*math.Sqrt2}, // = 4 + 3*(√2−1)
		{"Zero", Zero, 0},
	}

	for _, tc := range cases {
		got := tc.fn(a, b)
		if math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("%s(%v,%v) = %.9f, want %.9f", tc.name, a, b, got, tc.want)
		}
	}
}

// TestHeuristicOrderingAtSamePoint: for the same (a,b), the ranking
// Euclidean ≥ Chebyshev and Manhattan ≥ Euclidean must hold
// (geometric fact independent of obstacles).
func TestHeuristicOrderingAtSamePoint(t *testing.T) {
	pairs := [][2]Point{
		{{0, 0}, {3, 4}},
		{{1, 1}, {7, 8}},
		{{0, 0}, {9, 0}},
		{{2, 5}, {8, 1}},
	}
	for _, pp := range pairs {
		a, b := pp[0], pp[1]
		mnh := Manhanttan(a, b)
		euc := Euclidean(a, b)
		chb := Chebyshev(a, b)
		oct := Octile(a, b)

		if mnh < euc-1e-9 {
			t.Errorf("Manhattan(%v,%v)=%g < Euclidean=%g; impossible", a, b, mnh, euc)
		}
		if euc < chb-1e-9 {
			t.Errorf("Euclidean(%v,%v)=%g < Chebyshev=%g; impossible", a, b, euc, chb)
		}
		if mnh < oct-1e-9 {
			t.Errorf("Manhattan(%v,%v)=%g < Octile=%g; Octile must be ≤ Manhattan", a, b, mnh, oct)
		}
	}
}

// TestHeuristicAdmissibility4Dir: every heuristic must not exceed the true
// 4-dir BFS cost for reachable pairs in bigGrid (admissibility requirement).
func TestHeuristicAdmissibility4Dir(t *testing.T) {
	pairs := []struct{ start, goal Point }{
		{Point{0, 0}, Point{9, 9}},
		{Point{0, 0}, Point{9, 5}},
		{Point{0, 5}, Point{9, 5}},
		{Point{0, 0}, Point{3, 9}},
		{Point{5, 0}, Point{9, 9}},
		{Point{0, 0}, Point{0, 9}},
	}

	for _, pp := range pairs {
		trueCost := float64(dijkstraCost(bigGrid, pp.start, pp.goal, false))
		if trueCost < 0 {
			continue
		}
		for _, h := range allHeuristics {
			val := h.fn(pp.start, pp.goal)
			if val > trueCost+1e-9 {
				t.Errorf("%s(%v→%v) = %.4f > true 4-dir cost %.0f — inadmissible",
					h.name, pp.start, pp.goal, val, trueCost)
			}
		}
	}
}

// TestHeuristicAdmissibility8Dir: for 8-dir unit-cost movement only
// Chebyshev and Zero are provably admissible (Euclidean and Octile can
// exceed the Chebyshev distance which equals the true unit 8-dir cost).
func TestHeuristicAdmissibility8Dir(t *testing.T) {
	pairs := []struct{ start, goal Point }{
		{Point{0, 0}, Point{9, 9}},
		{Point{0, 5}, Point{9, 5}},
		{Point{0, 0}, Point{9, 0}},
		{Point{2, 2}, Point{8, 8}},
	}
	admissible8Dir := []struct {
		name string
		fn   HeuristicFn
	}{
		{"Chebyshev", Chebyshev},
		{"Zero", Zero},
	}

	for _, pp := range pairs {
		trueCost := float64(dijkstraCost(bigGrid, pp.start, pp.goal, true))
		if trueCost < 0 {
			continue
		}
		for _, h := range admissible8Dir {
			val := h.fn(pp.start, pp.goal)
			if val > trueCost+1e-9 {
				t.Errorf("%s(%v→%v) = %.4f > true 8-dir cost %.0f — inadmissible",
					h.name, pp.start, pp.goal, val, trueCost)
			}
		}
	}
}

// TestEuclideanAndOctileCanExceedChebyshev documents (does NOT fail) the
// fact that Euclidean and Octile can overestimate the 8-dir unit cost for
// diagonal steps, making them inadmissible for that movement model.
func TestEuclideanAndOctileCanExceedChebyshev(t *testing.T) {
	// One diagonal step: (0,0)→(1,1)  — true 8-dir cost = 1, Chebyshev = 1.
	a, b := Point{0, 0}, Point{1, 1}
	chb := Chebyshev(a, b) // 1.0
	euc := Euclidean(a, b) // √2 ≈ 1.414
	oct := Octile(a, b)    // √2 ≈ 1.414

	if euc <= chb {
		t.Errorf("expected Euclidean(%v,%v) > Chebyshev; got euc=%.4f chb=%.4f", a, b, euc, chb)
	}
	if oct <= chb {
		t.Errorf("expected Octile(%v,%v) > Chebyshev; got oct=%.4f chb=%.4f", a, b, oct, chb)
	}
	t.Logf("(0,0)→(1,1): Chebyshev=%.4f  Euclidean=%.4f  Octile=%.4f", chb, euc, oct)
}

// ─────────────────────────────────────────────
// Small-grid basic correctness
// ─────────────────────────────────────────────

func TestAstarFound(t *testing.T) {
	start, goal := Point{0, 0}, Point{4, 4}
	ans := Astar(smallGrid, start, goal, Manhanttan, false)
	requireFound(t, ans, start, goal)
	validatePath(t, smallGrid, ans.Path, start, goal, false)
}

func TestAstarPathEndsAtGoal(t *testing.T) {
	start, goal := Point{0, 0}, Point{4, 4}
	ans := Astar(smallGrid, start, goal, Manhanttan, false)
	requireFound(t, ans, start, goal)
	last := ans.Path[len(ans.Path)-1]
	if last != goal {
		t.Errorf("path ends at %v, want %v", last, goal)
	}
}

func TestAstarNoPath(t *testing.T) {
	start, goal := Point{0, 0}, Point{2, 2}
	ans := Astar(smallGrid, start, goal, Manhanttan, false)
	if ans.Found {
		t.Errorf("expected no path from %v to enclosed cell %v", start, goal)
	}
	if ans.NodeExp == 0 {
		t.Error("NodeExp must be > 0 even when no path exists")
	}
}

func TestAstarStartEqualsGoal(t *testing.T) {
	pos := Point{0, 0}
	ans := Astar(smallGrid, pos, pos, Manhanttan, false)
	if !ans.Found {
		t.Fatal("Found must be true when start == goal")
	}
	if ans.NodeExp != 1 {
		t.Errorf("NodeExp = %d, want 1", ans.NodeExp)
	}
	if len(ans.Path) != 1 || ans.Path[0] != pos {
		t.Errorf("Path = %v, want [%v]", ans.Path, pos)
	}
}

func TestAstarAdjacentCells(t *testing.T) {
	start, goal := Point{0, 0}, Point{1, 0}
	ans := Astar(smallGrid, start, goal, Manhanttan, false)
	requireFound(t, ans, start, goal)
	validatePath(t, smallGrid, ans.Path, start, goal, false)
	if got := len(ans.Path) - 1; got != 1 {
		t.Errorf("path cost = %d, want 1 (adjacent cells)", got)
	}
}

// ─────────────────────────────────────────────
// Big grid — 4-directional movement
// ─────────────────────────────────────────────

// TestBigGrid4DirCrossBarrier: (0,0)→(9,9) must route through the gap at (4,5).
// Every admissible heuristic must produce the same optimal cost of 18.
func TestBigGrid4DirCrossBarrier(t *testing.T) {
	start, goal := Point{0, 0}, Point{9, 9}
	const wantCost = 18

	for _, h := range allHeuristics {
		t.Run(h.name, func(t *testing.T) {
			ans := Astar(bigGrid, start, goal, h.fn, false)
			requireFound(t, ans, start, goal)
			validatePath(t, bigGrid, ans.Path, start, goal, false)
			if got := len(ans.Path) - 1; got != wantCost {
				t.Errorf("path cost = %d, want %d", got, wantCost)
			}
		})
	}
}

// TestBigGrid4DirStraightThroughGap: (0,5)→(9,5) at constant Y=5 passes
// directly through the gap — optimal cost is 9 for every heuristic.
func TestBigGrid4DirStraightThroughGap(t *testing.T) {
	start, goal := Point{0, 5}, Point{9, 5}
	const wantCost = 9

	for _, h := range allHeuristics {
		t.Run(h.name, func(t *testing.T) {
			ans := Astar(bigGrid, start, goal, h.fn, false)
			requireFound(t, ans, start, goal)
			validatePath(t, bigGrid, ans.Path, start, goal, false)
			if got := len(ans.Path) - 1; got != wantCost {
				t.Errorf("path cost = %d, want %d", got, wantCost)
			}
		})
	}
}

// TestBigGrid4DirEnclosedCell: (7,7) is surrounded by cardinal walls and must
// be reported as unreachable in 4-dir mode.
func TestBigGrid4DirEnclosedCell(t *testing.T) {
	// Verify the enclosure exists in the grid definition.
	for _, neighbour := range []Point{{6, 7}, {8, 7}, {7, 6}, {7, 8}} {
		if bigGrid[neighbour.X][neighbour.Y] != 1 {
			t.Fatalf("expected grid[%d][%d] = 1 (wall) to form enclosure around (7,7)", neighbour.X, neighbour.Y)
		}
	}
	if bigGrid[7][7] != 0 {
		t.Fatal("expected grid[7][7] = 0 (open) for the enclosed cell test")
	}

	ans := Astar(bigGrid, Point{0, 0}, Point{7, 7}, Manhanttan, false)
	if ans.Found {
		t.Error("expected Found=false for enclosed cell (7,7) in 4-dir mode")
	}
}

// TestBigGrid4DirLeftHalfNoObstacles: paths that remain entirely in the left
// half (X<4) have no obstacles, so cost equals Manhattan distance exactly.
func TestBigGrid4DirLeftHalfNoObstacles(t *testing.T) {
	cases := []struct {
		start, goal Point
		wantCost    int
	}{
		{start: Point{0, 0}, goal: Point{0, 9}, wantCost: 9},  // straight down column 0
		{start: Point{0, 0}, goal: Point{3, 0}, wantCost: 3},  // straight right to X=3
		{start: Point{0, 0}, goal: Point{3, 9}, wantCost: 12}, // |3|+|9|
		{start: Point{0, 9}, goal: Point{3, 0}, wantCost: 12}, // |3|+|9|
		{start: Point{1, 2}, goal: Point{3, 7}, wantCost: 7},  // |2|+|5|
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			ans := Astar(bigGrid, tc.start, tc.goal, Manhanttan, false)
			requireFound(t, ans, tc.start, tc.goal)
			validatePath(t, bigGrid, ans.Path, tc.start, tc.goal, false)
			if got := len(ans.Path) - 1; got != tc.wantCost {
				t.Errorf("%v→%v: cost=%d, want %d", tc.start, tc.goal, got, tc.wantCost)
			}
		})
	}
}

// ─────────────────────────────────────────────
// 4-dir Dijkstra cross-validation
// ─────────────────────────────────────────────

// TestDijkstraCrossValidation4Dir runs A* with every heuristic on a set of
// (start, goal) pairs and asserts that the resulting path cost matches the
// BFS-Dijkstra optimal.  All five heuristics are admissible for 4-dir
// unit-cost grids, so all must produce optimal paths.
func TestDijkstraCrossValidation4Dir(t *testing.T) {
	pairs := []struct{ start, goal Point }{
		{start: Point{0, 0}, goal: Point{9, 9}},
		{start: Point{0, 5}, goal: Point{9, 5}},
		{start: Point{0, 0}, goal: Point{3, 9}},
		{start: Point{5, 0}, goal: Point{9, 9}},
		{start: Point{0, 9}, goal: Point{9, 0}},
		{start: Point{0, 0}, goal: Point{9, 0}},
		{start: Point{2, 2}, goal: Point{2, 8}},
		{start: Point{1, 0}, goal: Point{9, 9}},
		{start: Point{0, 0}, goal: Point{5, 9}},
		{start: Point{3, 3}, goal: Point{9, 9}},
	}

	for _, pp := range pairs {
		optimal := dijkstraCost(bigGrid, pp.start, pp.goal, false)
		if optimal < 0 {
			t.Logf("skip unreachable pair %v→%v", pp.start, pp.goal)
			continue
		}

		for _, h := range allHeuristics {
			t.Run(h.name, func(t *testing.T) {
				ans := Astar(bigGrid, pp.start, pp.goal, h.fn, false)
				if !ans.Found {
					t.Fatalf("[%s] %v→%v: A* Found=false; Dijkstra cost=%d",
						h.name, pp.start, pp.goal, optimal)
				}
				validatePath(t, bigGrid, ans.Path, pp.start, pp.goal, false)
				if got := len(ans.Path) - 1; got != optimal {
					t.Errorf("[%s] %v→%v: cost=%d, Dijkstra optimal=%d",
						h.name, pp.start, pp.goal, got, optimal)
				}
			})
		}
	}
}

// ─────────────────────────────────────────────
// Big grid — 8-directional movement
// ─────────────────────────────────────────────

// TestBigGrid8DirCrossBarrier: (0,0)→(9,9) with diagonal moves.
// The barrier still forces routing through (4,5), giving optimal cost 10.
// Cross-validated against BFS-Dijkstra using the admissible heuristics.
func TestBigGrid8DirCrossBarrier(t *testing.T) {
	start, goal := Point{0, 0}, Point{9, 9}
	const wantCost = 10

	admissible := []struct {
		name string
		fn   HeuristicFn
	}{
		{"Chebyshev", Chebyshev},
		{"Zero", Zero},
	}

	// Confirm the Dijkstra reference agrees with wantCost.
	if got := dijkstraCost(bigGrid, start, goal, true); got != wantCost {
		t.Fatalf("dijkstraCost 8-dir %v→%v = %d, expected %d (grid definition error?)",
			start, goal, got, wantCost)
	}

	for _, h := range admissible {
		t.Run(h.name, func(t *testing.T) {
			ans := Astar(bigGrid, start, goal, h.fn, true)
			requireFound(t, ans, start, goal)
			validatePath(t, bigGrid, ans.Path, start, goal, true)
			if got := len(ans.Path) - 1; got != wantCost {
				t.Errorf("path cost = %d, want %d", got, wantCost)
			}
		})
	}
}

// TestBigGrid8DirAllHeuristicsFind: every heuristic must find *a* valid path
// in 8-dir mode (optimality not required for inadmissible ones).
func TestBigGrid8DirAllHeuristicsFind(t *testing.T) {
	start, goal := Point{0, 0}, Point{9, 9}
	for _, h := range allHeuristics {
		t.Run(h.name, func(t *testing.T) {
			ans := Astar(bigGrid, start, goal, h.fn, true)
			requireFound(t, ans, start, goal)
			validatePath(t, bigGrid, ans.Path, start, goal, true)
		})
	}
}

// TestBigGrid8DirEnclosedCellReachable: (7,7) is enclosed for 4-dir but
// becomes reachable in 8-dir via its diagonal neighbours.
func TestBigGrid8DirEnclosedCellReachable(t *testing.T) {
	start, goal := Point{0, 0}, Point{7, 7}

	// Confirm 4-dir cannot reach it.
	if dijkstraCost(bigGrid, start, goal, false) != -1 {
		t.Fatal("expected (7,7) to be unreachable in 4-dir mode (grid definition error?)")
	}

	// Confirm 8-dir can reach it (diagonal access via e.g. (6,6)→(7,7)).
	if dijkstraCost(bigGrid, start, goal, true) < 0 {
		t.Fatal("expected (7,7) to be reachable in 8-dir mode (diagonal neighbours are open)")
	}

	ans := Astar(bigGrid, start, goal, Chebyshev, true)
	requireFound(t, ans, start, goal)
	validatePath(t, bigGrid, ans.Path, start, goal, true)
}

// ─────────────────────────────────────────────
// 8-dir Dijkstra cross-validation
// ─────────────────────────────────────────────

// TestDijkstraCrossValidation8Dir: Chebyshev (the tightest admissible 8-dir
// heuristic for unit-cost diagonals) must always produce Dijkstra-optimal paths.
func TestDijkstraCrossValidation8Dir(t *testing.T) {
	pairs := []struct{ start, goal Point }{
		{start: Point{0, 0}, goal: Point{9, 9}},
		{start: Point{0, 5}, goal: Point{9, 5}},
		{start: Point{0, 0}, goal: Point{9, 0}},
		{start: Point{0, 0}, goal: Point{0, 9}},
		{start: Point{2, 2}, goal: Point{8, 8}},
		{start: Point{0, 0}, goal: Point{7, 7}}, // reachable via diagonal
		{start: Point{1, 0}, goal: Point{9, 9}},
		{start: Point{0, 9}, goal: Point{9, 0}},
	}

	for _, pp := range pairs {
		optimal := dijkstraCost(bigGrid, pp.start, pp.goal, true)
		if optimal < 0 {
			t.Logf("skip unreachable pair %v→%v (8-dir)", pp.start, pp.goal)
			continue
		}

		t.Run("Chebyshev", func(t *testing.T) {
			ans := Astar(bigGrid, pp.start, pp.goal, Chebyshev, true)
			if !ans.Found {
				t.Fatalf("%v→%v: A* Found=false; Dijkstra cost=%d", pp.start, pp.goal, optimal)
			}
			validatePath(t, bigGrid, ans.Path, pp.start, pp.goal, true)
			if got := len(ans.Path) - 1; got != optimal {
				t.Errorf("%v→%v: cost=%d, Dijkstra optimal=%d", pp.start, pp.goal, got, optimal)
			}
		})
	}
}

// ─────────────────────────────────────────────
// Efficiency / node-expansion tests
// ─────────────────────────────────────────────

// TestManhattanFewer4DirNodes: Manhattan (tightest admissible 4-dir heuristic)
// must expand no more nodes than Zero (uninformed Dijkstra).
func TestManhattanFewer4DirNodes(t *testing.T) {
	start, goal := Point{0, 0}, Point{9, 9}
	ansMnh := Astar(bigGrid, start, goal, Manhanttan, false)
	ansZero := Astar(bigGrid, start, goal, Zero, false)

	if !ansMnh.Found || !ansZero.Found {
		t.Fatal("both heuristics must find a path")
	}
	if ansMnh.NodeExp > ansZero.NodeExp {
		t.Errorf("Manhattan expanded %d nodes > Zero's %d — informed heuristic must be ≤ uninformed",
			ansMnh.NodeExp, ansZero.NodeExp)
	}
	t.Logf("node expansions — Manhattan: %d  Zero(Dijkstra): %d", ansMnh.NodeExp, ansZero.NodeExp)
}

// TestChebyshevFewer8DirNodes: Chebyshev must expand no more nodes than Zero
// in 8-dir mode (same argument — Chebyshev is the tightest admissible 8-dir heuristic).
func TestChebyshevFewer8DirNodes(t *testing.T) {
	start, goal := Point{0, 0}, Point{9, 9}
	ansChb := Astar(bigGrid, start, goal, Chebyshev, true)
	ansZero := Astar(bigGrid, start, goal, Zero, true)

	if !ansChb.Found || !ansZero.Found {
		t.Fatal("both heuristics must find a path")
	}
	if ansChb.NodeExp > ansZero.NodeExp {
		t.Errorf("Chebyshev expanded %d nodes > Zero's %d", ansChb.NodeExp, ansZero.NodeExp)
	}
	t.Logf("node expansions (8-dir) — Chebyshev: %d  Zero: %d", ansChb.NodeExp, ansZero.NodeExp)
}

// TestAllHeuristicsExpandPositiveNodes: NodeExp > 0 in every scenario.
func TestAllHeuristicsExpandPositiveNodes(t *testing.T) {
	cases := []struct {
		name        string
		start, goal Point
		diagonal    bool
	}{
		{name: "4dir_found", start: Point{0, 0}, goal: Point{9, 9}, diagonal: false},
		{name: "4dir_not_found", start: Point{0, 0}, goal: Point{7, 7}, diagonal: false},
		{name: "4dir_same", start: Point{5, 5}, goal: Point{5, 5}, diagonal: false},
		{name: "8dir_found", start: Point{0, 0}, goal: Point{9, 9}, diagonal: true},
		{name: "8dir_same", start: Point{5, 5}, goal: Point{5, 5}, diagonal: true},
	}
	for _, tc := range cases {
		for _, h := range allHeuristics {
			t.Run(tc.name+"/"+h.name, func(t *testing.T) {
				ans := Astar(bigGrid, tc.start, tc.goal, h.fn, tc.diagonal)
				if ans.NodeExp <= 0 {
					t.Errorf("NodeExp = %d, want > 0", ans.NodeExp)
				}
			})
		}
	}
}

// TestNodeExpansionSummary logs a side-by-side comparison table of all
// heuristics for a long path; useful for manual inspection.
func TestNodeExpansionSummary(t *testing.T) {
	start, goal := Point{0, 0}, Point{9, 9}
	t.Log("── 4-dir node-expansion summary ──────────────────")
	for _, h := range allHeuristics {
		ans := Astar(bigGrid, start, goal, h.fn, false)
		t.Logf("  %-12s  cost=%d  nodes=%d", h.name, len(ans.Path)-1, ans.NodeExp)
	}
	t.Log("── 8-dir node-expansion summary ──────────────────")
	for _, h := range allHeuristics {
		ans := Astar(bigGrid, start, goal, h.fn, true)
		t.Logf("  %-12s  cost=%d  nodes=%d", h.name, len(ans.Path)-1, ans.NodeExp)
	}
}
