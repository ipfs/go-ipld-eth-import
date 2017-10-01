package main

import (
  "bufio"
  "bytes"
  "encoding/hex"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "time"
)

func main() {
  DagPutCmd := exec.Command("ipfs", "dag", "put", "--input-enc", "raw", "--format", "eth-block")

  var bufOut, bufErr bytes.Buffer

  DagPutCmd.Stdin = bytes.NewBuffer(getBytes())

  DagPutCmd.Stdout = &bufOut
  DagPutCmd.Stderr = &bufErr

  // The gist of this business
  start := time.Now().UnixNano()
  err := DagPutCmd.Run()
  fmt.Printf("Time Diff: %.0d\n", time.Now().UnixNano()-start)

  if err != nil {
    fmt.Printf("%s\n", bufErr.String())
  } else {
    fmt.Printf("%s\n", bufOut.String())
  }
}

// Hygiene
func getBytes() []byte {
  var data []byte

  file := flag.String("file", "", "a file containing a valid block")
  isHex := flag.Bool("hex", false, "flag to indicate if data is in hex or raw bytes format")

  flag.Parse()

  if *file != "" {
    data, _ = ioutil.ReadFile(*file)
  } else {
    r := bufio.NewReader(os.Stdin)
    buf := make([]byte, 0, 4*1024)
    for {
      n, err := r.Read(buf[:cap(buf)])
      buf = buf[:n]
      if n == 0 {
        if err == nil {
          continue
        }
        if err == io.EOF {
          break
        }
        log.Fatal(err)
      }
    }
    data = buf
  }

  var output []byte
  var err error
  if *isHex == true {
    n := bytes.IndexByte(data, 0)
    s := string(data[:n])
    output, err = hex.DecodeString(s)
    if err != nil {
      panic(err)
    }
  }

  return output
}
