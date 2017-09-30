package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	var ethAddress string

	flag.StringVar(&ethAddress, "eth-address", "5abfec25f74cd88437631a7731906932776356f9", "Address which keccak-256 hash we want")
	flag.Parse()

	kb, err := hex.DecodeString(ethAddress)
	if err != nil {
		panic(err)
	}
	secureKey := crypto.Keccak256(kb)

	fmt.Printf("%x\n", secureKey)
}
