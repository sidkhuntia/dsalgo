package bloomfilter_test

import (
	"math"
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
	p := 0.01
	n := 1000000
	bf := bloomfilter.NewBloomFilter(p, n)
	assert.NotNil(t, bf)
	assert.Equal(t, int(math.Ceil(-float64(n)*math.Log(p)/(math.Log(2)*math.Log(2)))), bf.Size())
	assert.Equal(t, int64(math.Ceil(float64(bf.Size())/float64(n)*math.Log(2))), bf.HashCount())
	assert.Len(t, bf.BitSet(), int(math.Ceil(-float64(n)*math.Log(p)/(math.Log(2)*math.Log(2)))))
}

func TestAddAndContains(t *testing.T) {
	p := 0.01
	n := 1000000
	bf := bloomfilter.NewBloomFilter(p, n)

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
	testWordsCount := (math.MaxInt64)
	for range testWordsCount {
		word := randomWord(rand.Intn(10) + 1)
		if bf.Contains(word) && !slices.Contains(addedWords, word) {
			falsePositives++
		}
	}
	t.Logf("False positives Percentage: %.2f%%", float64(falsePositives)/float64(testWordsCount)*100)
}

func TestClear(t *testing.T) {
	p := 0.01
	n := 1000000
	bf := bloomfilter.NewBloomFilter(p, n)
	bf.Add("test")
	assert.True(t, bf.Contains("test"))
	bf.Clear()
	assert.False(t, bf.Contains("test"))
}
