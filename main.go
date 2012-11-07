package bloom 

import (
  "crypto/sha1"
  "encoding/binary"
  "bytes"
)

const FILTER_SIZE = 17 // in bytes
const FILTER_INDEX_SIZE = FILTER_SIZE * 8 // in bits

type BloomFilter struct {
  bloom_filter [FILTER_SIZE]byte
}

func Sha1(data []byte) []byte {
  var h = sha1.New()
  h.Write(data)
  return h.Sum(nil)
}

func ToIndex(data []byte, max_val int) int {
  buf := bytes.NewBuffer(data)
  result, err := binary.ReadUvarint(buf)
  if err != nil {
    panic(err)
  }
  return int(result) % max_val
}

func SetBit(bloom_filter *[FILTER_SIZE]byte, index int) {
  loc := index / 8
  bloom_filter[loc] |= 1 << uint(index % 8) 
}

func IsBitSet(bloom_filter [FILTER_SIZE]byte, index int) bool {
  loc := index / 8
  return (bloom_filter[loc] & byte(1 << uint(index % 8))) != 0
}

func DataToBloomIndex(d []byte) int {
  var hashes []byte = Sha1(d)
  index := ToIndex(hashes, FILTER_INDEX_SIZE)
  return index
}

func (b *BloomFilter) Add(d []byte) {
  output := DataToBloomIndex(d)
  SetBit(&b.bloom_filter, output)
}

func (b *BloomFilter) Has(d []byte) bool {
  output := DataToBloomIndex(d)
  return IsBitSet(b.bloom_filter, output)
}
