#!/bin/bash

set -euo pipefail

VERSION="$1"

TMP_DIR="$(cat /tmp/calamary_${VERSION}.run)"

bash "${TMP_DIR}/calamary/test/postExec.sh" "$VERSION"

# cleanup
if echo "$TMP_DIR" | grep -q '/tmp/calamary'
then
  rm -rf "$TMP_DIR"
fi
rm "/tmp/calamary_${VERSION}.run"
rm "/tmp/calamary_${VERSION}.exit"
