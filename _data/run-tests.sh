#!/bin/bash

export GOPATH=$(dirname $(dirname $PWD))
outfile=gotest.out

2>&1 go test -v | tee $outfile
go2xunit -fail -input $outfile -output tests.xml
