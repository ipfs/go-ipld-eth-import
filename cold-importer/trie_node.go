package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	goque "github.com/beeker1121/goque"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	metrics "github.com/hermanjunge/go-ipld-eth-import/metrics"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const MEthStateTrie = 0x96

// trieStack wraps the goque stack, enabling the adding of specific
// methods for dealing with the state trie.
type trieStack struct {
	*goque.Stack
}

func NewTrieStack(blockNumber uint64) *trieStack {
	var err error
	ts := &trieStack{}

	dataDirectoryName := "/tmp/trie_stack_data_dir/" + strconv.FormatUint(blockNumber, 10)

	// Clearing the directory if exists, as we want to start always
	// with a fresh stack database.
	os.RemoveAll(dataDirectoryName)
	ts.Stack, err = goque.OpenStack(dataDirectoryName)
	if err != nil {
		panic(err)
	}

	// Metrics in this operation
	metrics.NewLogger("traverse-state-trie")
	metrics.NewLogger("geth-leveldb-get-query")
	metrics.NewLogger("traverse-state-trie-iterations")
	metrics.NewLogger("new-nodes-bytes-tranferred")

	metrics.NewCounter("traverse-state-trie-branches")
	metrics.NewCounter("traverse-state-trie-extensions")
	metrics.NewCounter("traverse-state-trie-leaves")

	return ts
}

// TraverseStateTrie, traverses the entire state trie of a given block number
// from a "cold" geth database
func (ts *trieStack) TraverseStateTrie(db *gethDB, blockNumber uint64) {
	var err error

	metrics.StartLogDiff("traverse-state-trie")

	// From the block number, we get its canonical hash, and header RLP
	blockHash := db.GetCanonicalHash(blockNumber)
	headerRLP := db.GetHeaderRLP(blockHash, blockNumber)

	header := new(types.Header)
	if err = rlp.Decode(bytes.NewReader(headerRLP), header); err != nil {
		panic(err)
	}

	// Init the traversal with the state root
	_, err = ts.Push(header.Root[:])
	if err != nil {
		panic(err)
	}

	for {
		_tsti := metrics.StartLogDiff("traverse-state-trie-iterations")

		// Get the next item from the stack
		item, err := ts.Pop()
		if err == goque.ErrEmpty {
			break
		}
		if err != nil {
			panic(err)
		}

		// For clarity purposes
		key := item.Value

		// Create the cid
		mhash, err := mh.Encode(key, mh.KECCAK_256)
		if err != nil {
			panic(err)
		}
		c := cid.NewCidV1(MEthStateTrie, mhash)

		// DEBUG
		_ = c
		// DEBUG

		// TODO
		// Find out whether we already have this data imported in our local IPFS
		// * If we have it, it means that this branch has been already traversed,
		//   meaning that we need to `continue`
		// * If we don't let's add this data, before continuing the traversing.
		// NOTE
		// Pass the ipfs functions from outside this library,
		// to ensure some modularity.

		// Let's get that data
		_l := metrics.StartLogDiff("geth-leveldb-get-query")
		val, err := db.Get(key)
		if err != nil {
			panic(err)
		}
		metrics.StopLogDiff("geth-leveldb-get-query", _l)
		metrics.AddLog("new-nodes-bytes-tranferred", int64(len(val)))

		// Process this element
		// If it is a branch or an extension, add their children to the stack
		children := ts.processTrieNode(val)
		if children != nil {
			for _, child := range children {
				_, err = ts.Push(child)
				if err != nil {
					panic(err)
				}
			}
		}

		metrics.StopLogDiff("traverse-state-trie-iterations", _tsti)
	}

	metrics.StopLogDiff("traverse-state-trie", 0)
}

// processTrieNode will decode the given RLP. If the result is a branch or
// extension, it will return its children hashes, otherwise, nil will
// be returned.
func (ts *trieStack) processTrieNode(rlpTrieNode []byte) [][]byte {
	var (
		out [][]byte
		i   []interface{}
	)

	// Decode the node
	err := rlp.DecodeBytes(rlpTrieNode, &i)
	if err != nil {
		// Zero tolerance, if we have an err here,
		// it means our source database could be in bad shape.
		panic(err)
	}

	switch len(i) {
	case 2:
		first := i[0].([]byte)
		last := i[1].([]byte)

		switch first[0] / 16 {
		case '\x00':
			fallthrough
		case '\x01':
			// This is an extension
			metrics.IncCounter("traverse-state-trie-extensions")
			out = [][]byte{last}
		case '\x02':
			fallthrough
		case '\x03':
			// This is a leaf
			metrics.IncCounter("traverse-state-trie-leaves")
			out = nil
		default:
			// Zero tolerance
			panic("unknown hex prefix on trie node")
		}

	case 17:
		// This is a branch
		metrics.IncCounter("traverse-state-trie-branches")

		for _, vi := range i {
			v := vi.([]byte)
			switch len(v) {
			case 0:
				continue
			case 32:
				out = append(out, v)
			default:
				panic(fmt.Sprintf("unrecognized object: %v", v))
			}
		}

	default:
		panic("unknown trie node type")
	}

	return out
}
