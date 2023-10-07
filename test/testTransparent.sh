#!/bin/bash

set -uo pipefail
set +e

source ./util/base.sh
source ./util/route.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}1"}"

testsTransparent=()
testsTransparent[0]="basic"

log_header 'RUNNING TESTS: TRANSPARENT'

for test in "${testsTransparent[@]}"
do
  if ! runTest "transparent/$test"
  then
    fail
  fi
done

unset PROXY_PORT
