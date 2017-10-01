package main

import (
  "bytes"
  "flag"
  "fmt"
  "os/exec"
  "time"
)

func main() {
  dag := flag.String("dag", "", "DAG id to lookup")
  flag.Parse()

  if *dag == "" {
    panic("Missing DAG")
  }

  DagGet(*dag)
}

func DagGet(cidString string) {

  cmd := exec.Command("ipfs", "dag", "get", "--local", cidString)

  var out bytes.Buffer
  cmd.Stdout = &out

  var stderr bytes.Buffer
  cmd.Stderr = &stderr

  // Here!
  start := time.Now().UnixNano()
  err := cmd.Run()
  fmt.Printf("Time Diff: %.0d\n", time.Now().UnixNano()-start)

  if err != nil {
    fmt.Printf("%s\n", stderr.String())
  } else {
    fmt.Printf("%s\n", out.String())
  }
}
