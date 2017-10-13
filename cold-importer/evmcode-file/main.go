package main

import (
	"flag"

	"github.com/ipfs/go-ipld-eth-import/lib"
)

/*

## EVM CODE to FILE

Traverses the entire state of a given block, finding its accounts.
Whenever it finds an account with a non empty hash (i.e. a smart contract),
it fetches its contents, dumping them in a file, with its keccak256 hash
as a name.

## EXAMPLE USAGE

make evmcode-file && \
./build/bin/evmcode-file \
	--block-number 4352702 \
	--geth-db-filepath /Users/hj/Documents/data/fast-geth/geth/chaindata \
	--evmcode-directory /tmp/evmcode \
	--prefix 1a

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
	ts := lib.NewTrieStack(db, blockNumber) // Aca colocamos el prefijo y filepath despues
	defer ts.Close()

	/*
		// Launch Synchronization
		ts.TraverseStateTrie()

		// Print the metrics
		printReport()
	*/
}
