#!/bin/bash

# executing test-script as new process
function runTest {
  testScript="$1"
  echo ''
  echo "RUNNING TEST '${testScript}'"
  echo ''
  bash ./${testScript}.sh
  result="$?"
  if [[ "result" -ne "0" ]]
  then
    echo "FAILED TEST '${testScript}'"
    return 1
  fi
  echo ''
  return 0
}
