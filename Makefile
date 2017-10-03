build-cold:
	go build -o build/bin/cold-importer cold-importer/*.go

cold: hack build-cold unhack

clean:
	rm -rf build/bin/*

# There should be a better way than gx-go hook *-test
# If you know it, please share it :-)
hack:
	gx-go hook pre-test
	sed -i '' 's/package main/package plugin/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
	cd ${GOPATH}/src/github.com/ipfs/go-ipld-eth; gx-go hook pre-test

# Very ugly! Open to suggestions
# Will think on something later
unhack:
	cd ${GOPATH}/src/github.com/ipfs/go-ipld-eth; gx-go hook post-test
	sed -i '' 's/package plugin/package main/' ${GOPATH}/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
	gx-go hook post-test

.PHONY: build-cold cold clean hack unhack
