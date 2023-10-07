#!/bin/bash

set -uo pipefail
set +e

source ./base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}6"}"

testsSocks5=()

# log_header 'RUNNING TESTS: SOCKS5'

for test in "${testsSocks5[@]}"
do
  if ! runTest "socks5/$test"
  then
    fail
  fi
done

unset PROXY_PORT
