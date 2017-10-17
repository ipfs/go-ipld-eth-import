package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	goque "github.com/beeker1121/goque"
	types "github.com/ethereum/go-ethereum/core/types"
	crypto "github.com/ethereum/go-ethereum/crypto"
	rlp "github.com/ethereum/go-ethereum/rlp"
	metrics "github.com/ipfs/go-ipld-eth-import/metrics"
)

// MEthStateTrie is the cid codec for a Ethereum State Trie.
const MEthStateTrie = 0x96

var emptyCodeHash = crypto.Keccak256(nil)

// TrieStack wraps the goque stack, enabling the adding of specific
// methods for dealing with the state trie.
type TrieStack struct {
	*goque.Stack

	db                    *GethDB
	dumpDir               string
	operation             string
	firstNibbleInt        int
	iterationCheapCounter int
}

// NewTrieStack initializes the traversal stack, and finds the canonical
// block header, returning the TrieStack wrapper for further instructions
func NewTrieStack(db *GethDB, blockNumber uint64, dumpDir, nibble, operation string) *TrieStack {
	var err error
	ts := &TrieStack{}

	// Metrics in this operation
	metrics.NewLogger("traverse-state-trie")
	metrics.NewLogger("geth-leveldb-get-queries")
	metrics.NewLogger("trie-node-children-processes")
	metrics.NewLogger("traverse-state-trie-iterations")
	metrics.NewLogger("new-nodes-bytes-tranferred")
	metrics.NewLogger("file-creations")
	metrics.NewCounter("traverse-state-trie-branches")
	metrics.NewCounter("traverse-state-trie-extensions")
	metrics.NewCounter("traverse-state-trie-leaves")
	metrics.NewCounter("traverse-state-smart-contracts")

	// Add the reference to the database
	ts.db = db

	// Hardcoded stack directory. Sue me
	dataDirectoryName := "/tmp/trie_stack_data_dir/" + strconv.FormatUint(blockNumber, 10)

	// Clearing the directory if exists, as we want to start always
	// with a fresh stack database.
	os.RemoveAll(dataDirectoryName)
	ts.Stack, err = goque.OpenStack(dataDirectoryName)
	if err != nil {
		panic(err)
	}

	// Find the block header RLP we need
	blockHash := db.GetCanonicalHash(blockNumber)
	headerRLP := db.GetHeaderRLP(blockHash, blockNumber)
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(headerRLP), header); err != nil {
		panic(err)
	}

	// Finally, Init the traversal with the state root
	_, err = ts.Push(header.Root[:])
	if err != nil {
		panic(err)
	}

	// Assign these variables
	ts.dumpDir = dumpDir

	switch operation {
	case "evmcode":
		ts.operation = "evmcode"
	case "state-trie":
		ts.operation = "state-trie"
	default:
		panic("operation not supported")
	}

	if len(nibble) > 1 {
		panic("unsupported nibble lenght")
	}
	if len(nibble) == 1 {
		n := nibble[0]
		switch {
		case n >= '0' && n <= '9':
			ts.firstNibbleInt = int(n - 48)
		case n >= 'a' && n <= 'f':
			ts.firstNibbleInt = int(n + 10 - 97)
		default:
			panic("wrong value for nibble")
		}
	} else {
		ts.firstNibbleInt = -1
	}

	// Return the wrapped object
	ts.iterationCheapCounter = 0
	return ts
}

// TraverseStateTrie performs a stack assisted traversal
// over the state trie node.
func (ts *TrieStack) TraverseStateTrie() {
	_l := metrics.StartLogDiff("traverse-state-trie")

	for {
		ts.liveCounter()
		err := ts.traverseStateTrieIteration()
		if err == goque.ErrEmpty {
			break
		}
		if err != nil {
			panic(err)
		}
	}

	metrics.StopLogDiff("traverse-state-trie", _l)
}

// traverseStateTrieIteration is the atomic component of the
// loop in TraverseStateTrie.
func (ts *TrieStack) traverseStateTrieIteration() error {
	_l := metrics.StartLogDiff("traverse-state-trie-iterations")

	// Get the next item from the stack
	item, err := ts.Pop()
	if err != nil {
		return err
	}
	// This clarifies a bit the code below
	key := item.Value

	// Fetch the value
	val := ts.fetchFromGethDB(key)

	switch ts.operation {
	case "evmcode":
		// If it is a leaf, we will get its EVM Code
		evmCodeKey := getTrieNodeEVMCode(val)
		if evmCodeKey != nil {
			code := ts.fetchFromGethDB(evmCodeKey)
			ts.storeFile(crypto.Keccak256(code), code)
		}
	case "state-trie":
		// Just store the found element
		ts.storeFile(key, val)
	}

	// Find the children of this element.
	// If found, they will be pushed in the stack.
	ts.findChildrenToStack(val)

	metrics.StopLogDiff("traverse-state-trie-iterations", _l)
	return nil
}

// liveCounter gives the lonely user some company
func (ts *TrieStack) liveCounter() {
	ts.iterationCheapCounter++
	fmt.Printf("%d\r", ts.iterationCheapCounter)
}

// fetchFromGethDB returns the value from the cold LevelDB.
func (ts *TrieStack) fetchFromGethDB(key []byte) []byte {
	_l := metrics.StartLogDiff("geth-leveldb-get-queries")

	val, err := ts.db.Get(key)
	if err != nil {
		panic(err)
	}
	metrics.AddLog("new-nodes-bytes-tranferred", int64(len(val)))

	metrics.StopLogDiff("geth-leveldb-get-queries", _l)
	return val
}

// findChildrenToStack evaluates a trie node. If it finds any
// children, it will add them to the stack, to follow the traversal.
func (ts *TrieStack) findChildrenToStack(rawVal []byte) {
	_l := metrics.StartLogDiff("trie-node-children-processes")

	children := getTrieNodeChildren(rawVal)
	if children != nil {
		for idx, child := range children {
			// If we are in the first iteration (i.e. the root),
			// we see whether --nibble is set. If so, process only the given one.
			if ts.iterationCheapCounter == 1 && ts.firstNibbleInt != -1 {
				if idx != ts.firstNibbleInt {
					continue
				}
				// Tell the user what's going on
				fmt.Printf("Reduced traversing from the root, down to %d\n",
					ts.firstNibbleInt)
			}

			_, err := ts.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	metrics.StopLogDiff("trie-node-children-processes", _l)
}

// storeFile will take the trie node contents, and store them into
// the file system, with the given key as a file name.
// It will take the first three bytes as subdirectories,
// to make its lookup easier.
func (ts *TrieStack) storeFile(key, contents []byte) {
	_l := metrics.StartLogDiff("file-creations")

	fileName := fmt.Sprintf("%x", key)
	fileDir := filepath.Join(ts.dumpDir, fileName[0:2], fileName[2:4], fileName[4:6])
	err := os.MkdirAll(fileDir, 0755)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath.Join(fileDir, fileName), contents, 0644)
	if err != nil {
		panic(err)
	}

	metrics.StopLogDiff("file-creations", _l)
}
