#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

################################################################################
#
# Build the executable
#
################################################################################

docker run \
	-ti --rm \
	--name golang-ipfs-compiler \
	-v $HOME/.golang-linux-ipfs-compiler:/go \
	-v $DIR/..:/go/src/github.com/ipfs/go-ipld-eth-import \
	-w /go/src/github.com/ipfs/go-ipld-eth-import  \
	golang-ipfs-compiler go build -o build/bin/docker-cold-importer cold-importer/*.go

################################################################################
#
# Make the image
#
################################################################################

docker build -t go-ipld-eth-cold-importer -f ./docker/cold-importer.Dockerfile .
