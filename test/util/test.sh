#!/bin/bash

function curlRc {
  echo "HTTP Call to $1"
  curl -ss --connect-timeout 2 --fail "$1"
  echo "$?"
}

function anyFailed {
  results=("$@")
  for result in "${results[@]}"
  do
    if [[ "$result" != '0' ]]
    then
      echo '1'
      return
    fi
  done
  echo '0'
}

