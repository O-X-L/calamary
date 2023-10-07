#!/bin/bash

set -uo pipefail
set +e

source ./base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}4"}"

testsHttp=()

# log_header 'RUNNING TESTS: HTTP'

for test in "${testsHttp[@]}"
do
  if ! runTest "http/$test"
  then
    fail
  fi
done

unset PROXY_PORT
