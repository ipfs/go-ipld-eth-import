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
	iterationCheapCounter int
}

func InitWalker(ipfs *IPFS) *Walker {
	return &Walker{
		ipfs: ipfs,
		iterationCheapCounter: 1,
	}
}

func (w *Walker) TraverseDirectory(ipfs *IPFS, dirPath string) {
	// Metrics in this operation
	metrics.NewLogger("traverse-directory")
	metrics.NewLogger("process-file")
	metrics.NewLogger("read-file")
	metrics.NewLogger("ipfs-dag-put")
	metrics.NewLogger("bytes-tranferred")

	// Walk all files in directory
	_l := metrics.StartLogDiff("traverse-directory")

	filepath.Walk(dirPath, w.processFile)

	metrics.StopLogDiff("traverse-directory", _l)
}

func (w *Walker) processFile(path string, info os.FileInfo, err error) error {
	var _l int

	// Live counter, to give the lonely user some company
	fmt.Printf("%d\r", w.iterationCheapCounter)
	w.iterationCheapCounter += 1

	_lpf := metrics.StartLogDiff("process-file")

	// Skip directories, of course
	if info.IsDir() {
		metrics.StopLogDiff("process-file", _lpf)
		return nil
	}

	// Read the file
	_l = metrics.StartLogDiff("read-file")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	metrics.StopLogDiff("read-file", _l)

	// Import it into IPFS
	_l = metrics.StartLogDiff("ipfs-dag-put")
	_ = w.ipfs.DagPut(data, "importer-ipld-raw-data")
	metrics.AddLog("bytes-tranferred", info.Size())
	metrics.StopLogDiff("ipfs-dag-put", _l)

	metrics.StopLogDiff("process-file", _lpf)

	return nil
}
