#!/bin/sh

go get -d ./...
go run ./cmd/main.go

eval "$@"