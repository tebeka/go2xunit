========
go2xunit
========

Converts `go test -v` output to xunit compatible XML output. 


Install
=======
`go install bitbucket.org/tebeka/go2xunit`


Usage
=====
By default `go2xml` reads data from standard input and emits XML to standard
output. However you can use `-input` and `-output` flags to change this.

::

    go test -v | go2xunit -output tests.xml

Contact
=======
Miki Tebeka <miki.tebeka@gmail.com>

Bug reports go here_.

.. _here: https://bitbucket.org/tebeka/go2xunit/issues?status=new&status=open


