## USAGE

Bring Ethereum to IPFS.

## Cold Importer.

Grabs the information from a disconnected (hence "_cold_") levelDB from
go-ethereum and puts it into an IPFS client.

### Imported Information

* `eth-state-trie`
  State elements as nodes of the state trie.
  Its leaves are the ethereum accounts.

#### Want

* `eth-block`
  Block header.
* `eth-tx`
  Transactions.
* `eth-tx-trie`
  Transactions as nodes of the transactions tries.
* `evm-code`
  EVM Code (i.e. Smart Contracts)
* `eth-storage-trie`
  Storage of the accounts as nodes of its respective trie.

### Requirements

Install the following programs

```
	go get -v -u github.com/whyrusleeping/gx
	go get -v -u github.com/whyrusleeping/gx-go

	go get -v -u github.com/ipfs/go-ipld-eth
```

### Build

Compile with `make cold`.


### Example Usage

Execute doing

```
./build/bin/cold-importer \
	--geth-db-filepath /Users/hj/Documents/tmp/geth-data/geth/chaindata \
	--ipfs-repo-path ~/.ipfs \
	--block-number 0
```

#### Command Line Parameters

* `--geth-db-filepath`
  LevelDB Directory. As it only supports only one process, make sure it is
  not being used by go-ethereum or other program, hence, this importing is
  called _cold_.

* `--ipfs-repo-path`
  This is where you keep your IPFS files. Make sure these are not being used
  by an IPFS client.

* `--block-number`
  Specifies the block number data (canonical chain in this db) to fetch.
