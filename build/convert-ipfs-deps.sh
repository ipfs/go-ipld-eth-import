#!/bin/bash

################################################################################
#
# * Convert the libraries using gx, to the ones declared in package.json
# * Modify the go-ipld-eth library, to be included in this package
#
################################################################################

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Update the dependencies of this repository accordingly to package.json
cd $CURRENT_DIR/..
gx-go hook pre-test

# Transform go-ipld-eth to a library
if [[ `uname` == 'Darwin' ]]; then
	sed -i '' 's/package main/package plugin/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
else
	sed -i 's/package main/package plugin/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
fi

# Update the dependencies of go-ipld-eth accordingly to package.json
cd ${GOPATH}/src/github.com/ipfs/go-ipld-eth
gx-go hook pre-test

# echo
echo
