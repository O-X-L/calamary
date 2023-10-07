#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")"
source ./target.sh

VERSION="$1"
PATH_BADGE='/var/www/cicd/calamary'
BADGE_STATUS='FAILED-ENVIRONMENT'
BADGE_COLOR='#d9644d'
BADGE_LABEL="Integration Tests - ${VERSION}"

function update_badge {
  if [ -d "$PATH_BADGE" ]
  then
    cd "$PATH_BADGE"
    rm -f "${VERSION}.calamary.test.svg"
    anybadge --label="$BADGE_LABEL" --value="$BADGE_STATUS | $(date '+%Y-%m-%d %H:%M') GMT+2" --file="${VERSION}.calamary.test.svg" --color="$BADGE_COLOR"
  fi
}

if [[ cat "/tmp/calamary_${VERSION}.exit" != "0" ]]
then
  update_badge
fi
