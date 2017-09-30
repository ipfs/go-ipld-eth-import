package main

import (
	"flag"
	"fmt"

	metrics "github.com/ipfs/go-ipld-eth-import/metrics"
)

func main() {
	var (
		blockNumber uint64
		ipfsCmd     string
		dbFilePath  string
	)

	// Command line options
	flag.Uint64Var(&blockNumber, "block-number", 0, "Canonical number of the block state to import")
	flag.StringVar(&ipfsCmd, "ipfs-cmd", "ipfs", "IPFS command path")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.Parse()

	// IPFS
	ipfs := ipfsInit(ipfsCmd)

	// Cold Database
	db := gethDBInitStart(dbFilePath)
	defer db.Stop()

	// Launch State Traversal
	ts := NewTrieStack(blockNumber)
	defer ts.Close()

	ts.TraverseStateTrie(db, ipfs, blockNumber)

	// Print some stats
	printReport()
}

func printReport() {
	// Formatters
	separatorFmt := "=========================================================================\n\n"
	iterationsFmt := "%-25s: %12d\n"
	loggersFmt := "%-25s: %12.0f ns  -> Total: %12d (%d)\n"

	// Actual Content
	fmt.Printf("Traversal finished\n")

	fmt.Printf(separatorFmt)

	// Iterations
	n, sum, avg := metrics.GetAverageLogDiff("traverse-state-trie-iterations")
	fmt.Printf(iterationsFmt, "Number of iterations", n)
	fmt.Printf(iterationsFmt, "  Branches", metrics.GetCounter("traverse-state-trie-branches"))
	fmt.Printf(iterationsFmt, "  Extensions", metrics.GetCounter("traverse-state-trie-extensions"))
	fmt.Printf(iterationsFmt, "  Leaves", metrics.GetCounter("traverse-state-trie-leaves"))

	fmt.Printf(separatorFmt)

	// Logger Times
	n, sum, avg = metrics.GetAverageLogDiff("traverse-state-trie-iterations")
	fmt.Printf(loggersFmt, "Avg time per iteration", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("ipfs-dag-get-queries")
	fmt.Printf(loggersFmt, "Avg time ipfs dag get", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("ipfs-dag-put-queries")
	fmt.Printf(loggersFmt, "Avg time ipfs dag put", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("geth-leveldb-get-queries")
	fmt.Printf(loggersFmt, "Avg time levelDB Get()", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("trie-node-processes")
	fmt.Printf(loggersFmt, "Avg time Node processing", avg, sum, n)

	fmt.Printf(separatorFmt)

	// Totals
	_, sum, _ = metrics.GetAverageLogDiff("traverse-state-trie")
	fmt.Printf("%-25s: %12d ms\n", "Total Time elapsed", sum/(1000*1000))

	_, sum, avg = metrics.GetAverageLogDiff("new-nodes-bytes-tranferred")
	fmt.Printf("%-25s: %12d bytes\n", "Total bytes state", sum)
	fmt.Printf("%-25s: %12.0f bytes\n", "Average per iteration", avg)
}
