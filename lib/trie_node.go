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

const MEthStateTrie = 0x96

var emptyCodeHash = crypto.Keccak256(nil)

// TrieStack wraps the goque stack, enabling the adding of specific
// methods for dealing with the state trie.
type TrieStack struct {
	*goque.Stack

	db                    *GethDB
	iterationCheapCounter int
}

// NewTriwStack initializes the traversal stack, and finds the canonical
// block header, returning the TrieStack wrapper for further instructions
func NewTrieStack(db *GethDB, blockNumber uint64) *TrieStack {
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

	// Return the wrapped object
	ts.iterationCheapCounter = 0
	return ts
}

// TODO: Document
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

// TODO Document
func (ts *TrieStack) traverseStateTrieIteration() error {
	_l := metrics.StartLogDiff("traverse-state-trie-iterations")

	// Get the next item from the stack
	item, err := ts.Pop()
	if err != nil {
		return err
	}

	// Fetch the value
	val := ts.fetchFromGethDB(item.Value)

	// If it is a leaf, we will get its EVM Code
	evmCodeKey := getTrieNodeEVMCode(val)
	if evmCodeKey != nil {
		code := ts.fetchFromGethDB(evmCodeKey)
		storeFile("evmcode", item.Value, code)
	}

	// Find the children of this element.
	// If found, they will be pushed in the stack.
	findChildrenToStack(ts, val)

	metrics.StopLogDiff("traverse-state-trie-iterations", _l)
	return nil
}

// liveCounter gives the lonely user some company
func (t *TrieStack) liveCounter() {
	t.iterationCheapCounter += 1
	fmt.Printf("%d\r", t.iterationCheapCounter)
}

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

func findChildrenToStack(ts *TrieStack, rawVal []byte) {
	_l := metrics.StartLogDiff("trie-node-children-processes")

	children := getTrieNodeChildren(rawVal)
	if children != nil {
		for _, child := range children {
			_, err := ts.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	metrics.StopLogDiff("trie-node-children-processes", _l)
}

func storeFile(kind string, key, contents []byte) {
	_l := metrics.StartLogDiff("file-creations")

	fileName := fmt.Sprintf("%x", key)
	fileDir := filepath.Join("/tmp/cold-dump", kind, fileName[0:2], fileName[2:4], fileName[4:6])
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
