all: lint test

test:
	go test -v . ./lib

lint:
	golint -set_exit_status .
	golint -set_exit_status lib

.PHONY: test lib
