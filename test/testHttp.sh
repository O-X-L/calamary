#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}4"}"

testsHttp=()


for test in "${testsHttp[@]}"
do
  if ! runTest "http/$test"
  then
    fail
  fi
done

unset PROXY_PORT
