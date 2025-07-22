package bloomfilter

import (
	"hash"
	"math"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const (
	MAX_SAFE_PRIME = 4294967291
)

type BloomFilter struct {
	bitSet        []bool
	hashFunctions []hash.Hash64
	hashCount     int64
}

func NewBloomFilter(p float64, n int) *BloomFilter {
	m := -float64(n) * math.Log(p) / (math.Log(2) * math.Log(2))
	m_int := int64(math.Ceil(m))
	k := m / float64(n) * math.Log(2)
	k_int := int64(math.Ceil(k))
	hashFunctions := make([]hash.Hash64, k_int)
	for i := range k_int {
		hashFunctions[i] = murmur3.New64WithSeed(uint32(rand.Intn(MAX_SAFE_PRIME)))
	}
	return &BloomFilter{
		bitSet:        make([]bool, m_int),
		hashFunctions: hashFunctions,
		hashCount:     k_int,
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

func (bf *BloomFilter) HashCount() int64 {
	return bf.hashCount
}

func (bf *BloomFilter) BitSet() []bool {
	return bf.bitSet
}

func (bf *BloomFilter) computeHashes(item string) []uint64 {
	hashes := make([]uint64, bf.hashCount)
	for i := range bf.hashCount {
		bf.hashFunctions[i].Write([]byte(item))
		hashes[i] = bf.hashFunctions[i].Sum64() % uint64(len(bf.bitSet))
		bf.hashFunctions[i].Reset()
	}
	return hashes
}

