#!/bin/bash

PROXY_PORT="${PROXY_PORT:="${PORT_BASE}2"}"

testsTransparentTproxy=()


for test in "${testsTransparentTproxy[@]}"
do
  if ! runTest "transparentTproxy/$test"
  then
    fail
  fi
done

unset PROXY_PORT
