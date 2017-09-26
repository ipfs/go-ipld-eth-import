#!/bin/bash

################################################################################
#
# The reason we are doing this here and not in the Dockerfile is that
# you may want to update your packages without wanting to build the image
# everytime.
#
################################################################################

# IPFS + Plugin
go get -v github.com/ethereum/go-ethereum
go get -v github.com/ipfs/go-ipfs
cd /go/src/github.com/ipfs/go-ipfs/
make deps

# Cold Importer
go get -v github.com/beeker1121/goque
go get -v github.com/ipfs/go-cid
go get -v github.com/multiformats/go-multihash
go get -v github.com/syndtr/goleveldb/leveldb
go get -v github.com/syndtr/goleveldb/leveldb/errors
