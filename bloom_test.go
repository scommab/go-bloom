package bloom

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func GetLines(file_name string) [][]byte {
	var result [][]byte // empty array

	f, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	line, _, err := r.ReadLine()
	for err == nil {
		// turn to a sting for easier printing
		result = append(result, line)
		line, _, err = r.ReadLine()
	}
	return result
}

func printHas(filter *BloomFilter, d []byte) bool {
	has := filter.Has(d)

	fmt.Printf("%s", d)
	if has {
		fmt.Printf(" True\n")
	} else {
		fmt.Printf(" False\n")
	}
	return has
}

func printFilter(filter *BloomFilter) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, filter.bloom_filter)
	fmt.Printf("BLOOM FILTER: % x\n", buf.Bytes())
}

func TestBloom(t *testing.T) {
	var data = GetLines("test_data/test1_keys")
	var test = GetLines("test_data/test1_invalid")
	filter := MakeBloomFilter(1, 2) //new(BloomFilter)

	for _, d := range data {
		filter.Add(d)
	}
	printFilter(filter)

	for _, d := range data {
		if printHas(filter, d) != true {
			t.Fatalf("Failed Match")
		}
	}

	for _, d := range test {
		if printHas(filter, d) != false {
			t.Fatalf("Failed Match")
		}
	}
}

func BenchmarkBloomBuild(b *testing.B) {
	var data = GetLines("test_data/test1_keys")
	filter := MakeBloomFilter(1, 2) //new(BloomFilter)

	for i := 0; i < b.N; i++ {
		for _, d := range data {
			filter.Add(d)
		}
	}
}

func BenchmarkBloomLookup(b *testing.B) {
	var data = GetLines("test_data/test1_keys")
	var test = GetLines("test_data/test1_invalid")
	filter := MakeBloomFilter(1, 2) //new(BloomFilter)

	for _, d := range data {
		filter.Add(d)
	}
	for i := 0; i < b.N; i++ {
		for _, d := range data {
			filter.Has(d)
		}

		for _, d := range test {
			filter.Has(d)
		}
	}
}
