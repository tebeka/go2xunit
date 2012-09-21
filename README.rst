==============
go2xunit 0.1.1
==============

Converts `go test -v` output to xunit compatible XML output (used in
Jenkins/Hudson). 


Install
=======
`go install bitbucket.org/tebeka/go2xunit`


Usage
=====
By default `go2xunit` reads data from standard input and emits XML to standard
output. However you can use `-input` and `-output` flags to change this.

The `-fail` switch will cause `go2xunit` to exit with non zero status if there
are failed tests.

::

    go test -v | go2xunit -output tests.xml

Here's an example script (`run-tests.sh`) that can be used with Jenkins_/Hudson_.

::
    
    #!/bin/bash

    export GOPATH=$(dirname $(dirname $PWD))
    outfile=gotest.out

    go test -v | tee $outfile
    go2xunit -fail -input $outfile -output tests.xml


.. _Jenkins: http://jenkins-ci.org/
.. _Hudson: http://hudson-ci.org/

Contact
=======
Miki Tebeka <miki.tebeka@gmail.com>

Bug reports go here_.

.. _here: https://bitbucket.org/tebeka/go2xunit/issues?status=new&status=open


