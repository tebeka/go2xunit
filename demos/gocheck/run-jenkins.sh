#!/bin/bash
# Run Jenkins as daemon on port 8000

# You can download jenkins.war from http://mirrors.jenkins-ci.org/war/latest/jenkins.war

java -jar jenkins.war \
    --logfile=${PWD}/jenkins.log \
    --httpPort=8000 \
    --daemon
