test:
	cd lib && go test -v
	go test -v


lint:
	golint .
	golint lib


.PHONY: test lib
