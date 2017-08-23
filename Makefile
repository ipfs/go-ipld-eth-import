cold:
	go build -o build/bin/cold-importer cold-importer/*.go

clean:
	rm -rf build/bin/*

.PHONY: cold clean
