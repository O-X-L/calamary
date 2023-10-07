#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}5"}"

testsHttps=()


for test in "${testsHttps[@]}"
do
  if ! runTest "https/$test"
  then
    fail
  fi
done

unset PROXY_PORT
