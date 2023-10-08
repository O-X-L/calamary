#!/bin/bash

# executing test-script as new process
function runTest {
  testScript="$1"
  echo '' > "$(tty)"
  echo "RUNNING TEST '${testScript}'" > "$(tty)"
  echo '' > "$(tty)"
  bash ./${testScript}.sh
  result="$?"
  if [[ "result" -ne "0" ]]
  then
    echo "FAILED TEST '${testScript}'" > "$(tty)"
    return 1
  fi
  echo ''
  return 0
}
