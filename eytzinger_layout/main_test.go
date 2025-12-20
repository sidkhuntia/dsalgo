package eytzingerlayout

import (
	"fmt"
	"testing"
	"time"
)

func TestEytzingerLayout(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Printf("Input:           %v\n", input)
	fmt.Printf("Eytzinger Layout: ")
	eytzinger := NewEytzingerLayout(input)
	eytzinger.Print()

	fmt.Printf("Searching for 5: %d\n", eytzinger.Search(5))
	fmt.Printf("Searching for 9: %d\n", eytzinger.Search(9))
}

// Standard binary search implementation
func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

func generateSortedArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}
	return arr
}

func BenchmarkEytzingerSearch_100(b *testing.B) {
	arr := generateSortedArray(100)
	eytzinger := NewEytzingerLayout(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eytzinger.Search(i % 100)
	}
}

func BenchmarkBinarySearch_100(b *testing.B) {
	arr := generateSortedArray(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binarySearch(arr, i%100)
	}
}

func BenchmarkEytzingerSearch_1000(b *testing.B) {
	arr := generateSortedArray(1000)
	eytzinger := NewEytzingerLayout(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eytzinger.Search(i % 1000)
	}
}

func BenchmarkBinarySearch_1000(b *testing.B) {
	arr := generateSortedArray(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binarySearch(arr, i%1000)
	}
}

func BenchmarkEytzingerSearch_10000(b *testing.B) {
	arr := generateSortedArray(10000)
	eytzinger := NewEytzingerLayout(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eytzinger.Search(i % 10000)
	}
}

func BenchmarkBinarySearch_10000(b *testing.B) {
	arr := generateSortedArray(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binarySearch(arr, i%10000)
	}
}

func BenchmarkEytzingerSearch_100000(b *testing.B) {
	arr := generateSortedArray(100000)
	eytzinger := NewEytzingerLayout(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eytzinger.Search(i % 100000)
	}
}

func BenchmarkBinarySearch_100000(b *testing.B) {
	arr := generateSortedArray(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binarySearch(arr, i%100000)
	}
}

func TestAverageSearchTime(t *testing.T) {
	sizes := []int{100, 1000, 10000, 100000, 1000000, 10000000, 1000000000000}
	iterations := 100000

	fmt.Println("\n=== Average Search Time Comparison ===")
	fmt.Printf("%-10s %-20s %-20s %-15s\n", "Size", "Eytzinger (ns)", "Binary Search (ns)", "Speedup (%)")
	fmt.Println("----------------------------------------------------------------------------")

	for _, size := range sizes {
		arr := generateSortedArray(size)
		eytzinger := NewEytzingerLayout(arr)

		// Measure Eytzinger search
		start := time.Now()
		for i := 0; i < iterations; i++ {
			eytzinger.Search(i % size)
		}
		eytzingerTime := time.Since(start).Nanoseconds() / int64(iterations)

		// Measure Binary search
		start = time.Now()
		for i := 0; i < iterations; i++ {
			binarySearch(arr, i%size)
		}
		binaryTime := time.Since(start).Nanoseconds() / int64(iterations)

		// Calculate speedup percentage
		speedup := float64(binaryTime-eytzingerTime) / float64(binaryTime) * 100

		fmt.Printf("%-10d %-20d %-20d %-15.2f\n", size, eytzingerTime, binaryTime, speedup)
	}
}
