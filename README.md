# go-ipld-eth-import

Bring Ethereum to IPFS.

## Cold Importer.

Grabs the information from a disconnected (hence "_cold_") levelDB from
go-ethereum and puts it into an IPFS client.

### Information

* `eth-block`
  Block header.
* `eth-tx`
  Transactions.
* `eth-tx-trie`
  Transactions as nodes of the transactions tries.
* `eth-state-trie`
  State elements as nodes of the state trie.
  Its leaves are the ethereum accounts.

### Usage

```
./build/bin/cold-importer --dbpath <dbpath> <options>
```

#### Command Line Parameters

* `--dbpath`
  LevelDB

* `--block-number`
  Specifies the block number data (canonical chain in this db) to fetch.
  Will default to the latest block if not specified
