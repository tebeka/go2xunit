#!/bin/bash

case $1 in
    -h | --help ) echo "usage: $(basename $0) FILE [FILE...]"; exit;;
esac

if [ $# -lt 1 ]; then
    1>&2 echo "error: wrong number of arguments"
    exit 1
fi

for name in $@; do
    if [ ! -f "${name}" ]; then
	1>&2 echo "error: ${name} - no such file"
	exit 1
    fi

    out=_data/out/xunit/$(basename ${name}).xml
    echo "${name} -> ${out}"
    go run . -input ${name} -output ${out}
    out=_data/out/xunit.net/$(basename ${name}).xml
    echo "${name} -> ${out}"
    go run . -input ${name} -output ${out} -xunitnet
done
