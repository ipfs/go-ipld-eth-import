#!/bin/bash

################################################################################
#
# There are better ways to do this. This is the ugly shortcut
#
################################################################################

go get -v -u github.com/whyrusleeping/gx
go get -v -u github.com/whyrusleeping/gx-go

go get -v -u github.com/ipfs/go-ipld-eth

go get -v -u github.com/ipfs/go-ipfs
cd $GOPATH/src/github.com/ipfs/go-ipfs
make build

cd $GOPATH/src/github.com/ipfs/go-ipld-eth-import
./build/convert-ipfs-deps.sh
go get ./...
