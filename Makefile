.PHONY: test
test:
	go build ./src/...
	go vet ./src/...
	go test ./src/...
