## TODO
## make all
## make evmcode-file
## make state-trie-file
## make state-trie-ipfs

clean:
	rm -rf build/bin/*

clean-deps:
	build/un-convert-ipfs-deps.sh

evmcode-ipfs:
	build/convert-ipfs-deps.sh
	go build -o build/bin/evmcode-ipfs cold-importer/evmcode-ipfs/*.go
	build/un-convert-ipfs-deps.sh

.PHONY: clean clean-deps evmcode-ipfs
