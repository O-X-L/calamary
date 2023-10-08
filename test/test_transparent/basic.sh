#!/bin/bash

source ../util/test.sh
source ../util/route.sh

target='1.1.1.1'
route_add "$target"
c1=$(curlRc "http://${target}")
c2=$(curlRc "https://${target}")
route_rm "$target"

if [[ "$c1" != "0" ]] || [[ "$c2" != "0" ]]
then
  exit 1
fi
