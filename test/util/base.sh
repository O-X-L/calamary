#!/bin/bash

source ./route.sh

function runTest {
  testScript="$1"
  echo ''
  echo "RUNNING TEST '${testScript}'"
  echo ''
  ./${testScript}.sh
  result="$?"
  if [[ "result" -ne "0" ]]
  then
    echo "FAILED TEST '${testScript}'"
    return 1
  fi
  echo ''
  return 0
}

function curlRc {
  curl --connect-timeout 3 --fail "$1"
  return "$?"
}
