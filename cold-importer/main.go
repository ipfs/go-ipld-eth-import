package main

import (
	"flag"
	"fmt"

	metrics "github.com/hermanjunge/go-ipld-eth-import/metrics"
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

	n, sum, avg := metrics.GetAverageLogDiff("traverse-state-trie-iterations")
	fmt.Printf("Number of iterations:\t\t%d\n", n)

	fmt.Printf("\tBranches:\t\t%d\n", metrics.GetCounter("traverse-state-trie-branches"))
	fmt.Printf("\tExtensions:\t\t%d\n", metrics.GetCounter("traverse-state-trie-extensions"))
	fmt.Printf("\tLeaves:\t\t\t%d\n", metrics.GetCounter("traverse-state-trie-leaves"))

	fmt.Printf("==========================================\n\n")

	_, sum, _ = metrics.GetAverageLogDiff("traverse-state-trie")
	fmt.Printf("Time elapsed:\t\t\t%d \t\tms\n", sum/(1000*1000))

	n, sum, avg = metrics.GetAverageLogDiff("traverse-state-trie-iterations")
	fmt.Printf("Avg time per iteration:\t\t%.0f \t\tns\t(%d %d)\n", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("geth-leveldb-get-query")
	fmt.Printf("Avg time per levelDB:\t\t%.0f \t\tns\t(%d %d)\n", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("new-nodes-bytes-tranferred")
	fmt.Printf("Total bytes state:\t\t%d \tbytes\t(%.0f %d)\n", sum, avg, n)
}
