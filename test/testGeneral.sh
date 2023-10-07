#!/bin/bash

set -uo pipefail
set +e

# tests are targeting the 'transparent' mode

source ./base.sh

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
