#!/bin/bash

set -euo pipefail

source ./target.sh

TMP_BASE="/tmp/calamary_${VERSION}"  # could be problematic
PORT_BASE="$(date +'%H%M')"
CERT_CN="/C=AT/ST=Styria/CN=Calamary Forward Proxy"

# remove leading 0 as it is not valid as port
if [[ ${PORT_BASE:0:1} == "0" ]]
then
  PORT_BASE="3${PORT_BASE:1}"
fi

function log {
  echo ''
  echo "$1"
  echo ''
}

function stop_proxy {
  log 'STOPPING PROXY'
  ssh_cmd "sudo systemctl stop calamary@${VERSION}.service"
}

function cleanup {
  log 'CLEANUP'
  ssh_cmd "sudo rm -f ${TMP_BASE}*"
  rm -f ./*_tmp.*
  stop_proxy
}

cleanup

log 'PREPARING FOR TESTS'

cp 'config.yml' 'config_tmp.yml'

sed -i "s@PORT_BASE@$PORT_BASE@g" 'config_tmp.yml'
sed -i "s@CRT_BASE@$TMP_BASE@g" 'config_tmp.yml'

log 'GENERATING CERTS'
# todo: generate ca & subca
openssl req -x509 -newkey rsa:4096 -keyout 'cert_tmp.key' -out 'cert_tmp.crt' -sha256 -days 60 -nodes -subj "$CERT_CN" 2>/dev/null

log 'COPYING FILES TO PROXY-HOST'
function copy_file {
    scp -P "$PROXY_SSH_PORT" "$1" "$PROXY_USER"@"$PROXY_HOST":"$2" >/dev/null 2>&1
}

copy_file 'calamary' "$TMP_BASE"
copy_file 'config_tmp.yml' "${TMP_BASE}.yml"
copy_file 'cert_tmp.key' "${TMP_BASE}.key"
copy_file 'cert_tmp.crt' "${TMP_BASE}.crt"
ssh_cmd "sudo chown proxy:proxy ${TMP_BASE}*"

log 'STARTING PROXY'
ssh_cmd "sudo systemctl start calamary@${VERSION}.service"

function runTest {
  testScript="$1"
  echo ''
  echo "RUNNING TEST '${testScript}'"
  echo ''
  ./${testScript}.sh
  result="$?"
  if [[ "result" -ne "0" ]]
  then
    echo "FAILED TEST '${testScript}'"
    return 1
  fi
  echo ''
  return 0
}

function fail {
  log 'TEST-RUN FAILED!'
  status='FAILED'
  stop_proxy
  update_badge
  exit 99
}

log 'STARTING TESTS'

sed +e
find . -type f -name '*.sh' -exec chmod +x {} \;
source testTransparent.sh
source testGeneral.sh
source testTransparentTproxy.sh
source testHttp.sh
source testHttps.sh
source testProxyproto.sh
source testSocks5.sh

log 'TEST-RUN FINISHED SUCCESSFULLY!'
status='PASSED'

cleanup
update_badge

exit 0
