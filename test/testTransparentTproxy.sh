#!/bin/bash

set -uo pipefail
set +e

PORT_BASE="${PORT_BASE:-1000}"
PROXY_PORT="${PROXY_PORT:-"${PORT_BASE}2"}"

export PROXY_HOST="$PROXY_HOST"
export PROXY_PORT="$PROXY_PORT"

testsTransparentTproxy=()

# log_header 'RUNNING TESTS: TRANSPARENT-TPROXY'

for test in "${testsTransparentTproxy[@]}"
do
  if ! runTest "test_transparentTproxy/$test"
  then
    fail
  fi
done

unset PROXY_PORT
