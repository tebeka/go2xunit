# go2xunit

[![CircleCI](https://circleci.com/gh/tebeka/go2xunit.svg?style=svg)](https://circleci.com/gh/tebeka/go2xunit)

Converts `go test -v --json` output to [junit][junit] or [xunit.net][xnet]
compatible XML output (used in [Jenkins][jenkins]/[Hudson][hudson]).

Note that *the default is junit*, in Jenkins please pick `Publish JUnit test
result report` (not `XUnit`).


# Install

    go get github.com/tebeka/go2xunit


# Usage
By default `go2xunit` reads data from standard input and emits XML to standard
output. However you can use `-input` and `-output` flags to change this.

The `-no-fail` switch will cause `go2xunit` to not to exit with non zero status
if there are failed tests.

    2>&1 go test -v --json | go2xunit -output tests.xml

`go2xunit` also works with [gocheck][gocheck], and [testify][testify].

    2>&1 go test -gocheck.vv --json | go2xunit -output tests.xml

Here's an example script (`run-tests.sh`) that can be used with [Jenkins][jenkins]/[Hudson][hudson].

    #!/bin/bash

    outfile=gotest.out

    2>&1 go test -v | tee $outfile
    go2xunit -fail -input $outfile -output tests.xml


# Examples

* [go test](demos/gotest/)
* [gocheck](demos/gocheck/)
* [testify](demos/testify/)


# Related

* [testing: add -json flag for json results](https://github.com/golang/go/issues/2981) open bug

# Contact
Miki Tebeka <miki.tebeka@gmail.com>

Bug reports go [here][bugs].


[jenkins]: http://jenkins-ci.org/
[hudson]: http://hudson-ci.org/
[gocheck]: http://labix.org/gocheck
[testify]: http://godoc.org/github.com/stretchr/testify
[bugs]: https://github.com/tebeka/go2xunit/issues
[xnet]: https://xunit.codeplex.com/wikipage?title=XmlFormat
[junit]:  https://www.ibm.com/support/knowledgecenter/en/SSQ2R2_14.1.0/com.ibm.rsar.analysis.codereview.cobol.doc/topics/cac_useresults_junit.html
