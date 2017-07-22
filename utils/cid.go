package main

import (
	"encoding/hex"
	"fmt"

	cid "github.com/ipfs/go-cid"
	multihash "github.com/multiformats/go-multihash"
)

// This program just returns a cid against the input parameters
func main() {

	// TODO
	// Get the console arguments

	// TODO
	// I need the multihash
	// DEBUG
	buf, _ := hex.DecodeString("6263d74e77b2fdc85d359f95a04bec722ff91417154840f908e89652d202bdca")
	mHashBuf, _ := multihash.EncodeName(buf, "keccak-256")
	mHash, _ := multihash.Decode(mHashBuf)
	// DEBUG

	codecType := cid.Codecs["eth-block"]

	c := cid.NewCidV1(codecType, mHash)

	fmt.Printf("---> %v", c)
}
