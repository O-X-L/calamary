#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}6"}"

testsSocks5=()


for test in "${testsSocks5[@]}"
do
  if ! runTest "socks5/$test"
  then
    fail
  fi
done

unset PROXY_PORT
