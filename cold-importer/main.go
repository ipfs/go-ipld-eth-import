package main

import (
	"flag"
	"fmt"
)

func main() {
	var blockNumber uint64

	// Command line options
	flag.Uint64Var(&blockNumber, "block-number", 0, "Canonical number of the block state to import")
	flag.Parse()

	// Cold Database
	db := InitStart()
	defer db.Stop()

	// Launch State Traversal
	ts := NewTrieStack(blockNumber)
	defer ts.Close()

	ts.TraverseStateTrie(db, blockNumber)

	// Print some stats
	printReport(ts)
}

func printReport(ts *trieStack) {
	fmt.Printf("Traversal finished\n==================\n\n")
	fmt.Printf("Time elapsed:\t\t%d ms\n", ts.Stats.TotalTime)
	fmt.Printf("Number of iterations:\t%d\n", ts.Stats.IterationsCnt)
	fmt.Printf("\tBranches:\t%d\n", ts.Stats.BranchCnt)
	fmt.Printf("\tExtensions:\t%d\n", ts.Stats.ExtensionCnt)
	fmt.Printf("\tLeaves:\t\t%d\n", ts.Stats.LeafCnt)
}
