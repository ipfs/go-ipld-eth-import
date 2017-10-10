#!/bin/bash

################################################################################
#
# * Revert the changes made by ./build/convert-ipfs-deps.sh
#
################################################################################

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# echo
echo

# Take the dependencies back in the go-ipld-package
cd ${GOPATH}/src/github.com/ipfs/go-ipld-eth
gx-go hook post-test

# Transform go-ipld-eth back to normal
if [[ `uname` == 'Darwin' ]]; then
	sed -i '' 's/package plugin/package main/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
else
	sed -i 's/package plugin/package main/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
fi

# Take the dependencies back in the repository
cd $CURRENT_DIR/..
gx-go hook post-test
