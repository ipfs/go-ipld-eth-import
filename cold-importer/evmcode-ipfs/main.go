package main

import (
	"flag"

	"github.com/ipfs/go-ipld-eth-import/lib"
)

/*

## EVM CODE IPFS

Takes the files dumped from the geth database and imports them to IPFS

## EXAMPLE USAGE

make evmcode-ipfs && ./build/bin/evmcode-ipfs --ipfs-repo-path ~/.ipfs

*/

func main() {
	var (
		evmcodeDir   string
		ipfsRepoPath string
		prefix       string
	)

	// Command line options
	flag.StringVar(&evmcodeDir, "evmcode-directory", "/tmp/evmcode", "Directory where the EVM code files are")
	flag.StringVar(&ipfsRepoPath, "ipfs-repo-path", "~/.ipfs", "IPFS repository path")
	flag.StringVar(&prefix, "prefix", "", "If set, will only process the files which name starts with <prefix>")
	flag.Parse()

	// IPFS
	ipfs := lib.InitIPFSNode(ipfsRepoPath)

	// Launch the main loop
	walker := lib.InitWalker(ipfs)
	walker.TraverseDirectory(ipfs, evmcodeDir)

	// Print the metrics
	printReport()
}
