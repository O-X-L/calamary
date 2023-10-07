#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}3"}"

testsProxyproto=()

# log_header 'RUNNING TESTS: PROXYPROTO'

for test in "${testsProxyproto[@]}"
do
  if ! runTest "proxyproto/$test"
  then
    fail
  fi
done

unset PROXY_PORT
