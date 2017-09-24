#!/bin/bash

################################################################################
#
# The reason we are doing this here and not in the Dockerfile is that
# you may want to update your packages without wanting to build the image
# everytime.
#
################################################################################

go get -v github.com/ipfs/go-ipfs
cd /go/src/github.com/ipfs/go-ipfs/
make deps

go get -v github.com/ethereum/go-ethereum
