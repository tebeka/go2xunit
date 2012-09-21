#!/bin/bash

java -jar jenkins.war \
    --logfile=${PWD}/jenkins.log \
    --httpPort=8000 \
    --daemon
