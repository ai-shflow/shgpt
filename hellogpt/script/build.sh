#!/bin/bash

build=$(date +%FT%T%z)
version="$1"

ldflags="-s -w -X github.com/cligpt/shgpt/hellogpt/config.Build=$build -X github.com/cligpt/shgpt/hellogpt/config.Version=$version"
target="hellogpt"

go env -w GOPROXY=https://goproxy.cn,direct

# go tool dist list
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "$ldflags" -o bin/$target example/example.go
CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags "$ldflags" -o bin/$target.exe example/example.go
