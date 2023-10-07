#!/bin/bash

testsTransparent=()
testsTransparent[0]="basic"
testsTransparent[1]="dummyOk"
#testsTransparent[2]="dummyFail"

for test in "${testsTransparent[@]}"
do
  if ! runTest "transparent/$test"
  then
    fail
  fi
done
