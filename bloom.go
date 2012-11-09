package bloom

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
)

const BITS_IN_BLOOM_TYPE = 8

type BloomFilter struct {
	bloom_filter []byte
	size         uint
	probe        int
}

// Create a new Bloom filter.
// size - describes the number of bytes the Bloom filter.
// probe - the number of probes each insert/lookup uses.
func MakeBloomFilter(size uint, probe int) *BloomFilter {
	result := new(BloomFilter)
	result.bloom_filter = make([]byte, size)
	result.size = size
	result.probe = probe
	return result
}

func Sha1(data []byte) []byte {
	var h = sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func ToIndex(data []byte) uint {
	buf := bytes.NewBuffer(data)
	result, err := binary.ReadUvarint(buf)
	if err != nil {
		panic(err)
	}
	return uint(result)
}

func (b *BloomFilter) SetBit(index uint) {
	index = index % (b.size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	b.bloom_filter[loc] |= 1 << uint(index%BITS_IN_BLOOM_TYPE)
}

func (b *BloomFilter) IsBitSet(index uint) bool {
	index = index % (b.size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	return (b.bloom_filter[loc] & byte(1<<uint(index%BITS_IN_BLOOM_TYPE))) != 0
}

func DataToBloomIndex(d []byte) ([]byte, uint) {
	var hash []byte = Sha1(d)
	index := ToIndex(hash)
	return hash, index
}

// Add an item to the Bloom filter
func (b *BloomFilter) Add(d []byte) {
	hash, output := DataToBloomIndex(d)
	b.SetBit(output)
	for i := 0; i < b.probe-1; i++ {
		hash, output = DataToBloomIndex(hash)
		b.SetBit(output)
	}
}

// Check if a value is in the Bloom filter.
// Returns False if the value definitely isn't in the Bloom filter.
// Returns True if the value could be in the Bloom filter.
func (b *BloomFilter) Has(d []byte) bool {
	result := true
	hash, output := DataToBloomIndex(d)
	result = result && b.IsBitSet(output)
	for i := 0; i < b.probe-1; i++ {
		hash, output = DataToBloomIndex(hash)
		result = result && b.IsBitSet(output)
	}
	return result
}
