# go-ipld-eth-dump-star

`ETHEREUM CLIENT` ---> `IPFS`

Bag of hacks to get your ethereum information to IPFS.

## Overview

The best way to face this problem in an organically way, is to create simple
scripts that perform a simple task well (_sounds familiar?_). And then, use scripts to automatize
the tasks and control the flow.

There are three categories of programs in this repository:

* Extractors: Take data _from_ and ethereum source (mostly the Parity IPFS or RPC APIs).
* Feeders: Insert data _into_ IPFS, using some invocation of `ipfs dag put`. You won't find a lot of scripts here, as we will rely on the IPFS CLI.
* Control: Scripts to automate these processes ("_Get me the txs of the latests 100 blocks_").

## Requirements

* You will need to compile the `go` programs in the `utils` directory with `make`.

### Components

#### Extractors

##### `get-block-header-by-number <number in dec>`

Returns the block header in RLP.

#### Feeders

(TODO)

#### Control

##### `get-block-header-rlp-to-file <from> <to>`

Runs `get-block-header-by-number` around the given interval, dumping the
found contents into files `block-number-<number>`.
