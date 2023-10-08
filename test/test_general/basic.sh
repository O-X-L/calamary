#!/bin/bash

source "$(pwd)/util/test.sh"
source "$(pwd)/util/route.sh"

results=()



exit_code="$(anyFailed "${results[@]}")"
exit "$exit_code"
