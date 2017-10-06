## TODO
## make all
## make evmcode-file
## make state-trie-file
## make state-trie-ipfs

build-evmcode-ipfs:
	go build -o build/bin/evmcode-ipfs cold-importer/evmcode-ipfs/*.go

evmcode-ipfs: hack build-evmcode-ipfs unhack

clean:
	rm -rf build/bin/*

# There should be a better way than gx-go hook *-test
# If you know it, please share it :-)
hack:
	gx-go hook pre-test
	sed -i '' 's/package main/package plugin/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
	cd ${GOPATH}/src/github.com/ipfs/go-ipld-eth; gx-go hook pre-test
	echo

# Very ugly! Open to suggestions
# Will think on something later
unhack:
	echo
	cd ${GOPATH}/src/github.com/ipfs/go-ipld-eth; gx-go hook post-test
	sed -i '' 's/package plugin/package main/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
	gx-go hook post-test

.PHONY: build-evmcode-ipfs evmcode-ipfs clean hack unhack
