package main


import (
  "fmt"
  "bufio"
  "crypto/md5"
  "crypto/sha1"
  "crypto/sha256"
  "encoding/binary"
  "bytes"
  "os"
)

const FILTER_SIZE = 32 // in bytes
const FILTER_INDEX_SIZE = FILTER_SIZE * 8 // in bits

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

func PrintAsStrings(data [][]byte) {
  for _, d := range data {
    fmt.Println(string(d))
  }
}


func Md5sum(data []byte) []byte {
  var md5 = md5.New()
  md5.Write(data)
  return md5.Sum(nil)
}


func Sha1(data []byte) []byte {
  var h = sha1.New()
  h.Write(data)
  return h.Sum(nil)
}

func Sha256(data []byte) []byte {
  var h = sha256.New()
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
  bloom_filter[loc] |= 2 << uint(index % 8) 
}

func IsBitSet(bloom_filter [FILTER_SIZE]byte, index int) bool {
  loc := index / 8
  return (bloom_filter[loc] & byte(2 << uint(index % 8))) != 0
}

func DataToBloomIndex(d []byte) int {
  var hashes []byte = Md5sum(d)
  index := ToIndex(hashes, FILTER_INDEX_SIZE)

  fmt.Printf("%s, %x\n", d, hashes)
  fmt.Printf("%s, %x\n", d, index)

  return index
}
func main() {
  var data = GetLines("data")
  var test = GetLines("input")
  PrintAsStrings(data)
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
