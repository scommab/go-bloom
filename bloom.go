package bloom

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
)

const BITS_IN_BLOOM_TYPE = 64

type BloomFilter struct {
	bloom_filter []uint64
	size         uint
	probe        int
}

// Create a new Bloom filter.
// size - describes the number of bytes the Bloom filter.
// probe - the number of probes each insert/lookup uses.
func MakeBloomFilter(size uint, probe int) *BloomFilter {
	result := new(BloomFilter)
	result.bloom_filter = make([]uint64, size)
	result.size = size
	result.probe = probe
	return result
}

func sha1_hash(data []byte) []byte {
	var h = sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func toIndex(data []byte) uint {
	buf := bytes.NewBuffer(data)
	result, err := binary.ReadUvarint(buf)
	if err != nil {
		panic(err)
	}
	return uint(result)
}

func (b *BloomFilter) setBit(index uint) {
	index = index % (b.size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	b.bloom_filter[loc] |= 1 << uint(index % BITS_IN_BLOOM_TYPE)
}

func (b *BloomFilter) isBitSet(index uint) bool {
	index = index % (b.size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	return (b.bloom_filter[loc] & (1<<uint(index % BITS_IN_BLOOM_TYPE))) != 0
}

func dataToBloomIndex(d []byte) ([]byte, uint) {
	var hash []byte = sha1_hash(d)
	index := toIndex(hash)
	return hash, index
}

// Add an item to the Bloom filter
func (b *BloomFilter) Add(d []byte) {
	hash, output := dataToBloomIndex(d)
	b.setBit(output)
	for i := 0; i < b.probe-1; i++ {
		hash, output = dataToBloomIndex(hash)
		b.setBit(output)
	}
}

// Check if a value is in the Bloom filter.
// Returns False if the value definitely isn't in the Bloom filter.
// Returns True if the value could be in the Bloom filter.
func (b *BloomFilter) Has(d []byte) bool {
	result := true
	hash, output := dataToBloomIndex(d)
	result = result && b.isBitSet(output)
	for i := 0; i < b.probe-1; i++ {
		hash, output = dataToBloomIndex(hash)
		result = result && b.isBitSet(output)
	}
	return result
}
