FROM golang:1.12-alpine
ENV GOCACHE=/tmp/gocache
WORKDIR /workdir
ADD . /workdir
RUN apk add --no-cache git \
    && go build -o pwmm main.go

ENTRYPOINT ["/workdir/pwmm"]