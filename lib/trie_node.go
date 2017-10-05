package lib

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	goque "github.com/beeker1121/goque"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	cid "github.com/ipfs/go-cid"
	metrics "github.com/ipfs/go-ipld-eth-import/metrics"
	mh "github.com/multiformats/go-multihash"
)

const MEthStateTrie = 0x96

// trieStack wraps the goque stack, enabling the adding of specific
// methods for dealing with the state trie.
type TrieStack struct {
	*goque.Stack
}

func NewTrieStack(blockNumber uint64) *TrieStack {
	var err error
	ts := &TrieStack{}

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
	metrics.NewLogger("ipfs-block-get-queries")
	metrics.NewLogger("ipfs-dag-put-queries")
	metrics.NewLogger("geth-leveldb-get-queries")
	metrics.NewLogger("trie-node-children-processes")

	metrics.NewLogger("traverse-state-trie-iterations")
	metrics.NewLogger("new-nodes-bytes-tranferred")

	metrics.NewCounter("traverse-state-trie-branches")
	metrics.NewCounter("traverse-state-trie-extensions")
	metrics.NewCounter("traverse-state-trie-leaves")

	return ts
}

// TraverseStateTrie, traverses the entire state trie of a given block number
// from a "cold" geth database
func (ts *TrieStack) TraverseStateTrie(db *GethDB, ipfs *IPFS, syncMode string, blockNumber uint64) {
	var (
		err error
		val []byte
	)

	metrics.StartLogDiff("traverse-state-trie")

	// From the block number, we get its canonical hash, and header RLP
	blockHash := db.GetCanonicalHash(blockNumber)
	headerRLP := db.GetHeaderRLP(blockHash, blockNumber)

	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(headerRLP), header); err != nil {
		panic(err)
	}

	// Init the traversal with the state root
	_, err = ts.Push(header.Root[:])
	if err != nil {
		panic(err)
	}

	_iterationsCnt := 1

	for {
		_tsti := metrics.StartLogDiff("traverse-state-trie-iterations")

		// Live counter, to give the lonely user some company
		fmt.Printf("%d\r", _iterationsCnt)
		_iterationsCnt += 1

		// Get the next item from the stack
		item, err := ts.Pop()
		if err == goque.ErrEmpty {
			break
		}
		if err != nil {
			panic(err)
		}

		switch syncMode {
		case "state":
			// Do not work twice.
			// Ask the good IPFS Blockstore if we have this key already
			if hasTrieNode(ipfs, item.Value) {
				metrics.StopLogDiff("traverse-state-trie-iterations", _tsti)
				continue
			}

			// We don't have it. We fetch it from geth and import it.
			val = fetchFromGethDB(db, item.Value)
			importToIPFS(ipfs, item.Value)
		case "evmcode":
			// Just fetch the value
			val = fetchFromGethDB(db, item.Value)
			// If it is a leaf, we import its EVM Code
			// TODO
		default:
			panic("unsupported sync option")
		}

		// Find the children of this element
		findChildrenToStack(ts, val)

		metrics.StopLogDiff("traverse-state-trie-iterations", _tsti)
	}

	metrics.StopLogDiff("traverse-state-trie", 0)
}

// hasTrieNode will query our IPFS blockstore, and tell us whether we have
// this key or not.
func hasTrieNode(ipfs *IPFS, key []byte) bool {
	var _l int

	// Create the cid
	mhash, err := mh.Encode(key, mh.KECCAK_256)
	if err != nil {
		panic(err)
	}
	c := cid.NewCidV1(MEthStateTrie, mhash)

	// Do we have this node imported already?
	_l = metrics.StartLogDiff("ipfs-block-get-queries")
	ipfsBlockFound := ipfs.HasBlock(c.String())
	metrics.StopLogDiff("ipfs-block-get-queries", _l)

	return ipfsBlockFound
}

func fetchFromGethDB(db *GethDB, key []byte) []byte {
	_l := metrics.StartLogDiff("geth-leveldb-get-queries")
	val, err := db.Get(key)
	if err != nil {
		panic(err)
	}
	metrics.StopLogDiff("geth-leveldb-get-queries", _l)
	metrics.AddLog("new-nodes-bytes-tranferred", int64(len(val)))

	return val
}

func importToIPFS(ipfs *IPFS, val []byte) {
	_l := metrics.StartLogDiff("ipfs-dag-put-queries")
	_ = ipfs.DagPut(val, "eth-state-trie")
	metrics.StopLogDiff("ipfs-dag-put-queries", _l)
}

func findChildrenToStack(ts *TrieStack, rawVal []byte) {
	_l := metrics.StartLogDiff("trie-node-children-processes")
	children := processTrieNode(rawVal)
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

// processTrieNode will decode the given RLP. If the result is a branch or
// extension, it will return its children hashes, otherwise, nil will
// be returned.
func processTrieNode(rlpTrieNode []byte) [][]byte {
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
