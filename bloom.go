package bloom

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
)

const BITS_IN_BLOOM_TYPE = 8

type BloomFilter struct {
	bloom_filter []byte
	size         int
	probe        int
}

func MakeBloomFilter(size int, probe int) *BloomFilter {
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

func ToIndex(data []byte) int {
	buf := bytes.NewBuffer(data)
	result_64, err := binary.ReadVarint(buf)
	result := int(result_64)
	if err != nil {
		panic(err)
	}
	if result < 0 {
		result = -result
	}
	return result
}

func (b *BloomFilter) SetBit(index int) {
	index = index % (b.size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	b.bloom_filter[loc] |= 1 << uint(index%BITS_IN_BLOOM_TYPE)
}

func (b *BloomFilter) IsBitSet(index int) bool {
	index = index % (b.size * BITS_IN_BLOOM_TYPE)
	loc := index / BITS_IN_BLOOM_TYPE
	return (b.bloom_filter[loc] & byte(1<<uint(index%BITS_IN_BLOOM_TYPE))) != 0
}

func DataToBloomIndex(d []byte) ([]byte, int) {
	var hash []byte = Sha1(d)
	index := ToIndex(hash)
	return hash, index
}

func (b *BloomFilter) Add(d []byte) {
	hash, output := DataToBloomIndex(d)
	b.SetBit(output)
	for i := 0; i < b.probe-1; i++ {
		hash, output = DataToBloomIndex(hash)
		b.SetBit(output)
	}
}

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
