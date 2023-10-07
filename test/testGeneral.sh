#!/bin/bash

# tests are targeting the 'transparent' mode

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}1"}"

testsGeneral=()

for test in "${testsGeneral[@]}"
do
  if ! runTest "general/$test"
  then
    fail
  fi
done

unset PROXY_PORT
