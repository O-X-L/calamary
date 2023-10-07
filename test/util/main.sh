#!/bin/bash

function log {
  echo "$1"
}

function log_header {
  echo "##### $1 #####"
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

function copy_file {
    scp -P "$PROXY_SSH_PORT" "$1" "$PROXY_USER"@"$PROXY_HOST":"$2" >/dev/null 2>&1
}

function fail {
  log 'TEST-RUN FAILED!'
  status='FAILED'
  stop_proxy
  update_badge
  exit 99
}
