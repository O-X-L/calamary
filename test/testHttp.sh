#!/bin/bash

set -uo pipefail
set +e

PORT_BASE="${PORT_BASE:-1000}"
PROXY_PORT="${PROXY_PORT:-"${PORT_BASE}4"}"

export PROXY_HOST="$PROXY_HOST"
export PROXY_PORT="$PROXY_PORT"
export http_proxy="http://${PROXY_HOST}:${PROXY_PORT}"
export https_proxy="http://${PROXY_HOST}:${PROXY_PORT}"
export HTTP_PROXY="http://${PROXY_HOST}:${PROXY_PORT}"
export HTTPS_PROXY="http://${PROXY_HOST}:${PROXY_PORT}"

testsHttp=()

# log_header 'RUNNING TESTS: HTTP'

for test in "${testsHttp[@]}"
do
  if ! runTest "test_http/$test"
  then
    fail
  fi
done

unset PROXY_PORT
unset http_proxy
unset https_proxy
unset HTTP_PROXY
unset HTTPS_PROXY
