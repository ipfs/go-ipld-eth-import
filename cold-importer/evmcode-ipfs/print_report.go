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
	separatorFmt := "=========================================================================\n\n"
	loggersFmt := "%-25s: %12.0f ns  -> Total: %18d (%d)\n"

	// Actual Content
	fmt.Printf("Traversal finished\n")

	fmt.Printf(separatorFmt)

	// Logger Times
	n, sum, avg = metrics.GetAverageLogDiff("process-file")
	fmt.Printf(loggersFmt, "Avg time per iteration", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("read-file")
	fmt.Printf(loggersFmt, "Avg time per file read", avg, sum, n)

	n, sum, avg = metrics.GetAverageLogDiff("ipfs-dag-put")
	fmt.Printf(loggersFmt, "Avg time ipfs dag put", avg, sum, n)

	fmt.Printf(separatorFmt)

	// Totals
	_, sum, _ = metrics.GetAverageLogDiff("traverse-directory")
	fmt.Printf("%-25s: %12d ms\n", "Total Time elapsed", sum/(1000*1000))

	_, sum, avg = metrics.GetAverageLogDiff("bytes-tranferred")
	fmt.Printf("%-25s: %12d bytes\n", "Total bytes state", sum)
	fmt.Printf("%-25s: %12.0f bytes\n", "Average per iteration", avg)
}
