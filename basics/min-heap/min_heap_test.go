package minheap

import (
	"testing"
)

func TestLen_EmptyHeap(t *testing.T) {
	mh := &MinHeap{}
	if mh.Len() != 0 {
		t.Errorf("expected Len 0 on empty heap, got %d", mh.Len())
	}
}

func TestInsert_SingleElement(t *testing.T) {
	mh := &MinHeap{}
	mh.Insert(42)
	if mh.Len() != 1 {
		t.Errorf("expected Len 1 after one insert, got %d", mh.Len())
	}
}

func TestInsert_MultipleElements_Len(t *testing.T) {
	mh := &MinHeap{}
	values := []int{5, 3, 8, 1, 9, 2}
	for i, v := range values {
		mh.Insert(v)
		if mh.Len() != i+1 {
			t.Errorf("expected Len %d after %d inserts, got %d", i+1, i+1, mh.Len())
		}
	}
}

func TestExtractMin_SingleElement(t *testing.T) {
	mh := &MinHeap{}
	mh.Insert(7)
	min := mh.ExtractMin()
	if min != 7 {
		t.Errorf("expected ExtractMin to return 7, got %d", min)
	}
	if mh.Len() != 0 {
		t.Errorf("expected Len 0 after extracting only element, got %d", mh.Len())
	}
}

func TestExtractMin_ReturnsMinimum(t *testing.T) {
	mh := &MinHeap{}
	for _, v := range []int{10, 3, 7, 1, 5} {
		mh.Insert(v)
	}
	min := mh.ExtractMin()
	if min != 1 {
		t.Errorf("expected ExtractMin to return 1, got %d", min)
	}
}

func TestExtractMin_ReducesLen(t *testing.T) {
	mh := &MinHeap{}
	for _, v := range []int{4, 2, 6} {
		mh.Insert(v)
	}
	mh.ExtractMin()
	if mh.Len() != 2 {
		t.Errorf("expected Len 2 after one ExtractMin, got %d", mh.Len())
	}
}

// Extracting all elements should yield values in ascending (sorted) order.
func TestExtractMin_YieldsSortedOrder(t *testing.T) {
	mh := &MinHeap{}
	input := []int{9, 3, 7, 1, 5, 2, 8, 4, 6}
	for _, v := range input {
		mh.Insert(v)
	}

	prev := mh.ExtractMin()
	for mh.Len() > 0 {
		curr := mh.ExtractMin()
		if curr < prev {
			t.Errorf("heap order violated: extracted %d after %d", curr, prev)
		}
		prev = curr
	}
}

// Insert in strictly descending order — every insert triggers a full bubble-up.
func TestInsert_DescendingOrder_HeapProperty(t *testing.T) {
	mh := &MinHeap{}
	for _, v := range []int{10, 8, 6, 4, 2} {
		mh.Insert(v)
	}
	min := mh.ExtractMin()
	if min != 2 {
		t.Errorf("expected minimum 2 after descending inserts, got %d", min)
	}
}

// Insert in strictly ascending order — no bubble-up needed.
func TestInsert_AscendingOrder_HeapProperty(t *testing.T) {
	mh := &MinHeap{}
	for _, v := range []int{1, 3, 5, 7, 9} {
		mh.Insert(v)
	}
	min := mh.ExtractMin()
	if min != 1 {
		t.Errorf("expected minimum 1 after ascending inserts, got %d", min)
	}
}

func TestExtractMin_WithDuplicates(t *testing.T) {
	mh := &MinHeap{}
	for _, v := range []int{4, 4, 4, 1, 1} {
		mh.Insert(v)
	}

	expected := []int{1, 1, 4, 4, 4}
	for _, want := range expected {
		got := mh.ExtractMin()
		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestExtractMin_NegativeNumbers(t *testing.T) {
	mh := &MinHeap{}
	for _, v := range []int{0, -5, 3, -1, -10, 7} {
		mh.Insert(v)
	}
	min := mh.ExtractMin()
	if min != -10 {
		t.Errorf("expected minimum -10, got %d", min)
	}
}

func TestHeap_RepeatedInsertAndExtract(t *testing.T) {
	mh := &MinHeap{}

	mh.Insert(5)
	mh.Insert(2)
	if got := mh.ExtractMin(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}

	mh.Insert(1)
	mh.Insert(8)
	if got := mh.ExtractMin(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}

	if got := mh.ExtractMin(); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
	if got := mh.ExtractMin(); got != 8 {
		t.Errorf("expected 8, got %d", got)
	}
	if mh.Len() != 0 {
		t.Errorf("expected empty heap, got Len %d", mh.Len())
	}
}
