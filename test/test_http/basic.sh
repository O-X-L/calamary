#!/bin/bash

source "$(pwd)/util/test.sh"

results=()



exit_code="$(anyFailed "${results[@]}")"
exit "$exit_code"
