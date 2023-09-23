#!/bin/bash

cd "$(dirname "$0")/.."
golangci-lint run --config .golangci.yml
