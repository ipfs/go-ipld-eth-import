#!/bin/bash

go get -v github.com/ipfs/go-ipfs
cd /go/src/github.com/ipfs/go-ipfs/
make deps

go get -v github.com/ethereum/go-ethereum
