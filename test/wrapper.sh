#!/bin/bash

# this is the main entrypoint for testing
# it will run the tests of the currently checked-out version on the provided lib-version

set -eo pipefail

echo ''

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
status='RUNNING'
PATH_BADGE='/var/www/cicd/calamary'
BADGE_LABEL="Integration Tests - ${VERSION}"

declare -A BADGE_COLORS
BADGE_COLORS[UNKNOWN]='#404040'
BADGE_COLORS[RUNNING]='#1f77aa'
BADGE_COLORS[PASSED]='#97CA00'
BADGE_COLORS[FAILED]='#d9644d'
BADGE_COLORS[FAILED-ENVIRONMENT]='#d9644d'

function log {
  echo "$1"
}

function update_badge {
  if [ -d "$PATH_BADGE" ]
  then
    cd "$PATH_BADGE"
    rm -f "${VERSION}.calamary.test.svg"
    anybadge --label="$BADGE_LABEL" --value="$status | $(date '+%Y-%m-%d %H:%M') GMT+2" --file="${VERSION}.calamary.test.svg" --color="${BADGE_COLORS[$status]}"
    cd "$WD"
  fi
}

cd "$(dirname "$0")"
WD="$(pwd)"

update_badge

mkdir -p "$TMP_DIR"
cp -r ./* "${TMP_DIR}/"

tar -C "${TMP_DIR}/" -xzf ./tools/easyrsa.tgz

cd ..
REPO_DIR="$(pwd)"
VERSION_TEST_TAG="$(git rev-parse --abbrev-ref HEAD)"
VERSION_TEST_COMMIT="$(git rev-parse HEAD)"
VERSION_TEST="${VERSION_TEST_TAG}-${VERSION_TEST_COMMIT:0:8}"

log "TESTING VERSION '${VERSION}' WITH TEST-VERSION '${VERSION_TEST}'"

log "BUILDING BINARY (${TMP_BIN})"

if [[ "$VERSION_TEST_TAG" != "$VERSION" ]]
then
  git checkout "$VERSION"
fi
cd lib/
go mod download
cd main/
go build -o "$TMP_BIN"
chmod +x "$TMP_BIN"

cd "$TMP_DIR"
WD="$(pwd)"

# start actual testing
log 'STARTING TESTS'
source ./main.sh

echo ''
