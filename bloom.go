package bloom

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
)

const BITS_IN_BLOOM_TYPE = 64

type BloomFilter struct {
	Bloom_filter []uint64
	Size         uint
	Probe        int
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
	index = index % (b.Size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	b.Bloom_filter[loc] |= 1 << (index % BITS_IN_BLOOM_TYPE)
}

func (b *BloomFilter) isBitSet(index uint) bool {
	index = index % (b.Size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	return (b.Bloom_filter[loc] & (1 << (index % BITS_IN_BLOOM_TYPE))) != 0
}

func dataToBloomIndex(d []byte) ([]byte, uint) {
	var hash []byte = sha1_hash(d)
	index := toIndex(hash)
	return hash, index
}

// Create a new Bloom filter.
// size - the number of int64 to use for the Bloom filter.
// probe - the number of probes each insert/lookup uses.
func MakeBloomFilter(size uint, probe int) *BloomFilter {
	result := new(BloomFilter)
	result.Bloom_filter = make([]uint64, size)
	result.Size = size
	result.Probe = probe
	return result
}

// Export the bloom filter (encoded in Json)
func (b *BloomFilter) Export() ([]byte, error) {
	r, err := json.Marshal(b)
	return r, err
}

// Load an exported bloom filter
func MakeBloomFilterFromJson(jsonBlob []byte) (*BloomFilter, error) {
	result := new(BloomFilter)
	err := json.Unmarshal(jsonBlob, &result)
	return result, err
}

//Add an item to the Bloom filter
func (b *BloomFilter) Add(d []byte) {
	hash, output := dataToBloomIndex(d)
	b.setBit(output)
	for i := 0; i < b.Probe-1; i++ {
		hash, output = dataToBloomIndex(hash)
		b.setBit(output)
	}
}

// Check if a value is in the Bloom filter.
// Returns False if the value definitely isn't in the Bloom filter.
// Returns True if the value could be in the Bloom filter.
func (b *BloomFilter) Has(d []byte) bool {
	hash, output := dataToBloomIndex(d)
	if !b.isBitSet(output) {
		return false
	}
	for i := 0; i < b.Probe-1; i++ {
		hash, output = dataToBloomIndex(hash)
		if !b.isBitSet(output) {
			return false
		}
	}
	return true
}
