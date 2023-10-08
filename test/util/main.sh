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
  echo "Copying file $1 => $2"
  rc=0
  for i in {1..5}
  do
    scp -P "$PROXY_SSH_PORT" "$1" "$PROXY_USER"@"$PROXY_HOST":"$2" >/dev/null 2>&1 || rc="$?"
    if [[ "$rc" == '0' ]]
    then
      return
    fi
  done
}

function ssh_cmd {
  echo "Running remote command: '$1'"
  ssh -p "$PROXY_SSH_PORT" "$PROXY_USER"@"$PROXY_HOST" "$1" >/dev/null 2>&1
}

function fail {
  echo ''
  log 'TEST-RUN FAILED!'
  status='FAILED'
  cleanup
  update_badge
  exit 99
}
