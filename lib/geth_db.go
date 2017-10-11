package lib

import (
	"encoding/binary"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type GethDB struct {
	db *leveldb.DB
}

// GethDBInit creates the connection with the "cold" Geth LevelDB.
func GethDBInit(path string) *GethDB {
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

	return &GethDB{db: db}
}

// Stop Closes the DB
func (g *GethDB) Stop() {
	g.db.Close()
}

// Get returns the value associated to that key in the DB
func (g *GethDB) Get(key []byte) ([]byte, error) {
	return g.db.Get(key, nil)
}

// GetCanonicalHash returns the stored CHT Hash for a given number
func (g *GethDB) GetCanonicalHash(number uint64) []byte {
	headerPrefix := []byte("h")
	numSuffix := []byte("n")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), numSuffix...)
	val, _ := g.db.Get(key, nil)

	return val
}

// GetHeaderRLP returns the RLP of the block header
// for a pair (hash, number) as key
func (g *GethDB) GetHeaderRLP(hash []byte, number uint64) []byte {
	headerPrefix := []byte("h")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), hash...)

	val, _ := g.db.Get(key, nil)
	return val
}

// GetBodyRLP returns the RLP of the block header, plus ommer list
// transactions for a pair (hash, number) as key
func (g *GethDB) GetBodyRLP(hash []byte, number uint64) []byte {
	bodyPrefix := []byte("b")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(bodyPrefix, encodedNumber...), hash...)

	val, _ := g.db.Get(key, nil)
	return val
}
