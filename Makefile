all: lint test

test:
	go test -v . ./lib

lint:
	golint .
	golint lib

.PHONY: test lib
