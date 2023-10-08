#!/bin/bash

set -uo pipefail
set +e

PORT_BASE="${PORT_BASE:-1000}"
PROXY_PORT="${PROXY_PORT:-"${PORT_BASE}1"}"

export PROXY_HOST="$PROXY_HOST"
export PROXY_PORT="$PROXY_PORT"

testsTransparent=()
testsTransparent[0]="basic"

log_header 'RUNNING TESTS: TRANSPARENT'

for test in "${testsTransparent[@]}"
do
  if ! runTest "test_transparent/$test"
  then
    fail
  fi
done

unset PROXY_PORT
