all: lint test

test:
	go test -v ./...

lint:
	golint -set_exit_status ./...

.PHONY: test lib
