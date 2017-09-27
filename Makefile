all: lint test

test:
	go test -v ./...

lint:
	golint -set_exit_status ./...

publish:
	git tag v$(shell perl -ne 'print "$$1\n" if /Version = "(.*)"/' main.go)
	git push
	git push --tags

.PHONY: test lib
