# go2xunit [![Build Status](https://travis-ci.org/tischda/go2xunit.svg)](https://travis-ci.org/tischda/go2xunit)

Converts `go test -v` (or `gocheck -vv`) output to xunit compatible XML output
(used in [Jenkins][jenkins]/[Hudson][hudson]).

This is a fork from go2xunit by Miki Tebeka <miki.tebeka@gmail.com>.

# Install

    go get github.com/tischda/go2xunit


# Usage
~~~
Usage of go2xunit.exe:
  -bamboo=false: xml compatible with Atlassian's Bamboo
  -fail=false: fail (non zero exit) if any test failed
  -fail-on-race=false: mark test as failing if it exposes a data race
  -gocheck=false: parse gocheck output
  -input="": input file (default to stdin)
  -output="": output file (default to stdout)
  -version=false: print version and exit
  -xunitnet=false: xml compatible with xunit.net
~~~

By default `go2xunit` reads data from standard input and emits XML to standard
output. However you can use `-input` and `-output` flags to change this.

The `-fail` switch will cause `go2xunit` to exit with non zero status if there
are failed tests.

    go test -v | go2xunit -output tests.xml

`go2xunit` also works with [gocheck][gocheck], and [testify][testify].

    go test -gocheck.vv | go2xunit -gocheck -output tests.xml


[jenkins]: http://jenkins-ci.org/
[hudson]: http://hudson-ci.org/
[gocheck]: http://labix.org/gocheck
[testify]: http://godoc.org/github.com/stretchr/testify
[bugs]: https://bitbucket.org/tebeka/go2xunit/issues
