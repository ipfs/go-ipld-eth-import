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

Install the following programs

```
go get -v -u github.com/whyrusleeping/gx
go get -v -u github.com/whyrusleeping/gx-go

go get -v -u github.com/ipfs/go-ipld-eth
```

### Importers

#### EVM Code from GethDB to File

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
