test: $(wildcard *.go **/*.go)
	go test -v ./...
