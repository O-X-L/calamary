#!/bin/bash

set -uo pipefail
set +e

source ./base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}1"}"

testsTransparent=()
testsTransparent[0]="basic"
testsTransparent[1]="dummyOk"
#testsTransparent[2]="dummyFail"

log_header 'RUNNING TESTS: TRANSPARENT'

for test in "${testsTransparent[@]}"
do
  if ! runTest "transparent/$test"
  then
    fail
  fi
done

unset PROXY_PORT
