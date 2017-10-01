package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"

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
	account := db.TraverseToLeaf(stateRoot, stateTriePath)
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
	// TODO

	// STEP
	// We are done. Go grab some beverage of your preference.
	// TODO
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
func (g *gethDB) TraverseToLeaf(root []byte, path []byte) []byte {
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
