package bloomfilter

import (
	"hash"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const (
	MAX_SAFE_PRIME = 4294967291
)

type BloomFilter struct {
	bitSet        []bool
	hashFunctions []hash.Hash32
	hashCount     int
}

func NewBloomFilter(size int, hashCount int) *BloomFilter {
	hashFunctions := make([]hash.Hash32, hashCount)
	for i := 0; i < hashCount; i++ {
		hashFunctions[i] = murmur3.New32WithSeed(uint32(rand.Intn(MAX_SAFE_PRIME)))
	}
	return &BloomFilter{
		bitSet:        make([]bool, size),
		hashFunctions: hashFunctions,
		hashCount:     hashCount,
	}
}

func (bf *BloomFilter) Add(item string) {
	hashes := bf.computeHashes(item)
	for _, hash := range hashes {
		bf.bitSet[hash] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	hashes := bf.computeHashes(item)
	for _, hash := range hashes {
		if !bf.bitSet[hash] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) Clear() {
	for i := range bf.bitSet {
		bf.bitSet[i] = false
	}
}

func (bf *BloomFilter) Size() int {
	return len(bf.bitSet)
}

func (bf *BloomFilter) HashCount() int {
	return bf.hashCount
}

func (bf *BloomFilter) BitSet() []bool {
	return bf.bitSet
}

func (bf *BloomFilter) computeHashes(item string) []uint32 {
	hashes := make([]uint32, bf.hashCount)
	for i := 0; i < bf.hashCount; i++ {
		bf.hashFunctions[i].Write([]byte(item))
		hashes[i] = bf.hashFunctions[i].Sum32() % uint32(len(bf.bitSet))
		bf.hashFunctions[i].Reset()
	}
	return hashes
}

