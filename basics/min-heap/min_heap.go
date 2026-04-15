package minheap

import "cmp"

type MinHeap[T cmp.Ordered] struct {
	nodes []T
}

func (mh MinHeap[T]) Len() int {
	return len(mh.nodes)
}

func (mh *MinHeap[T]) Insert(a T) {
	mh.nodes = append(mh.nodes, a)

	// BUBBLE UP
	i := mh.Len() - 1
	for i > 0 && mh.nodes[i] < mh.nodes[(i-1)/2] {
		parent := (i - 1) / 2
		mh.nodes[parent], mh.nodes[i] = mh.nodes[i], mh.nodes[parent]
		i = parent
	}
}

func (mh *MinHeap[T]) ExtractMin() T {
	ans := mh.nodes[0]

	mh.nodes[0] = mh.nodes[mh.Len()-1]
	mh.nodes = mh.nodes[:mh.Len()-1]

	i := 0

	for {
		left, right, smallest := 2*i+1, 2*i+2, i

		if left < mh.Len() && mh.nodes[left] < mh.nodes[smallest] {
			smallest = left
		}
		if right < mh.Len() && mh.nodes[right] < mh.nodes[smallest] {
			smallest = right
		}

		if smallest == i {
			break
		}

		mh.nodes[i], mh.nodes[smallest] = mh.nodes[smallest], mh.nodes[i]
		i = smallest
	}

	return ans
}
