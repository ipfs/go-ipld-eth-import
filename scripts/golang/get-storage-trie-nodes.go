package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/beeker1121/goque"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

/*
## EXAMPLE USAGE

go build -o /tmp/getnodes scripts/golang/get-storage-trie-nodes.go &&
/tmp/getnodes \
	--geth-db-filepath /Users/hj/Documents/tmp/fast-geth-data/geth/chaindata \
	--block-number 4321849 \
	--eth-address 6810e776880c02933d47db1b9fc05908e5386b96

*/

// Pretty Printing
const hashFmt = "%-20s: %32x\n"

type EthAccount struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte // This is the storage root trie
	CodeHash []byte // This is the hash of the EVM code
}

func main() {
	var (
		ethAddress, dbFilePath string
		blockNumber            uint64
	)

	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.Uint64Var(&blockNumber, "block-number", 4321849, "Canonical chain block number")
	flag.StringVar(&ethAddress, "eth-address", "6810e776880c02933d47db1b9fc05908e5386b96", "Address which keccak-256 hash we want")
	flag.Parse()

	// Initialization of the DB
	db := InitStartDB(dbFilePath)
	defer db.Stop()

	fmt.Printf("======================================================================================\n")

	// STEP
	// Get the block canonical hash from the block number given.
	blockHash := db.GetCanonicalHash(blockNumber)
	fmt.Printf(hashFmt, "BlockHash", blockHash)

	// STEP
	// Fetch the block by its hash, decode it, get its state root.
	headerRLP := db.GetHeaderRLP(blockHash, blockNumber)
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(headerRLP), header); err != nil {
		panic(err)
	}
	stateRoot := header.Root[:]
	fmt.Printf(hashFmt, "State Root", stateRoot)

	// Additional STEP
	// Save this header RLP into a file
	dirPath := fmt.Sprintf("/tmp/get-trie-nodes/%d/0x%s", blockNumber, ethAddress)
	filePath := fmt.Sprintf("%s/eth-block-header-rlp-%d", dirPath, blockNumber)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filePath, headerRLP, 0644)
	if err != nil {
		panic(err)
	}

	// STEP
	// Compute the keccak256 hash of the ETH Address, to get the path to that data.
	kb, err := hex.DecodeString(ethAddress)
	if err != nil {
		panic(err)
	}
	stateTriePath := crypto.Keccak256(kb)
	fmt.Printf(hashFmt, "State Trie Path", stateTriePath)

	// STEP
	// Traverse the state trie from the root, using the path above.
	account := db.TraverseToLeaf(stateRoot, stateTriePath, dirPath)
	fmt.Printf(hashFmt, "Account RLP", account[:32])

	// STEP
	// The leaf should be decoded, get its storage trie.
	var decodedAccount EthAccount
	err = rlp.DecodeBytes(account, &decodedAccount)
	if err != nil {
		panic(err)
	}
	fmt.Printf(hashFmt, "Storage Trie Root", decodedAccount.Root)

	// STEP
	// We full traverse the whole trie, storing each node as a file, in /tmp.
	db.FullTraverseStorageTrie(decodedAccount.Root, dirPath)

	// STEP
	// We are done. Go grab some beverage of your preference.
	fmt.Printf("\nWe are done, check up your files at %s\n", dirPath)
}

/*
  GETH DB HELPERS
*/

type gethDB struct {
	db *leveldb.DB
}

func InitStartDB(path string) *gethDB {
	if path == "" {
		panic("Path to the Geth's DB must be specified (--geth-db-filepath option)")
	}
	db, err := leveldb.OpenFile(path, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		fmt.Println("Corrupt")
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		panic(err)
	}

	return &gethDB{db: db}
}

func (g *gethDB) Stop() {
	g.db.Close()
}

func (g *gethDB) Get(key []byte) ([]byte, error) {
	return g.db.Get(key, nil)
}

func (g *gethDB) GetCanonicalHash(number uint64) []byte {
	headerPrefix := []byte("h")
	numSuffix := []byte("n")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), numSuffix...)
	val, _ := g.db.Get(key, nil)

	return val
}

func (g *gethDB) GetHeaderRLP(hash []byte, number uint64) []byte {
	headerPrefix := []byte("h")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), hash...)

	val, _ := g.db.Get(key, nil)
	return val
}

/*
  TRIE TRAVERSAL HELPERS
*/
func (g *gethDB) TraverseToLeaf(root []byte, path []byte, dirPath string) []byte {
	var (
		out []byte
		i   []interface{}
	)

	// convert the whole path to nibbles
	var nibbledPath []byte
	for _, p := range path {
		nibbledPath = append(nibbledPath, p/16)
		nibbledPath = append(nibbledPath, p%16)
	}
	path = nibbledPath

	// Just a separator
	fmt.Printf("\n")

	for {
		if out != nil {
			break
		}

		raw, err := g.db.Get(root, nil)
		if err != nil {
			panic(err)
		}

		// Let's save these found nodes into files
		filePath := fmt.Sprintf("%s/eth-state-trie-rlp-%x", dirPath, root[:3])
		err = ioutil.WriteFile(filePath, raw, 0644)
		if err != nil {
			panic(err)
		}

		err = rlp.DecodeBytes(raw, &i)
		if err != nil {
			panic(err)
		}

		// What is it?
		switch len(i) {
		case 2:
			// Extension or leaf
			first := i[0].([]byte)
			last := i[1].([]byte)

			var nibbledFirst []byte
			for _, f := range first {
				nibbledFirst = append(nibbledFirst, f/16)
				nibbledFirst = append(nibbledFirst, f%16)
			}
			first = nibbledFirst

			switch first[0] {
			// 0 and 1 are extensions
			case '\x00':
				fallthrough
			case '\x01':
				root = last
				path = path[len(first):]

				fmt.Printf(hashFmt, "Got Hash", root)
				fmt.Printf(hashFmt, "  Consumed Path", first)

			// 2 and 3 are leaves
			// Just return the last element
			case '\x02':
				fallthrough
			case '\x03':
				out = last

				fmt.Printf(hashFmt, "Got RLP", last[:32])
			default:
				panic("unknown hex prefix on trie node")

			}
		case 17:
			// We got a Branch
			// Pick the first nibble from the path
			nibble := path[0]

			// Follow the hash mapped by that nibble
			root = i[int(nibble)].([]byte)
			path = path[1:]

			fmt.Printf(hashFmt, "Got Hash", root)
			fmt.Printf(hashFmt, "  Consumed Path", nibble)
		}
	}

	fmt.Printf("\n")
	return out
}

// getHexIndex returns to you the integer 0 - 15 equivalent to your
// string character if applicable, or -1 otherwise.
func getHexIndex(s string) int {
	if len(s) != 1 {
		return -1
	}

	c := byte(s[0])
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	}

	return -1
}

func (g *gethDB) FullTraverseStorageTrie(root []byte, dirPath string) {
	// Just a separator
	fmt.Printf("\n")

	stackDirectoryName := "/tmp/trie_stack_data_dir/" + dirPath[len(dirPath)-30:]

	// Clearing the directory if exists, as we want to start always
	// with a fresh stack database.
	os.RemoveAll(stackDirectoryName)
	stack, err := goque.OpenStack(stackDirectoryName)
	if err != nil {
		panic(err)
	}
	defer stack.Close()

	// Init the traversal with the state root
	_, err = stack.Push(root[:])
	if err != nil {
		panic(err)
	}

	storageTrieNodeCounter := 0

	for {
		storageTrieNodeCounter += 1
		fmt.Printf("Processing Storage Trie Node %d\r", storageTrieNodeCounter)

		// Get the next item from the stack
		item, err := stack.Pop()
		if err == goque.ErrEmpty {
			break
		}
		if err != nil {
			panic(err)
		}

		// For easy reading
		key := item.Value

		// Search the Geth LevelDB for the value
		val, err := g.db.Get(key, nil)
		if err != nil {
			panic(err)
		}

		// Create a file
		filePath := fmt.Sprintf("%s/eth-storage-trie-rlp-%x", dirPath, key[:3])
		err = ioutil.WriteFile(filePath, val, 0644)
		if err != nil {
			panic(err)
		}

		// Process the raw node for children
		children := processTrieNode(val)
		if children != nil {
			for _, child := range children {
				_, err = stack.Push(child)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	fmt.Printf("Processed %d Storage Trie Node(s)\n", storageTrieNodeCounter)
}

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
			out = [][]byte{last}
		case '\x02':
			fallthrough
		case '\x03':
			// This is a leaf
			out = nil
		default:
			// Zero tolerance
			panic("unknown hex prefix on trie node")
		}

	case 17:
		// This is a branch
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
