COVERFILE=.count

test: $(wildcard *.go **/*.go)
	go test -v ./...

cover-test:
	go test -covermode=count -coverprofile=$(COVERFILE) .
	go tool cover -html=$(COVERFILE)
