package main

import (
  "encoding/hex"
  "fmt"
  "flag"

  cid "github.com/ipfs/go-cid"
  multihash "github.com/multiformats/go-multihash"
)

// This program just returns a keccak-256 cid against the input parameters
func main() {

  // TODO
  // Get the console arguments
  cidString := flag.String("cid", "", "cid to decode")
  flag.Parse()

  if *cidString == "" {
    panic("cid is required")
  }


  // TODO
  // Should use a parameter
  buf, _ := hex.DecodeString(*cidString)
  mHashBuf, _ := multihash.EncodeName(buf, "keccak-256")

  // TODO
  // Should use a parameter
  codecType := cid.Codecs["eth-block"]

  c := cid.NewCidV1(codecType, mHashBuf)

  fmt.Printf("%s\n", c)
}
