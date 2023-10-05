#!/bin/bash

# this is the main entrypoint for testing
# it will run the tests of the currently checked-out version on the provided lib-version

set -eo pipefail

if [ -z "$1" ]
then
  echo ''
  echo 'USAGE:'
  echo ' 1 > Version to test'
  echo ''
  exit 1
fi

set -u
VERSION="$1"

TMP_DIR="/tmp/calamary_$(date +%s)"
TMP_BIN="${TMP_DIR}/calamary"

function log {
  echo ''
  echo "$1"
  echo ''
}

cd "$(dirname "$0")"
mkdir -p "$TMP_DIR"
cp -r ./* "${TMP_DIR}/"

cd ..
REPO_DIR="$(pwd)"
VERSION_TEST="$(git rev-parse HEAD)"

log "TESTING VERSION '${VERSION}' WITH TEST-VERSION '${VERSION_TEST}'"

log "BUILDING BINARY (${TMP_BIN})"

git checkout "$VERSION"
cd lib/
go mod download
cd main/
go build -o "$TMP_BIN"

cd "$TMP_DIR"

# start actual testing
source ./main.sh
