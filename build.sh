#!/bin/bash

mkdir -p builds

# linux 64
echo building linux
env GOOS=linux GOARCH=amd64 go build -a -o builds/linux64/pwmm

# mac 64
echo building mac
env GOOS=darwin GOARCH=amd64 go build -a -o builds/macos64/pwmm
