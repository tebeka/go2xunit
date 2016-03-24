#!/bin/bash
# Push changes to github

case $1 in
    -h | --help ) printf "usage: %s\n\nPush to github\n" $(basename $0); exit;;
esac

hg bookmark -f -r default master
hg push git+ssh://git@github.com/tebeka/go2xunit.git
