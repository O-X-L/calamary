#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}3"}"

testsProxyproto=()

if [ "${#testsProxyproto}" -gt "0" ]
then
  log_header 'RUNNING TESTS: PROXYPROTO'
fi

for test in "${testsProxyproto[@]}"
do
  if ! runTest "proxyproto/$test"
  then
    fail
  fi
done

unset PROXY_PORT
