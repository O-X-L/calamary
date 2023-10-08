#!/bin/bash

set -uo pipefail
set +e

PORT_BASE="${PORT_BASE:-1000}"
PROXY_PORT="${PROXY_PORT:-"${PORT_BASE}3"}"

export PROXY_HOST="$PROXY_HOST"
export PROXY_PORT="$PROXY_PORT"

testsProxyproto=()

# log_header 'RUNNING TESTS: PROXYPROTO'

for test in "${testsProxyproto[@]}"
do
  if ! runTest "test_proxyproto/$test"
  then
    fail
  fi
done

unset PROXY_PORT
