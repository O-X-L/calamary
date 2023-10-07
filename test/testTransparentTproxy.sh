#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}2"}"

testsTransparentTproxy=()

if [ "${#testsTransparentTproxy}" -gt "0" ]
then
  log_header 'RUNNING TESTS: TRANSPARENT-TPROXY'
fi

for test in "${testsTransparentTproxy[@]}"
do
  if ! runTest "transparentTproxy/$test"
  then
    fail
  fi
done

unset PROXY_PORT
