#!/bin/bash

set -uo pipefail
set +e

source ./util/base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}2"}"

testsTransparentTproxy=()

# log_header 'RUNNING TESTS: TRANSPARENT-TPROXY'

for test in "${testsTransparentTproxy[@]}"
do
  if ! runTest "transparentTproxy/$test"
  then
    fail
  fi
done

unset PROXY_PORT
