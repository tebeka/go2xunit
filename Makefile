all: lint test

test:
	go test -v ./...

lint:
	golint -set_exit_status ./...

publish:
	git tag v$(shell perl -ne 'print "$$1\n" if /Version = "(.*)"/' main.go)
	git push
	git push --tags

binaries:
	pwd
	GOARCH=amd64 GOOS=linux go build -o go2xunit-linux-amd64
	GOARCH=amd64 GOOS=darwin go build -o go2xunit-darwin-amd64
	GOARCH=amd64 GOOS=windows go build -o go2xunit-windows-amd64.exe

.PHONY: test lib
