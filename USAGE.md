## USAGE

Bring Ethereum to IPFS.

## Cold Importer.

Set of tools that

* Grab the information from a disconnected (hence "_cold_") levelDB from go-ethereum to files.
* Traverse directories with these files to import them into IPFS.

By separating those functions, and allowing the use of prefixes, these activities can have a degree of scaling.

### Imported Information

* `evm-code`
  EVM Code (i.e. Smart Contracts)

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
* `eth-storage-trie`
  Storage of the accounts as nodes of its respective trie.

### Requirements

Just do

```
./build/first-load.sh
```

If you want to become a weekend contributor, here is your low hanging fruit:
_make this script elegant_.

### Importers

#### EVM Code from GethDB to File

##### Build

```
make evmcode-file
```

##### Example Usage

```
./build/bin/evmcode-file \
	--block-number 4339465 \
	--geth-db-filepath /Users/hj/Documents/data/fast-geth/geth/chaindata \
	--dump-directory /tmp/evmcode \
	--nibble 2
```

##### Command Line Parameters

* `--block-number`
  Specifies the block number data (canonical chain in this db) to fetch.

* `--geth-db-filepath`
  LevelDB Directory. As it only supports only one process, make sure it is
  not being used by go-ethereum or other program, hence, this importing is
  called _cold_.

* `--dump-directory`
  The directory where the `evmcode` files will be dumped.

* `--nibble`
  Supports just one nibble (hex character). If set, it will traverse the state
  trie down the chosen branch of the root, making your processing time about
  `15/16` faster.

#### EVM Code from File to IPFS

##### Build

```
make evvmcode-ipfs
```

##### Example Usage

Execute doing

```
./build/bin/evmcode-ipfs \
	--evmcode-directory /Users/hj/Documents/data/cold/evmcode
	--ipfs-repo-path ~/.ipfs \
	--prefix 1a
```

##### Command Line Parameters

* `--evmcode-directory`
  The directory where the `evmcode` files where dumped after processing the
  geth levelDB (using `evmcode-file`).

* `--ipfs-repo-path`
  The IPFS repository. Must be unlocked, i.e. `ipfs daemon` should not be using it.

* `--prefix`
  Useful to scale the effort: It will only process the files which name starts
  with the given prefix. It only support prefixes of two (2) characters (ex: `1a`).

#### State Trie Nodes from GethDB to File

##### Build

```
make state-trie-file
```

##### Example Usage

```
./build/bin/state-trie-file \
	--block-number 4371405 \
	--geth-db-filepath /Users/hj/Documents/data/fast-geth/geth/chaindata \
	--dump-directory /tmp/state-trie \
	--nibble 2
```

##### Command Line Parameters

* `--block-number`
  Specifies the block number data (canonical chain in this db) to fetch.

* `--geth-db-filepath`
  LevelDB Directory. As it only supports only one process, make sure it is
  not being used by go-ethereum or other program, hence, this importing is
  called _cold_.

* `--dump-directory`
  The directory where the `state trie node` files will be dumped.

* `--nibble`
  Supports just one nibble (hex character). If set, it will traverse the state
  trie down the chosen branch of the root, making your processing time about
  `15/16` faster.
