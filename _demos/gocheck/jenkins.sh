#!/bin/bash
# Test script run by jenkins

outfile=gotest.out

2>&1 go test -gocheck.vv | tee $outfile
go2xunit -gocheck -fail -input $outfile -output tests.xml
