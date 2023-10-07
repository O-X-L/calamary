#!/bin/bash

set -euo pipefail

VERSION="$1"

source ~/venv/bin/activate

TMP_DIR_REPO="$(cat /tmp/calamary_${VERSION}.run)"

bash "${TMP_DIR_REPO}/calamary/test/postExec.sh" "$VERSION"

# cleanup
if echo "$TMP_DIR_REPO" | grep -q '/tmp/calamary'
then
  rm -rf "$TMP_DIR_REPO"
fi
rm "/tmp/calamary_${VERSION}.run"
rm "/tmp/calamary_${VERSION}.exit"
