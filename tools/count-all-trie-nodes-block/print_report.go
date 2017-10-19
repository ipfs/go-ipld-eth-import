package main

import (
	"fmt"

	"github.com/ipfs/go-ipld-eth-import/metrics"
)

func printReport() {
	var (
		n   int
		sum int64
		avg float64
	)

	// Formatters
	separatorFmt := "=========================================================================\n"
	iterationsFmt := "%-25s: %12d\n"
	loggersFmt := "%-25s: %12.0f ns  -> Total: %18d (%d)\n"

	// Actual Content
	fmt.Printf("Traversal finished\n")

	fmt.Println(separatorFmt)

	// Iterations
	// Count per kind of trie node
	// Count of smart contracts
	n, _, _ = metrics.GetAverageLogDiff("traverse-state-trie-iterations")
	fmt.Printf(iterationsFmt, "Number of iterations", n)
	fmt.Printf(iterationsFmt, "  Branches", metrics.GetCounter("traverse-state-trie-branches"))
	fmt.Printf(iterationsFmt, "  Extensions", metrics.GetCounter("traverse-state-trie-extensions"))
	fmt.Printf(iterationsFmt, "  Leaves", metrics.GetCounter("traverse-state-trie-leaves"))
	fmt.Printf(iterationsFmt, "  Smart Contracts", metrics.GetCounter("traverse-state-smart-contracts"))

	fmt.Println(separatorFmt)

	// Logger Times (quantity, average, sum)
	n, sum, avg = metrics.GetAverageLogDiff("traverse-state-trie-iterations")
	fmt.Printf(loggersFmt, "Avg time per iteration", avg, sum, n)

	fmt.Println(separatorFmt)

	// Totals
	_, sum, _ = metrics.GetAverageLogDiff("traverse-state-trie")
	fmt.Printf("%-25s: %12d ms\n", "Total Time elapsed", sum/(1000*1000))

	_, sum, avg = metrics.GetAverageLogDiff("new-nodes-bytes-tranferred")
	fmt.Printf("%-25s: %12d bytes\n", "Total bytes", sum)
	fmt.Printf("%-25s: %12.0f bytes\n", "Average per iteration", avg)
}
