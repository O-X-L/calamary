#!/bin/bash

set -e

cd "$(dirname "$0")/../main"
go build
./main
