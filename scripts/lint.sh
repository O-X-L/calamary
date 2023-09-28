#!/bin/bash

cd "$(dirname "$0")/../lib"
golangci-lint run --config ../.golangci.yml
