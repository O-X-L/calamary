#!/bin/bash

source "$(pwd)/util/test.sh"
source "$(pwd)/util/route.sh"

results=()

target='1.1.1.1'
route_add "$target"
results[0]=$(curlRc "http://${target}")
results[1]=$(curlRc "https://${target}")
route_rm "$target"


exit_code="$(anyFailed "${results[@]}")"
exit "$exit_code"
