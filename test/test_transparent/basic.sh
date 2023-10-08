#!/bin/bash

source "$(dirname "$0")/../../util/test.sh"
source "$(dirname "$0")/../../util/route.sh"

target='1.1.1.1'
route_add "$target"
c1=$(curlRc "http://${target}")
c2=$(curlRc "https://${target}")
route_rm "$target"

if [[ "$c1" != "0" ]] || [[ "$c2" != "0" ]]
then
  exit 1
fi
