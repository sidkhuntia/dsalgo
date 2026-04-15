package minheap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// MinHeap[int]
// ---------------------------------------------------------------------------

func TestInt_Len_EmptyHeap(t *testing.T) {
	mh := &MinHeap[int]{}
	assert.Equal(t, 0, mh.Len())
}

func TestInt_Insert_SingleElement(t *testing.T) {
	mh := &MinHeap[int]{}
	mh.Insert(42)
	assert.Equal(t, 1, mh.Len())
}

func TestInt_Insert_MultipleElements_Len(t *testing.T) {
	mh := &MinHeap[int]{}
	values := []int{5, 3, 8, 1, 9, 2}
	for i, v := range values {
		mh.Insert(v)
		assert.Equal(t, i+1, mh.Len())
	}
}

func TestInt_ExtractMin_SingleElement(t *testing.T) {
	mh := &MinHeap[int]{}
	mh.Insert(7)
	assert.Equal(t, 7, mh.ExtractMin())
	assert.Equal(t, 0, mh.Len())
}

func TestInt_ExtractMin_ReturnsMinimum(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{10, 3, 7, 1, 5} {
		mh.Insert(v)
	}
	assert.Equal(t, 1, mh.ExtractMin())
}

func TestInt_ExtractMin_ReducesLen(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{4, 2, 6} {
		mh.Insert(v)
	}
	mh.ExtractMin()
	assert.Equal(t, 2, mh.Len())
}

// Repeated ExtractMin calls must yield values in ascending order.
func TestInt_ExtractMin_YieldsSortedOrder(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{9, 3, 7, 1, 5, 2, 8, 4, 6} {
		mh.Insert(v)
	}
	prev := mh.ExtractMin()
	for mh.Len() > 0 {
		curr := mh.ExtractMin()
		assert.LessOrEqual(t, prev, curr, "heap order violated")
		prev = curr
	}
}

// Every insert triggers a full bubble-up when values arrive in descending order.
func TestInt_Insert_DescendingOrder(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{10, 8, 6, 4, 2} {
		mh.Insert(v)
	}
	assert.Equal(t, 2, mh.ExtractMin())
}

// No bubble-up needed when values arrive in ascending order.
func TestInt_Insert_AscendingOrder(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{1, 3, 5, 7, 9} {
		mh.Insert(v)
	}
	assert.Equal(t, 1, mh.ExtractMin())
}

func TestInt_ExtractMin_WithDuplicates(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{4, 4, 4, 1, 1} {
		mh.Insert(v)
	}
	expected := []int{1, 1, 4, 4, 4}
	for _, want := range expected {
		assert.Equal(t, want, mh.ExtractMin())
	}
}

func TestInt_ExtractMin_NegativeNumbers(t *testing.T) {
	mh := &MinHeap[int]{}
	for _, v := range []int{0, -5, 3, -1, -10, 7} {
		mh.Insert(v)
	}
	assert.Equal(t, -10, mh.ExtractMin())
}

func TestInt_RepeatedInsertAndExtract(t *testing.T) {
	mh := &MinHeap[int]{}

	mh.Insert(5)
	mh.Insert(2)
	assert.Equal(t, 2, mh.ExtractMin())

	mh.Insert(1)
	mh.Insert(8)
	assert.Equal(t, 1, mh.ExtractMin())
	assert.Equal(t, 5, mh.ExtractMin())
	assert.Equal(t, 8, mh.ExtractMin())
	assert.Equal(t, 0, mh.Len())
}

// ---------------------------------------------------------------------------
// MinHeap[float64]
// ---------------------------------------------------------------------------

func TestFloat64_ExtractMin_ReturnsMinimum(t *testing.T) {
	mh := &MinHeap[float64]{}
	for _, v := range []float64{3.14, 2.71, 1.41, 1.73} {
		mh.Insert(v)
	}
	assert.InDelta(t, 1.41, mh.ExtractMin(), 1e-9)
}

func TestFloat64_ExtractMin_YieldsSortedOrder(t *testing.T) {
	mh := &MinHeap[float64]{}
	for _, v := range []float64{0.5, 3.3, 1.1, 2.2, 0.1} {
		mh.Insert(v)
	}
	prev := mh.ExtractMin()
	for mh.Len() > 0 {
		curr := mh.ExtractMin()
		assert.LessOrEqual(t, prev, curr, "heap order violated")
		prev = curr
	}
}

func TestFloat64_ExtractMin_WithNegativeAndZero(t *testing.T) {
	mh := &MinHeap[float64]{}
	for _, v := range []float64{0.0, -1.5, 2.5, -3.7} {
		mh.Insert(v)
	}
	assert.InDelta(t, -3.7, mh.ExtractMin(), 1e-9)
}

func TestFloat64_ExtractMin_WithDuplicates(t *testing.T) {
	mh := &MinHeap[float64]{}
	for _, v := range []float64{1.1, 1.1, 0.5, 0.5} {
		mh.Insert(v)
	}
	expected := []float64{0.5, 0.5, 1.1, 1.1}
	for _, want := range expected {
		assert.InDelta(t, want, mh.ExtractMin(), 1e-9)
	}
}

// ---------------------------------------------------------------------------
// MinHeap[string]
// ---------------------------------------------------------------------------

func TestString_ExtractMin_ReturnsLexicographicMinimum(t *testing.T) {
	mh := &MinHeap[string]{}
	for _, v := range []string{"banana", "apple", "cherry", "date"} {
		mh.Insert(v)
	}
	assert.Equal(t, "apple", mh.ExtractMin())
}

func TestString_ExtractMin_YieldsLexicographicOrder(t *testing.T) {
	mh := &MinHeap[string]{}
	for _, v := range []string{"zebra", "mango", "apple", "fig", "kiwi"} {
		mh.Insert(v)
	}
	prev := mh.ExtractMin()
	for mh.Len() > 0 {
		curr := mh.ExtractMin()
		assert.LessOrEqual(t, prev, curr, "heap order violated")
		prev = curr
	}
}

func TestString_ExtractMin_WithDuplicates(t *testing.T) {
	mh := &MinHeap[string]{}
	for _, v := range []string{"b", "a", "a", "c", "b"} {
		mh.Insert(v)
	}
	expected := []string{"a", "a", "b", "b", "c"}
	for _, want := range expected {
		assert.Equal(t, want, mh.ExtractMin())
	}
}

func TestString_Len_AfterInsertAndExtract(t *testing.T) {
	mh := &MinHeap[string]{}
	assert.Equal(t, 0, mh.Len())
	mh.Insert("hello")
	mh.Insert("world")
	assert.Equal(t, 2, mh.Len())
	mh.ExtractMin()
	assert.Equal(t, 1, mh.Len())
	mh.ExtractMin()
	assert.Equal(t, 0, mh.Len())
}
