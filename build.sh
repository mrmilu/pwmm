#!/bin/bash
function dockerway()
{
echo building with docker
echo building mac
docker run --rm -w /workdir -v $(pwd):/workdir -u 1000 -e GOCACHE=/tmp/cache -e GOOS=darwin -e GOARCH=amd64 golang go build -a -o builds/macos64/pwmm
echo building linux
docker run --rm -w /workdir -v $(pwd):/workdir -u 1000 -e GOCACHE=/tmp/cache -e GOOS=linux -e GOARCH=amd64 golang go build -a -o builds/linux64/pwmm
}

function nodocker()
{

echo building linux
env GOOS=linux GOARCH=amd64 go build -a -o builds/linux64/pwmm

# mac 64
echo building mac
env GOOS=darwin GOARCH=amd64 go build -a -o builds/macos64/pwmm 
}


mkdir -p builds
if which docker
then
    dockerway
else
    nodocker
fi