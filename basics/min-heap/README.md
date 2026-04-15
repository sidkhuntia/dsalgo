# Min-Heap

A **Min-Heap** is a complete binary tree where every parent node is less than or equal to its children. This means the smallest element is always at the root, allowing O(1) access to the minimum value.

This implementation stores the heap as a flat `[]int` slice using the standard array-based binary tree layout:

- Parent of node `i` → `(i - 1) / 2`
- Left child of node `i` → `2*i + 1`
- Right child of node `i` → `2*i + 2`

---

## API

### `Insert(val int)`
Appends the value to the end of the heap and **bubbles it up** until the heap property is restored.

**Time complexity:** O(log n)

### `ExtractMin() int`
Removes and returns the minimum element (the root). The last element is moved to the root and **sifted down** to restore the heap property.

**Time complexity:** O(log n)

### `Len() int`
Returns the number of elements currently in the heap.

**Time complexity:** O(1)

---

## Usage

```go
package main

import (
    "fmt"
    minheap "dsalgo/basics/min-heap"
)

func main() {
    mh := &minheap.MinHeap{}

    mh.Insert(10)
    mh.Insert(3)
    mh.Insert(7)
    mh.Insert(1)
    mh.Insert(5)

    fmt.Println(mh.Len())        // 5
    fmt.Println(mh.ExtractMin()) // 1
    fmt.Println(mh.ExtractMin()) // 3
    fmt.Println(mh.ExtractMin()) // 5
}
```

Calling `ExtractMin` repeatedly will yield elements in ascending (sorted) order — this is the basis of **Heap Sort**.

---

## Complexity Summary

| Operation    | Time       | Space |
|--------------|------------|-------|
| `Insert`     | O(log n)   | O(1)  |
| `ExtractMin` | O(log n)   | O(1)  |
| `Len`        | O(1)       | O(1)  |
| Build heap   | O(n log n) | O(n)  |

---

## Common Use Cases

- **Priority queues** — always process the smallest (highest-priority) item next.
- **Heap Sort** — extract min repeatedly to sort a collection.
- **Dijkstra's algorithm** — efficiently fetch the unvisited node with the smallest tentative distance.
- **Kth smallest element** — maintain a heap of size K over a stream of numbers.