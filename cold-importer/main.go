package main

import (
	"fmt"
	"gx/ipfs/QmbBhyDKsY4mbY6xsKt3qu9Y7FPvMJ6qbD8AMjYYvPRw1g/goleveldb/leveldb"
	"gx/ipfs/QmbBhyDKsY4mbY6xsKt3qu9Y7FPvMJ6qbD8AMjYYvPRw1g/goleveldb/leveldb/errors"
)

func main() {
	file := "/Users/hj/.parity-data/chains/ethereum/db/906a34e69aec8c0d/overlayrecent/db"
	db, err := leveldb.OpenFile(file, nil)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		fmt.Println("Corrupta")
		db, err = leveldb.RecoverFile(file, nil)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", db)
}
