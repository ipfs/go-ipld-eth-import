package lib

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ipfs/go-ipld-eth-import/metrics"
)

var emptyCodeHash = crypto.Keccak256(nil)

// This is the known root hash of an empty trie.
var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

// getTrieNodeChildren will decode the given RLP.
// If the result is a branch or extension, it will return its
// children hashes, otherwise, nil will be returned.
func getTrieNodeChildren(rlpTrieNode []byte) [][]byte {
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

// getTrieNodeEVMCode will decode the given RLP.
// If the result is a leaf, it will return its EVM Code.
// If the codehash is equal to the empty value, it will return nil.
func getTrieNodeEVMCode(rlpTrieNode []byte) []byte {
	var (
		out []byte
		i   []interface{}
	)

	// Decode the node
	err := rlp.DecodeBytes(rlpTrieNode, &i)
	if err != nil {
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
			out = nil
		case '\x02':
			fallthrough
		case '\x03':
			// This is a leaf
			var account []interface{}
			err = rlp.DecodeBytes(last, &account)
			if err != nil {
				panic(err)
			}

			codeHash := account[3].([]byte)
			if !bytes.Equal(codeHash, emptyCodeHash) {
				metrics.IncCounter("traverse-state-smart-contracts")
				out = codeHash
			} else {
				out = nil
			}

		default:
			// Zero tolerance
			panic("unknown hex prefix on trie node")
		}

	case 17:
		// This is a branch
		out = nil

	default:
		panic("unknown trie node type")
	}

	return out
}

// getTrieNodeStorageRoot will decode the given RLP.
// If the result is a leaf, it will return its Storage Trie.
// If the Storage Trie is equal to the empty trie, it will return nil.
func getTrieNodeStorageRoot(rlpTrieNode []byte) []byte {
	var (
		out []byte
		i   []interface{}
	)

	// Decode the node
	err := rlp.DecodeBytes(rlpTrieNode, &i)
	if err != nil {
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
			out = nil
		case '\x02':
			fallthrough
		case '\x03':
			// This is a leaf
			var account []interface{}
			err = rlp.DecodeBytes(last, &account)
			if err != nil {
				// Ignore the error, as we are probably dealing
				// with a storage trie leaf
				// TODO
				// Improve the stack we are using so we can know this
				// And procede better
				return nil
			}

			storageRoot := account[2].([]byte)
			if !bytes.Equal(storageRoot, emptyRoot[:]) {
				out = storageRoot
			} else {
				out = nil
			}

		default:
			// Zero tolerance
			panic("unknown hex prefix on trie node")
		}

	case 17:
		// This is a branch
		out = nil

	default:
		panic("unknown trie node type")
	}

	return out
}
