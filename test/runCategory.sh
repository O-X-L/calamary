#!/bin/bash

set -uo pipefail
set +e

source ./target.sh
source ./util/main.sh
source ./util/base.sh

# overwrite - without badge
function fail {
  echo ''
  log 'TEST-RUN FAILED!'
  stop_proxy
  exit 99
}

category="$1"
VERSION="$2"

port_base="${PORT_BASE:=1000}"
echo "TESTING WITH PROXY: ${PROXY_HOST}:${PROXY_PORT:-"${port_base}?"}"

script="$(ls | grep "${category:1}.sh")"
echo "RUNNING: ${script}"
echo ''
source $script
