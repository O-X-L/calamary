#!/bin/bash

set -euo pipefail

VERSION="$1"

source ~/venv/bin/activate
export PATH=$PATH:/usr/local/go/bin

TMP_DIR_REPO="/tmp/calamary_$(date +%s)"
echo "$TMP_DIR_REPO" > "/tmp/calamary_${VERSION}.run"
mkdir "$TMP_DIR_REPO"
cd "$TMP_DIR_REPO"

git clone 'https://github.com/superstes/calamary'

set +e
bash "${TMP_DIR_REPO}/calamary/test/wrapper.sh" "$VERSION"

ec="$?"
echo "$ec" > "/tmp/calamary_${VERSION}.exit"
exit "$ec"
