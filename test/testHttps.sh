#!/bin/bash

set -uo pipefail
set +e

source ./base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}5"}"

testsHttps=()

# log_header 'RUNNING TESTS: HTTPS'

for test in "${testsHttps[@]}"
do
  if ! runTest "https/$test"
  then
    fail
  fi
done

unset PROXY_PORT
