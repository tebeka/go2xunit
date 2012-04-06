export GOPATH := $(shell dirname $(shell dirname $(PWD)))
PACKAGE := go2xunit

all:
	go build $(PACKAGE)

test:
	go test -v $(PACKAGE)

fix:
	go fix $(PACKAGE)

doc:
	go doc $(PACKAGE)

install:
	go install $(PACKAGE)

README.html: README.rst
	rst2html $< > $@

.PHONY: all test install fix doc
