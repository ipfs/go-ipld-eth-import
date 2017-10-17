## TODO
## make all
## make state-trie-file
## make state-trie-ipfs

all: evmcode-file evmcode-ipfs state-trie-file

clean:
	rm -rf build/bin/*

clean-deps:
	build/un-convert-ipfs-deps.sh

evmcode-file:
	build/convert-ipfs-deps.sh
	go build -v -o build/bin/evmcode-file cold-importer/evmcode-file/*.go
	build/un-convert-ipfs-deps.sh

evmcode-ipfs:
	build/convert-ipfs-deps.sh
	go build -v -o build/bin/evmcode-ipfs cold-importer/evmcode-ipfs/*.go
	build/un-convert-ipfs-deps.sh

state-trie-file:
	build/convert-ipfs-deps.sh
	go build -v -o build/bin/state-trie-file cold-importer/state-trie-file/*.go
	build/un-convert-ipfs-deps.sh

vet:
	build/convert-ipfs-deps.sh
	unused ./...
	staticcheck ./...
	gosimple ./...
	golint ./...
	build/un-convert-ipfs-deps.sh

.PHONY: all lean clean-deps evmcode-file evmcode-ipfs state-trie-file vet
