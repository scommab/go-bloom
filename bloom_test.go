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
	binary.Write(buf, binary.LittleEndian, filter.Bloom_filter)
	fmt.Printf("BLOOM FILTER: % x\n", buf.Bytes())
}

func getTest1Data() ([][]byte, [][]byte, *BloomFilter) {
	return GetLines("test_data/test1_keys"),
		GetLines("test_data/test1_invalid"),
		MakeBloomFilter(1, 2)
}

func TestBloom(t *testing.T) {
	data, test, filter := getTest1Data()

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

func TestBloomExportImport(t *testing.T) {
	data, test, filter := getTest1Data()

	for _, d := range data {
		filter.Add(d)
	}

	blob, _ := filter.Export()
	filter2, _ := MakeBloomFilterFromJson(blob)

	for _, d := range data {
		if printHas(filter2, d) != true {
			t.Fatalf("Failed Match")
		}
	}

	for _, d := range test {
		if printHas(filter2, d) != false {
			t.Fatalf("Failed Match")
		}
	}
}

func BenchmarkBloomBuild(b *testing.B) {
	data, _, filter := getTest1Data()

	for i := 0; i < b.N; i++ {
		for _, d := range data {
			filter.Add(d)
		}
	}
}

func BenchmarkBloomExport(b *testing.B) {
	data, _, filter := getTest1Data()

	for _, d := range data {
		filter.Add(d)
	}

	for i := 0; i < b.N; i++ {
		filter.Export()
	}
}

func BenchmarkBloomImport(b *testing.B) {
	data, _, filter := getTest1Data()

	for _, d := range data {
		filter.Add(d)
	}
	blob, _ := filter.Export()

	for i := 0; i < b.N; i++ {
		MakeBloomFilterFromJson(blob)
	}
}

func BenchmarkBloomLookup(b *testing.B) {
	data, test, filter := getTest1Data()

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
