package main

import (
	"flag"

	"github.com/ipfs/go-ipld-eth-import/lib"
)

/*

## COUNT ALL TRIE NODES IN A BLOCK

Starts from the state root of a block, and counts all the leaves,
extensions and branches it finds.

## BUILDING IT

build/convert-ipfs-deps.sh
go build -v -o build/bin/tool-count-all tools/count-all-trie-nodes-block/*.go
build/un-convert-ipfs-deps.sh

## EXAMPLE USAGE

./build/bin/tool-count-all \
	--block-number 4352702 \
	--geth-db-filepath /Users/hj/Documents/data/fast-geth/geth/chaindata

*/

func main() {
	var (
		blockNumber uint64
		dbFilePath  string
	)

	// Command line options
	flag.Uint64Var(&blockNumber, "block-number", 0, "Canonical number of the block state to import")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.Parse()

	// Cold Database
	db := lib.GethDBInit(dbFilePath)
	defer db.Stop()

	// Init the synchronization stack
	ts := lib.NewTrieStack(db, blockNumber, "", "", "count-all")
	defer ts.Close()

	// Launch Synchronization
	ts.TraverseStateTrie()

	// Print the metrics
	printReport()
}
