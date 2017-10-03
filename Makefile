# There should be a better way than gx-go hook *-test
# If you know it, please share it :-D
cold:
	gx-go hook pre-test
	go build -o build/bin/cold-importer cold-importer/*.go
	gx-go hook post-test

clean:
	rm -rf build/bin/*

# Very ugly! Open to suggestions
ungx:
	gx-go hook post-test

.PHONY: cold clean
