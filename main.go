package main


import (
  "fmt"
  "bufio"
  "crypto/md5"
  "crypto/sha1"
  "crypto/sha256"
  "os"
)

const FILTER_SIZE = 32

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


func ToXBytes(data []byte) [FILTER_SIZE]byte {
  var result [FILTER_SIZE]byte
  for i, d := range data {
    if i >= FILTER_SIZE {
      break
    }
    result[i] = d
  }
  return result
}

func Md5sum(data []byte) [FILTER_SIZE]byte {
  var md5 = md5.New()
  md5.Write(data)
  return ToXBytes(md5.Sum(nil))
}


func Sha1(data []byte) [FILTER_SIZE]byte {
  var h = sha1.New()
  h.Write(data)
  return ToXBytes(h.Sum(nil))
}

func Sha256(data []byte) [FILTER_SIZE]byte {
  var h = sha256.New()
  h.Write(data)
  return ToXBytes(h.Sum(nil))
}

func Or(d1 [FILTER_SIZE]byte, d2[FILTER_SIZE]byte) [FILTER_SIZE]byte {
  var result [FILTER_SIZE]byte
  for i, _ := range d1 {
    result[i] = d1[i] | d2[i]
  }
  return result
}

func OrHashes(data []byte) [FILTER_SIZE]byte {
  var result [FILTER_SIZE]byte

  result = Or(result, Md5sum(data))
  result = Or(result, Sha1(data))
  result = Or(result, Sha256(data))

  return result
}

func main() {
  var data = GetLines("data")
  PrintAsStrings(data)
  var bloom_filter [FILTER_SIZE]byte
  
  for _, d := range data {
    var hashes [FILTER_SIZE]byte = OrHashes(d)
    fmt.Printf("%s, %x\n", d, hashes)
    bloom_filter = Or(bloom_filter, hashes)
  }

  fmt.Printf("BLOOM-FILTER is %x\n", bloom_filter)

}
