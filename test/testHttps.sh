#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}5"}"

testsHttps=()

if [ "${#testsHttps}" -gt "0" ]
then
  log_header 'RUNNING TESTS: HTTPS'
fi

for test in "${testsHttps[@]}"
do
  if ! runTest "https/$test"
  then
    fail
  fi
done

unset PROXY_PORT
