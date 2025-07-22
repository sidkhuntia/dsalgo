package bloomfilter_test

import (
	"math/rand"
	"slices"
	"testing"
	"time"

	bloomfilter "dsalgo/bloom-filter" // Corrected import path

	"github.com/stretchr/testify/assert"
)

func randomWord(length int) string {
	rand.NewSource(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	word := make([]rune, length)
	for i := range word {
		word[i] = letters[rand.Intn(len(letters))]
	}
	return string(word)
}

func TestNewBloomFilter(t *testing.T) {
	size := 1000
	hashCount := 3
	bf := bloomfilter.NewBloomFilter(size, hashCount)
	assert.NotNil(t, bf)
	assert.Equal(t, size, bf.Size())
	assert.Equal(t, hashCount, bf.HashCount())
	assert.Len(t, bf.BitSet(), size)
}

func TestAddAndContains(t *testing.T) {
	bf := bloomfilter.NewBloomFilter(1000000, 100)

	// Add random words
	addedWords := make([]string, 10000)
	for i := range addedWords {
		word := randomWord(rand.Intn(10) + 1)
		addedWords[i] = word
		bf.Add(word)
	}

	// Check added words
	for _, word := range addedWords {
		assert.True(t, bf.Contains(word), "Expected true for added word: %s", word)
	}

	// Check non-added words (may have false positives, but for testing we can check a few)
	falsePositives := 0
	testWordsCount := 1000000
	for range testWordsCount {
		word := randomWord(rand.Intn(10) + 1)
		if bf.Contains(word) && !slices.Contains(addedWords, word) {
			falsePositives++
		}
	}
	t.Logf("False positives Percentage: %.2f%%", float64(falsePositives)/float64(testWordsCount)*100)
}

func TestClear(t *testing.T) {
	bf := bloomfilter.NewBloomFilter(1000, 3)
	bf.Add("test")
	assert.True(t, bf.Contains("test"))
	bf.Clear()
	assert.False(t, bf.Contains("test"))
}
