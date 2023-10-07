#!/bin/bash

set -uo pipefail
set +e

source ./util/base.sh

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}6"}"

export http_proxy="socks5://${PROXY_HOST}:${PROXY_PORT}"
export https_proxy="socks5://${PROXY_HOST}:${PROXY_PORT}"
export HTTP_PROXY="socks5://${PROXY_HOST}:${PROXY_PORT}"
export HTTPS_PROXY="socks5://${PROXY_HOST}:${PROXY_PORT}"

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
