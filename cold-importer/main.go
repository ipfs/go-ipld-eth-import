package main

import (
	"flag"

	lib "github.com/ipfs/go-ipld-eth-import/lib"
)

/*
## EXAMPLE USAGE

make cold && ./build/bin/cold-importer \
	--geth-db-filepath /Users/hj/Documents/tmp/geth-data/geth/chaindata \
	--ipfs-repo-path ~/.ipfs \
	--block-number 0
*/

func main() {
	var (
		blockNumber  uint64
		ipfsRepoPath string
		dbFilePath   string
		syncMode     string
	)

	// Command line options
	flag.Uint64Var(&blockNumber, "block-number", 0, "Canonical number of the block state to import")
	flag.StringVar(&ipfsRepoPath, "ipfs-repo-path", "~/.ipfs", "IPFS repository path")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.StringVar(&syncMode, "sync-mode", "state", "What to synchronize")
	flag.Parse()

	// IPFS
	ipfs := lib.IpfsInit(ipfsRepoPath)

	// Cold Database
	db := lib.GethDBInit(dbFilePath)
	defer db.Stop()

	// Init the synchronization stack
	ts := lib.NewTrieStack(blockNumber)
	defer ts.Close()

	// Launch Synchronization
	ts.TraverseStateTrie(db, ipfs, syncMode, blockNumber)

	// Print the metrics
	printReport(syncMode)
}
