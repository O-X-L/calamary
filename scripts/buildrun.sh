#!/bin/bash

set -e

cd "$(dirname "$0")/../lib/main"
go build -o calamary
./calamary
