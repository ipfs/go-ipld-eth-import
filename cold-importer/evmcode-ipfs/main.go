package main

import (
	"flag"

	"github.com/ipfs/go-ipld-eth-import/lib"
)

/*

## EXAMPLE USAGE

make evmcode-ipfs && ./build/bin/evmcode-ipfs --ipfs-repo-path ~/.ipfs

*/

func main() {
	var (
		ipfsRepoPath string
		evmcodeDir   string
	)

	// Command line options
	flag.StringVar(&ipfsRepoPath, "ipfs-repo-path", "~/.ipfs", "IPFS repository path")
	flag.StringVar(&evmcodeDir, "evmcode-directory", "/tmp/evmcode", "Directory where the EVM code files are")
	flag.Parse()

	// IPFS
	ipfs := lib.IpfsInit(ipfsRepoPath)

	// Launch the main loop
	walker := lib.InitWalker(ipfs)
	walker.TraverseDirectory(ipfs, evmcodeDir)

	// Print the metrics
	printReport()
}
