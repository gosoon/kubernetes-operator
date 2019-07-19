#!/bin/bash

export GO111MODULE="on"
export GOPROXY="https://goproxy.io"
export CGO_ENABLED="0"
export GOOS="linux"
export GOARCH=amd64

go mod vendor
go build -ldflags "-s -w" -a -installsuffix cgo -o kubernetes-operator .
