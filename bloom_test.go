package main

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

func TestBloomTest(t *testing.T) {
  var data = GetLines("data")
  var test = GetLines("input")
  //PrintAsStrings(data)
  var bloom_filter [FILTER_SIZE]byte
  
  for _, d := range data {
    output := DataToBloomIndex(d)

    SetBit(&bloom_filter, output)
  }

  for _, d := range test {
    output := DataToBloomIndex(d)

    fmt.Printf("%s", d)
    if IsBitSet(bloom_filter, output) {
      fmt.Printf(" True\n")
    } else {
      fmt.Printf(" False\n")
    }
  }

 fmt.Printf("BLOOM-FILTER is %x\n", bloom_filter)

}
