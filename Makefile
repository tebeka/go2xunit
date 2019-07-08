all: lint test

test:
	go test -v ./...

lint:
	staticcheck ./...

publish:
	git tag v$(shell perl -ne 'print "$$1\n" if /Version = "(.*)"/' main.go)
	git push
	git push --tags

circleci:
	docker build -f Dockerfile.test .

.PHONY: test lib
