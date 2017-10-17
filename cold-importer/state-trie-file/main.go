package main

import (
	"flag"

	"github.com/ipfs/go-ipld-eth-import/lib"
)

/*

## STATE TRIE NODES to FILE

Traverses the entire state trie of a given block, storing
the find nodes into files.

## EXAMPLE USAGE

./build/bin/state-trie-file \
	--block-number 4371405 \
	--geth-db-filepath /Users/hj/Documents/data/fast-geth/geth/chaindata \
	--dump-directory /tmp/state-trie \
	--nibble 2

*/

func main() {
	var (
		blockNumber uint64
		dbFilePath  string
		dumpDir     string
		nibble      string
	)

	// Command line options
	flag.Uint64Var(&blockNumber, "block-number", 0, "Canonical number of the block state to import")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.StringVar(&dumpDir, "dump-directory", "/tmp/state-trie", "Path to the directory to create the files to be dumped")
	flag.StringVar(&nibble, "nibble", "",
		"If set, selects one of the 16 branches of the state root. Only support one nibble {0,1,2,3,4,5,6,7,8,9,0,a,b,c,d,e,f}")
	flag.Parse()

	// Cold Database
	db := lib.GethDBInit(dbFilePath)
	defer db.Stop()

	// Init the synchronization stack
	ts := lib.NewTrieStack(db, blockNumber, dumpDir, nibble, "state-trie")
	defer ts.Close()

	// Launch Synchronization
	ts.TraverseStateTrie()

	// Print the metrics
	printReport()
}
