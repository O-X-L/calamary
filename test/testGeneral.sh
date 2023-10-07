#!/bin/bash

# tests are targeting the 'transparent' mode

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}1"}"

testsGeneral=()

# log_header 'RUNNING TESTS: GENERAL'

for test in "${testsGeneral[@]}"
do
  if ! runTest "general/$test"
  then
    fail
  fi
done

unset PROXY_PORT
