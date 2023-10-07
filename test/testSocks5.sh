#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}6"}"

testsSocks5=()

if [ "${#testsSocks5}" -gt "0" ]
then
  log_header 'RUNNING TESTS: SOCKS5'
fi

for test in "${testsSocks5[@]}"
do
  if ! runTest "socks5/$test"
  then
    fail
  fi
done

unset PROXY_PORT
