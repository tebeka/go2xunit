#!/bin/bash
# Test script run by jenkins

outfile=gotest.out

go test -gocheck.vv | tee $outfile
go2xunit -gocheck -fail -input $outfile -output tests.xml
