all: lint test

test: build
	go test -v . ./lib


build:
	go build ./...
	go build

lint:
	golint -set_exit_status .
	golint -set_exit_status lib

.PHONY: test lib
