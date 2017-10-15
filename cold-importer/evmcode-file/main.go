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
	--nibble 1a

*/

func main() {
	var (
		blockNumber uint64
		dbFilePath  string
		evmCodeDir  string
		nibble      string
	)

	// Command line options
	flag.Uint64Var(&blockNumber, "block-number", 0, "Canonical number of the block state to import")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.StringVar(&evmCodeDir, "evmcode-directory", "/tmp/evmcode", "Path to the directory to create the files to be dumped")
	flag.StringVar(&nibble, "nibble", "",
		"If set, selects one of the 16 branches of the state root. Only support one nibble {0,1,2,3,4,5,6,7,8,9,0,a,b,c,d,e,f}")
	flag.Parse()

	// Cold Database
	db := lib.GethDBInit(dbFilePath)
	defer db.Stop()

	// Init the synchronization stack
	ts := lib.NewTrieStack(db, blockNumber, evmCodeDir, nibble)
	defer ts.Close()

	// Launch Synchronization
	ts.TraverseStateTrie()

	// Print the metrics
	printReport()
}
