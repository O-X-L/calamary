#!/bin/bash

PROXY_HOST='172.17.1.81'
PROXY_USER='tester'
PROXY_SSH_PORT=22

function ssh_cmd {
  echo "Running remote command: '$1'"
  ssh -p "$PROXY_SSH_PORT" "$PROXY_USER"@"$PROXY_HOST" "$1" >/dev/null 2>&1
  return "$?"
}
