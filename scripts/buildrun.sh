#!/bin/bash

set -e

binary="/tmp/calamary_$(date +"%s")"
echo "TEST: Binary: '${binary}'"

cd "$(dirname "$0")/../lib/main"
go build -o "$binary"
chmod +x "$binary"
sudo su proxy --shell /bin/bash -c "$binary"
