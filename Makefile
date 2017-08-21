cold:
	go build -o build/bin/cold-importer cold-importer/main.go

run-cold:
	./build/bin/cold-importer

clean:
	rm -rf build/bin/*

.PHONY: cold run-cold clean
