cold:
	go build -o build/bin/cold-importer cold-importer/*.go

clean:
	rm -rf build/bin/*

docker-cold:
	 ./docker/build-cold-importer-image.sh

.PHONY: cold clean docker-cold
