package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-ipld-eth-import/metrics"
)

type Walker struct {
	ipfs                  *IPFS
	dirPath               string
	prefix                string
	iterationCheapCounter int
}

// InitWalker gives us the Walker object, and set up the metrics
// of this exercise.
func InitWalker(ipfs *IPFS, dirPath, prefix string) *Walker {
	// Metrics in this operation
	metrics.NewLogger("traverse-directory")
	metrics.NewLogger("process-file")
	metrics.NewLogger("read-file")
	metrics.NewLogger("ipfs-dag-put")
	metrics.NewLogger("bytes-tranferred")

	return &Walker{
		ipfs:                  ipfs,
		dirPath:               dirPath,
		prefix:                prefix,
		iterationCheapCounter: 0,
	}
}

// TraverseDirectory is the main loop of this importer,
// it calls processFile as it goes encountering nodes.
func (w *Walker) TraverseDirectory() {
	_l := metrics.StartLogDiff("traverse-directory")

	// option --prefix makes the directory walk shorter.
	if w.prefix != "" {
		w.dirPath = filepath.Join(w.dirPath, w.prefix)
	}

	// Walk all files in directory
	filepath.Walk(w.dirPath, w.processFile)

	metrics.StopLogDiff("traverse-directory", _l)
}

// processFile is the core component of the traversal loop.
// If the node encountered is not a directory,
// it will get its data and import it into IPFS.
func (w *Walker) processFile(path string, info os.FileInfo, err error) error {
	_l := metrics.StartLogDiff("process-file")

	// Output a number to the user
	w.liveCounter()

	// Skip directories, of course
	if info.IsDir() {
		metrics.StopLogDiff("process-file", _l)
		return nil
	}

	// Get the file contents
	data := readFile(path)

	// And call `ipfs dag put`
	importIntoIPFS(w.ipfs, data)

	metrics.StopLogDiff("process-file", _l)
	return nil
}

// liveCounter gives the lonely user some company
func (w *Walker) liveCounter() {
	w.iterationCheapCounter += 1
	fmt.Printf("%d\r", w.iterationCheapCounter)
}

// readFile just calls ioutil.ReadFile and take metrics
func readFile(path string) []byte {
	_l := metrics.StartLogDiff("read-file")

	// Do it
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	metrics.StopLogDiff("read-file", _l)
	return data
}

// importIntoIPFS invokes our customized methods, leveraging
// the DAG.
func importIntoIPFS(ipfs *IPFS, rawData []byte) {
	_l := metrics.StartLogDiff("ipfs-dag-put")

	// Import it into IPFS,
	// with our stripped down functionality
	ipfs.DagPut(rawData, "importer-ipld-raw-data")

	metrics.AddLog("bytes-tranferred", int64(len(rawData)))
	metrics.StopLogDiff("ipfs-dag-put", _l)
}
