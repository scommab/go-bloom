package bloom

import (
  "fmt"
  "testing"
  "os"
  "bufio"
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

func TestBloomTest(t *testing.T) {
  var data = GetLines("data")
  var test = GetLines("input")
  filter := new(BloomFilter)
  
  for _, d := range data {
    filter.Add(d)
  }


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

 fmt.Printf("BLOOM-FILTER is %x\n", filter.bloom_filter)

}
