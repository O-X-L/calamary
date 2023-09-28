#!/bin/bash

cd "$(dirname "$0")/../lib"
go version
go test -v ./...
