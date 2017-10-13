package lib

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ipfs/go-ipld-eth-import/metrics"
)

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
			if bytes.Compare(codeHash, emptyCodeHash) != 0 {
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
