#!/bin/bash
# Test script run by jenkins

outfile=gotest.out

2>&1 go test -v | tee $outfile
go2xunit -fail -input $outfile -output tests.xml
