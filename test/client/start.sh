#!/bin/bash

set -euo pipefail

VERSION="$1"

source ~/venv/bin/activate
export PATH=$PATH:/usr/local/go/bin

TMP_DIR="/tmp/calamary_$(date +%s)"
echo "$TMP_DIR" > "/tmp/calamary_${VERSION}.run"
mkdir "$TMP_DIR"
cd "$TMP_DIR"

git clone 'https://github.com/superstes/calamary'

set +e
bash "${TMP_DIR}/calamary/test/wrapper.sh" "$VERSION"

echo $? > "/tmp/calamary_${VERSION}.exit"
