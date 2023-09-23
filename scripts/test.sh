#!/bin/bash

cd "$(dirname "$0")/.."
go version
go test -v ./...
